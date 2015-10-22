package tparse

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Parse will return the time corresponding to the layout and value.  It also parses floating point
// epoch values, and values of "now", "now+DURATION", and "now-DURATION".
//
// In addition to the duration abbreviations recognized by time.ParseDuration, it recognizes the
// following abbreviations:
//
//   year: y
//   month: mo, mon, mth, mn
//   week: w
//   day: d
func Parse(layout, value string) (time.Time, error) {
	if epoch, err := strconv.ParseFloat(value, 64); err == nil && epoch >= 0 {
		trunc := math.Trunc(epoch)
		nanos := fractionToNanos(epoch - trunc)
		return time.Unix(int64(trunc), int64(nanos)), nil
	}
	var t time.Time
	var y, m, d int
	if strings.HasPrefix(value, "now") {
		var duration time.Duration
		var direction = 1
		var err error

		if len(value) > 3 {
			switch value[3] {
			case '+':
				// no-op
			case '-':
				direction = -1
			default:
				return t, fmt.Errorf("can only subtract or add to now")
			}
			var nv string
			y, m, d, nv, err = ymd(value[4:])
			if err != nil {
				return t, err
			}
			if len(nv) > 0 {
				duration, err = time.ParseDuration(nv)
				if err != nil {
					return t, err
				}
			}
		}
		if direction < 0 {
			y = -y
			m = -m
			d = -d
		}
		return time.Now().Add(time.Duration(int(duration)*direction)).AddDate(y, m, d), nil
	}
	return time.Parse(layout, value)
}

func fractionToNanos(fraction float64) int64 {
	return int64(fraction * float64(time.Second/time.Nanosecond))
}

func ymd(value string) (int, int, int, string, error) {
	// alternating numbers and strings
	var y, m, d int
	var accum int     // accumulates digits
	var unit []byte   // accumulates units
	var unproc []byte // accumulate unprocessed durations to return

	unitComplete := func() {
		// NOTE: compare byte slices because some units, i.e. ms, are multi-rune
		if bytes.Equal(unit, []byte{'d'}) {
			d += accum
		} else if bytes.Equal(unit, []byte{'w'}) {
			d += 7 * accum
		} else if bytes.Equal(unit, []byte{'m', 'o'}) || bytes.Equal(unit, []byte{'m', 'o', 'n'}) || bytes.Equal(unit, []byte{'m', 't', 'h'}) || bytes.Equal(unit, []byte{'m', 'n'}) {
			m += accum
		} else if bytes.Equal(unit, []byte{'y'}) {
			y += accum
		} else {
			unproc = append(append(unproc, strconv.Itoa(accum)...), unit...)
		}
	}

	expectDigit := true
	for _, rune := range value {
		if unicode.IsDigit(rune) {
			if expectDigit {
				accum = accum*10 + int(rune-'0')
			} else {
				unitComplete()
				unit = unit[:0]
				accum = int(rune - '0')
			}
			continue
		}
		unit = append(unit, string(rune)...)
		expectDigit = false
	}
	if len(unit) > 0 {
		unitComplete()
		accum = 0
		unit = unit[:0]
	}
	return y, m, d, string(unproc), nil
}
