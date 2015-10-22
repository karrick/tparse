package tparse

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Parse will return the time corresponding to the layout and value.  It also parses floating point
// epoch values, and values of "now", "now+DURATION", and "now-DURATION". It uses time.ParseDuration
// to parse the specified DURATION.
func Parse(layout, value string) (time.Time, error) {
	if epoch, err := strconv.ParseFloat(value, 64); err == nil && epoch >= 0 {
		trunc := math.Trunc(epoch)
		nanos := fractionToNanos(epoch - trunc)
		return time.Unix(int64(trunc), int64(nanos)), nil
	}
	var t time.Time
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
			duration, err = time.ParseDuration(value[4:])
			if err != nil {
				return t, err
			}
		}
		return time.Now().Add(time.Duration(int(duration) * direction)), nil
	}
	return time.Parse(layout, value)
}

func fractionToNanos(fraction float64) int64 {
	return int64(fraction * float64(time.Second/time.Nanosecond))
}
