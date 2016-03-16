package tparse

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Parse will return the time value corresponding to the specified layout and value.  It also parses
// floating point and integer epoch values.
func Parse(layout, value string) (time.Time, error) {
	return ParseWithMap(layout, value, make(map[string]time.Time))
}

// ParseNow will return the time value corresponding to the specified layout and value.  It also
// parses floating point and integer epoch values.  It recognizes the special string `now` and
// replaces that with the time ParseNow is called.  This allows a suffix adding or subtracting
// various values from the base time.  For instance, ParseNow(time.ANSIC, "now+1d") will return a
// time corresponding to 24 hours from the moment the function is invoked.
//
// In addition to the duration abbreviations recognized by time.ParseDuration, it recognizes various
// tokens for days, weeks, months, and years.
//
//	package main
//
//	import (
//		"fmt"
//		"os"
//		"time"
//
//		tparse "gopkg.in/karrick/tparse.v2"
//	)
//
//	func main() {
//		actual, err := tparse.ParseNow(time.RFC3339, "now+1d3w4mo7y6h4m")
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "error: %s\n", err)
//			os.Exit(1)
//		}
//
//		fmt.Printf("time is: %s\n", actual)
//	}
func ParseNow(layout, value string) (time.Time, error) {
	m := map[string]time.Time{"now": time.Now()}
	return ParseWithMap(layout, value, m)
}

// ParseWithMap will return the time value corresponding to the specified layout and value.  It also
// parses floating point and integer epoch values.  It accepts a map of strings to time.Time values,
// and if the value string starts with one of the keys in the map, it replaces the string with the
// corresponding time.Time value.
//
//	package main
//
//	import (
//		"fmt"
//		"os"
//		"time"
//
//		tparse "gopkg.in/karrick/tparse.v2"
//	)
//
//	func main() {
//		m := make(map[string]time.Time)
//		m["start"] = start
//
//		end, err := tparse.ParseWithMap(time.RFC3339, "start+8h", m)
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "error: %s\n", err)
//			os.Exit(1)
//		}
//
//		fmt.Printf("start: %s; end: %s\n", start, end)
//	}
func ParseWithMap(layout, value string, dict map[string]time.Time) (time.Time, error) {
	if epoch, err := strconv.ParseFloat(value, 64); err == nil && epoch >= 0 {
		trunc := math.Trunc(epoch)
		nanos := fractionToNanos(epoch - trunc)
		return time.Unix(int64(trunc), int64(nanos)), nil
	}

	var matchKey []byte
	var matchTime time.Time
	// find longest matching key in dict
	for k, v := range dict {
		if strings.HasPrefix(value, k) && len(k) > len(matchKey) {
			matchKey = []byte(k)
			matchTime = v
		}
	}
	if len(matchKey) > 0 {
		return addDuration(matchTime, value[len(matchKey):])
	}
	return time.Parse(layout, value)
}

func fractionToNanos(fraction float64) int64 {
	return int64(fraction * float64(time.Second/time.Nanosecond))
}

// on err, returns epoch and error
func addDuration(base time.Time, value string) (time.Time, error) {
	if len(value) == 0 {
		return base, nil
	}
	var epoch time.Time
	var ty, tm, td int
	var tdur time.Duration
	var identifier, setComplete bool
	positive := true
	var iUnit, iNumber int
	var startNumberNextRune bool

	for i, rune := range value {
		if startNumberNextRune {
			iNumber = i
			startNumberNextRune = false
		}
		// [+-][0-9]+[^-+0-9]+
		if identifier {
			switch {
			case rune == '+', rune == '-':
				identifier = false
				setComplete = true
				startNumberNextRune = true
			case unicode.IsDigit(rune):
				identifier = false
				setComplete = true
			}
			if setComplete {
				if i > 0 {
					// we should have all we need for previous set
					y, m, d, dur, err := bar(value, positive, iNumber, iUnit, i)
					if err != nil {
						return epoch, err
					}
					ty += y
					tm += m
					td += d
					tdur += dur
					iNumber = i
				}
				setComplete = false
			}
			switch {
			case rune == '+':
				positive = true
			case rune == '-':
				positive = false
			}
		} else { // number
			switch {
			case rune == '+':
				positive = true
				startNumberNextRune = true
			case rune == '-':
				positive = false
				startNumberNextRune = true
			case unicode.IsDigit(rune):
				// nop
			default:
				identifier = true
				iUnit = i
			}
		}
	}

	if iNumber < iUnit && iUnit < len(value) {
		y, m, d, dur, err := bar(value, positive, iNumber, iUnit, len(value))
		if err != nil {
			return epoch, err
		}
		ty += y
		tm += m
		td += d
		tdur += dur
	} else {
		return epoch, fmt.Errorf("extra characters: %s", value[iNumber:])
	}
	return base.Add(tdur).AddDate(ty, tm, td), nil
}

func bar(value string, positive bool, iNumber, iUnit, i int) (int, int, int, time.Duration, error) {
	number := value[iNumber:iUnit]
	unit := value[iUnit:i]
	return calcDuration(positive, number, unit)
}

func calcDuration(positive bool, number, unit string) (int, int, int, time.Duration, error) {
	value, err := strconv.Atoi(number)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	var y, m, d int
	var duration time.Duration

	// NOTE: compare byte slices because some units, i.e. ms, are multi-rune
	switch unit {
	case "d", "day", "days":
		d = value
	case "w", "week", "weeks":
		d = 7 * value
	case "mo", "mon", "month", "months", "mth", "mn":
		m = value
	case "y", "year", "years":
		y = value
	case "sec", "second", "seconds":
		duration = time.Duration(value) * time.Second
	case "min", "minute", "minutes":
		duration = time.Duration(value) * time.Minute
	case "hr", "hour", "hours":
		duration = time.Duration(value) * time.Hour
	default:
		duration, err = time.ParseDuration(number + unit)
	}

	if !positive {
		y = -y
		m = -m
		d = -d
		duration = -duration
	}

	return y, m, d, duration, nil
}
