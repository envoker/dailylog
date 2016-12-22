package dailylog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/toelsiba/date"
)

type fileWriter struct {
	dir string
	pfn paramsFileName

	d    date.Date
	file *os.File
	bw   *bufio.Writer
}

func newFileWriter(dir string, pfn paramsFileName) *fileWriter {
	return &fileWriter{
		dir: dir,
		pfn: pfn,
	}
}

func (fw *fileWriter) write(data []byte) (int, error) {

	t := time.Now()
	d, _ := date.FromTime(t)

	if fw.file != nil {
		if !d.Equal(fw.d) {
			fw.Close()
		}
	}

	if fw.file == nil {

		fileName := filepath.Join(fw.dir, compileFileName(fw.pfn, d))

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return 0, fmt.Errorf("dailylog: open file error: %s", err.Error())
		}

		fw.d = d
		fw.file = file
		fw.bw = bufio.NewWriter(file)
	}

	return fw.bw.Write(data)
}

func (fw *fileWriter) Close() error {

	if fw.bw != nil {
		fw.bw.Flush()
		fw.bw = nil
	}

	if fw.file != nil {
		err := fw.file.Close()
		fw.file = nil
		return err
	}

	return nil
}

func (fw *fileWriter) Flush() error {
	if fw.bw != nil {
		if err := fw.bw.Flush(); err != nil {
			return err
		}
	}
	return nil
}
