package tparse

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
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
	if strings.HasPrefix(value, "now") {
		return AddDuration(time.Now(), value[3:])
	}
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
		return AddDuration(matchTime, value[len(matchKey):])
	}
	return time.Parse(layout, value)
}

func fractionToNanos(fraction float64) int64 {
	return int64(fraction * float64(time.Second/time.Nanosecond))
}

var unitMap = map[string]int64{
	"ns":      int64(time.Nanosecond),
	"us":      int64(time.Microsecond),
	"µs":      int64(time.Microsecond), // U+00B5 = micro symbol
	"μs":      int64(time.Microsecond), // U+03BC = Greek letter mu
	"ms":      int64(time.Millisecond),
	"s":       int64(time.Second),
	"sec":     int64(time.Second),
	"second":  int64(time.Second),
	"seconds": int64(time.Second),
	"m":       int64(time.Minute),
	"min":     int64(time.Minute),
	"minute":  int64(time.Minute),
	"minutes": int64(time.Minute),
	"h":       int64(time.Hour),
	"hr":      int64(time.Hour),
	"hour":    int64(time.Hour),
	"hours":   int64(time.Hour),
	"d":       int64(time.Hour * 24),
	"day":     int64(time.Hour * 24),
	"days":    int64(time.Hour * 24),
	"w":       int64(time.Hour * 24 * 7),
	"week":    int64(time.Hour * 24 * 7),
	"weeks":   int64(time.Hour * 24 * 7),
}

// AddDuration parses the duration string, and adds it to the base time. On error, it returns the
// base time and the error.
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
//              now := time.Now()
//		another, err := tparse.AddDuration(now, "now+1d3w4mo7y6h4m")
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "error: %s\n", err)
//			os.Exit(1)
//		}
//
//		fmt.Printf("time is: %s\n", another)
//	}
func AddDuration(base time.Time, s string) (time.Time, error) {
	if len(s) == 0 {
		return base, nil
	}
	var totalYears, totalMonths int64
	var totalDuration int64
	var number int64
	var isNegative bool

	for s != "" {
		// consume possible sign
		if s[0] == '+' {
			isNegative = false
			s = s[1:]
		} else if s[0] == '-' {
			isNegative = true
			s = s[1:]
		}
		// consume digits
		for ; s[0] >= '0' && s[0] <= '9'; s = s[1:] {
			number *= 10
			number += int64(s[0] - '0')
		}
		if isNegative {
			number *= -1
		}
		// find end of unit
		var i int
		for ; i < len(s) && s[i] != '+' && s[i] != '-' && (s[i] < '0' || s[i] > '9'); i++ {
			// identifier bytes: no-op
		}
		unit := s[:i]
		s = s[i:]
		// fmt.Printf("number: %d; unit: %q\n", number, unit)
		if dur, ok := unitMap[unit]; ok {
			totalDuration += number * dur
		} else {
			switch unit {
			case "mo", "mon", "month", "months", "mth", "mn":
				totalMonths += number
			case "y", "year", "years":
				totalYears += number
			default:
				return base, fmt.Errorf("unknown unit in duration: %q", unit)
			}
		}

		number = 0
	}
	return base.AddDate(int(totalYears), int(totalMonths), 0).Add(time.Duration(totalDuration)), nil
}
