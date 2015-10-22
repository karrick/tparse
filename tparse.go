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
	return ParseDict(layout, value, make(map[string]time.Time))
}

// ParseDict parses time values exactly like Parse, but allows a customizable dictionary of base
// time names and their respective time values.
func ParseDict(layout, value string, dict map[string]time.Time) (time.Time, error) {
	if epoch, err := strconv.ParseFloat(value, 64); err == nil && epoch >= 0 {
		trunc := math.Trunc(epoch)
		nanos := fractionToNanos(epoch - trunc)
		return time.Unix(int64(trunc), int64(nanos)), nil
	}
	var base time.Time
	var y, m, d int
	var duration time.Duration
	var direction = 1
	var err error

	if _, ok := dict["now"]; !ok {
		dict["now"] = time.Now()
	}

	for k, v := range dict {
		if strings.HasPrefix(value, k) {
			base = v
			if len(value) > len(k) {
				// maybe has +, -
				switch dir := value[len(k)]; dir {
				case '+':
					// no-op
				case '-':
					direction = -1
				default:
					return base, fmt.Errorf("expected '+' or '-': %q", dir)
				}
				var nv string
				y, m, d, nv = ymd(value[len(k)+1:])
				if len(nv) > 0 {
					duration, err = time.ParseDuration(nv)
					if err != nil {
						return base, err
					}
				}
			}
			if direction < 0 {
				y = -y
				m = -m
				d = -d
			}
			return base.Add(time.Duration(int(duration)*direction)).AddDate(y, m, d), nil
		}
	}
	return time.Parse(layout, value)
}

func fractionToNanos(fraction float64) int64 {
	return int64(fraction * float64(time.Second/time.Nanosecond))
}

func ymd(value string) (int, int, int, string) {
	// alternating numbers and strings
	var y, m, d int
	var accum int     // accumulates digits
	var unit []byte   // accumulates units
	var unproc []byte // accumulate unprocessed durations to return

	unitComplete := func() {
		// NOTE: compare byte slices because some units, i.e. ms, are multi-rune
		if bytes.Equal(unit, []byte{'d'}) || bytes.Equal(unit, []byte{'d', 'a', 'y'}) || bytes.Equal(unit, []byte{'d', 'a', 'y', 's'}) {
			d += accum
		} else if bytes.Equal(unit, []byte{'w'}) || bytes.Equal(unit, []byte{'w', 'e', 'e', 'k'}) || bytes.Equal(unit, []byte{'w', 'e', 'e', 'k', 's'}) {
			d += 7 * accum
		} else if bytes.Equal(unit, []byte{'m', 'o'}) || bytes.Equal(unit, []byte{'m', 'o', 'n'}) || bytes.Equal(unit, []byte{'m', 'o', 'n', 't', 'h'}) || bytes.Equal(unit, []byte{'m', 'o', 'n', 't', 'h', 's'}) || bytes.Equal(unit, []byte{'m', 't', 'h'}) || bytes.Equal(unit, []byte{'m', 'n'}) {
			m += accum
		} else if bytes.Equal(unit, []byte{'y'}) || bytes.Equal(unit, []byte{'y', 'e', 'a', 'r'}) || bytes.Equal(unit, []byte{'y', 'e', 'a', 'r', 's'}) {
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
	// log.Printf("y: %d; m: %d; d: %d; nv: %q", y, m, d, unproc)
	return y, m, d, string(unproc)
}
