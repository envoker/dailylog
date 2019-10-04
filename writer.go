package dailylog

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/envoker/dailylog/period"
)

// Intervals in minutes
type Config struct {
	Dirname        string        `json:"dirname"         yaml:"dirname"         toml:"dirname"`
	FilePrefix     string        `json:"file-prefix"     yaml:"file-prefix"     toml:"file-prefix"`
	FileExt        string        `json:"file-ext"        yaml:"file-ext"        toml:"file-ext"`
	KeepMaxDays    int           `json:"keep-max-days"   yaml:"keep-max-days"   toml:"keep-max-days"`
	RotateInterval period.Period `json:"rotate-interval" yaml:"rotate-interval" toml:"rotate-interval"`
}

type Writer struct {
	guard    sync.Mutex
	quit chan struct{}
	wg   *sync.WaitGroup
	fw   *fileWriter
	rm   *oldRemover
}

func New(c Config) (*Writer, error) {

	if err := os.MkdirAll(c.Dirname, os.ModePerm); err != nil {
		return nil, err
	}

	pfn := paramsFileName{
		prefix:     c.FilePrefix,
		dateFormat: defaultDateFormat,
		ext:        c.FileExt,
	}

	w := &Writer{
		quit: make(chan struct{}),
		wg:   new(sync.WaitGroup),
		fw:   newFileWriter(c.Dirname, pfn),
		rm:   newOldRemover(c.Dirname, c.KeepMaxDays, pfn),
	}

	if c.RotateInterval > 0 {
		w.wg.Add(1)
		go rotateWorker(w.wg, w.quit, c.RotateInterval.Duration(), w)
	}

	return w, nil
}

var errWriterIsClosed = errors.New("dailylog: Writer is closed or not created")

var _ io.WriteCloser = &Writer{}

func (w *Writer) Close() error {

	w.guard.Lock()
	defer w.guard.Unlock()

	if w.fw == nil {
		return errWriterIsClosed
	}

	// Stop worker
	close(w.quit)
	w.wg.Wait()

	err := w.fw.close()
	w.fw = nil
	return err
}

func (w *Writer) Write(data []byte) (n int, err error) {

	w.guard.Lock()
	defer w.guard.Unlock()

	if w.fw == nil {
		return 0, errWriterIsClosed
	}

	return w.fw.write(data)
}

func (w *Writer) Rotate() error {

	w.guard.Lock()
	defer w.guard.Unlock()

	if w.fw == nil {
		return errWriterIsClosed
	}

	w.fw.close()

	return w.rm.Run() // Remove old
}
