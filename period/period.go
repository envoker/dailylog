package period

import (
	"encoding"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// duration
type Period int64 // seconds

const (
	Second Period = 1
	Minute        = 60 * Second
	Hour          = 60 * Minute
	Day           = 24 * Hour
)

// units or designators: s, m, h, d

func PeriodByDuration(d time.Duration) Period {
	return Period(d / time.Second)
}

func (p Period) Duration() time.Duration {
	return time.Second * time.Duration(p)
}

// func (p Period) Days() int64 {
// 	return int64(p / Day)
// }

// func (p Period) Hours() int64 {
// 	return int64(p / Hour)
// }

func (p Period) String() string {

	v := int64(p)

	var b strings.Builder

	if v < 0 {
		b.WriteByte('-') // write minus sign
		v = -v
	}

	days, v := quoRemInt64(v, int64(Day))
	if days != 0 {
		b.WriteString(strconv.FormatInt(days, 10))
		b.WriteByte('d')
	}

	hours, v := quoRemInt64(v, int64(Hour))
	if hours != 0 {
		b.WriteString(strconv.FormatInt(hours, 10))
		b.WriteByte('h')
	}

	minutes, v := quoRemInt64(v, int64(Minute))
	if minutes != 0 {
		b.WriteString(strconv.FormatInt(minutes, 10))
		b.WriteByte('m')
	}

	seconds := v
	if (p == 0) || (seconds != 0) {
		b.WriteString(strconv.FormatInt(seconds, 10))
		b.WriteByte('s')
	}

	return b.String()
}

func quoRemInt64(a, b int64) (quo, rem int64) {
	quo = a / b
	rem = a % b
	return
}

func errorParsePeriod(s string) error {
	return errors.New("ParsePeriod: " + s)
}

func ParsePeriod(s string) (Period, error) {

	var p Period

	has := struct {
		seconds bool
		minutes bool
		hours   bool
		days    bool
	}{}

	data := []byte(s)
	pos := 0

	var negative bool
	if (len(data[pos:]) > 0) && (data[pos] == '-') {
		negative = true
		pos++
	}

	for len(data[pos:]) > 0 {

		// read number
		number, k := readNumber(data[pos:])
		if k == 0 {
			id := data[pos] // id - invalid digit
			return 0, fmt.Errorf("ParsePeriod: it has invalid digit (%q, 0x%02X) of number by index %d", id, id, pos)
		}
		pos += k

		// read designator
		if len(data[pos:]) == 0 {
			return 0, errorParsePeriod("it has not last designator")
		}
		designator := data[pos]

		switch designator {
		case 's':
			if has.seconds {
				return 0, errorParsePeriod("it has duplicate of seconds")
			} else {
				p += Period(number) * Second
				has.seconds = true
			}

		case 'm':
			if has.minutes {
				return 0, errorParsePeriod("it has duplicate of minutes")
			} else {
				p += Period(number) * Minute
				has.minutes = true
			}

		case 'h':
			if has.hours {
				return 0, errorParsePeriod("it has duplicate of hours")
			} else {
				p += Period(number) * Hour
				has.hours = true
			}

		case 'd':
			if has.days {
				return 0, errorParsePeriod("it has duplicate of days")
			} else {
				p += Period(number) * Day
				has.days = true
			}

		default:
			return 0, fmt.Errorf("ParsePeriod: it has invalid designator (%q, 0x%02X) by index %d", designator, designator, pos)
		}

		pos++ // skip designator
	}

	if !has.seconds && !has.minutes && !has.hours && !has.days {
		return 0, errorParsePeriod("it has not any designators")
	}

	if negative {
		p = -p
	}

	return p, nil
}

func byteIsDigit(b byte) bool {
	return ('0' <= b) && (b <= '9')
}

func byteToDigit(b byte) (digit int, ok bool) {
	if ('0' <= b) && (b <= '9') {
		digit = int(b - '0')
		return digit, true
	}
	return 0, false
}

func readNumber(bs []byte) (number int, k int) {
	for k < len(bs) {
		d, ok := byteToDigit(bs[k])
		if !ok {
			break
		}
		number = number*10 + d
		k++
	}
	return number, k
}

func _() {
	var p Period
	var (
		_ encoding.TextMarshaler   = p
		_ encoding.TextUnmarshaler = &p
	)
}

func (p Period) MarshalText() (text []byte, err error) {
	s := p.String()
	return []byte(s), nil
}

func (p *Period) UnmarshalText(text []byte) error {
	q, err := ParsePeriod(string(text))
	if err != nil {
		return err
	}
	*p = q
	return nil
}
