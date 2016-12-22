package dailylog

import (
	"errors"
	"io"
	"os"
	"sync"
)

// Intervals in minutes
type Config struct {
	Dirname        string `json:"dirname" yaml:"dirname"`
	FilePrefix     string `json:"file-prefix" yaml:"file-prefix"`
	FileExt        string `json:"file-ext" yaml:"file-ext"`
	KeepMaxDays    int    `json:"keep-max-days" yaml:"keep-max-days"`
	RotateInterval int    `json:"rotate-interval" yaml:"rotate-interval"`
	FlushInterval  int    `json:"flush-interval" yaml:"flush-interval"`
}

type Writer struct {
	m    sync.Mutex
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
		go loopRotate(w.wg, w.quit, c.RotateInterval, w)
	}
	if c.FlushInterval > 0 {
		w.wg.Add(1)
		go loopFlush(w.wg, w.quit, c.FlushInterval, w)
	}

	return w, nil
}

var errWriterIsClosed = errors.New("dailylog: Writer is closed or not created")

var _ io.WriteCloser = &Writer{}

func (w *Writer) Close() error {

	w.m.Lock()
	defer w.m.Unlock()

	if w.fw == nil {
		return errWriterIsClosed
	}

	// Stop worker
	close(w.quit)
	w.wg.Wait()

	err := w.fw.Close()
	w.fw = nil
	return err
}

func (w *Writer) Write(data []byte) (n int, err error) {

	w.m.Lock()
	defer w.m.Unlock()

	if w.fw == nil {
		return 0, errWriterIsClosed
	}

	return w.fw.write(data)
}

func (w *Writer) Rotate() error {

	w.m.Lock()
	defer w.m.Unlock()

	if w.fw == nil {
		return errWriterIsClosed
	}

	w.fw.Close()

	return w.rm.Run() // Remove old
}

func (w *Writer) Flush() error {

	w.m.Lock()
	defer w.m.Unlock()

	if w.fw == nil {
		return errWriterIsClosed
	}

	return w.fw.Flush()
}
