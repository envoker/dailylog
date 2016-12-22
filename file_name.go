package dailylog

import (
	"errors"
	"strings"

	"github.com/toelsiba/date"
)

// fileName format: <prefix> | <date> | <ext>

type paramsFileName struct {
	prefix     string
	dateFormat string
	ext        string
}

const defaultDateFormat = "2006-01-02"

func checkDateFormat(dateFormat string) error {
	d := date.CurrentDate()
	s := d.Format(dateFormat)
	q, err := date.Parse(dateFormat, s)
	if (err != nil) || !d.Equal(q) {
		return errors.New("dailylog: invalid date format")
	}
	return nil
}

func compileFileName(p paramsFileName, d date.Date) (fileName string) {
	return p.prefix + d.Format(p.dateFormat) + p.ext
}

var (
	errHasNotPrefix = errors.New("dailylog: has not prefix")
	errHasNotExt    = errors.New("dailylog: has not ext")
)

func parseDate(p paramsFileName, fileName string) (date.Date, error) {

	s := fileName

	if len(p.prefix) > 0 {
		if !strings.HasPrefix(fileName, p.prefix) {
			return date.Date{}, errHasNotPrefix
		}
		s = s[len(p.prefix):]
	}

	if len(p.ext) > 0 {
		if !strings.HasSuffix(fileName, p.ext) {
			return date.Date{}, errHasNotExt
		}
		s = s[:len(s)-len(p.ext)]
	}

	return date.Parse(p.dateFormat, s)
}
