package dailylog

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/toelsiba/date"
)

type oldRemover struct {
	dir         string
	keepMaxDays int
	pfn         paramsFileName
}

func newOldRemover(dir string, keepMaxDays int, pfn paramsFileName) *oldRemover {
	return &oldRemover{
		dir:         dir,
		keepMaxDays: keepMaxDays,
		pfn:         pfn,
	}
}

func (rm *oldRemover) Run() error {

	files, err := ioutil.ReadDir(rm.dir)
	if err != nil {
		return err
	}

	dateNow := date.CurrentDate()

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		dateFile, err := parseDate(rm.pfn, fileName)
		if err != nil {
			continue
		}
		if dateFile.DaysTo(dateNow) < rm.keepMaxDays {
			continue
		}
		if err := os.Remove(filepath.Join(rm.dir, fileName)); err != nil {
			return err
		}
	}

	return nil
}
