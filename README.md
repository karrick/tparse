# tparse

`Parse` will return the time corresponding to the layout and value.
It also parses floating point epoch values, and values of "now",
"now+DURATION", and "now-DURATION".

In addition to the duration abbreviations recognized by
`time.ParseDuration`, tparse recognizes various tokens for days,
weeks, months, and years, as well as tokens for seconds, minutes, and
hours.

Like `time.ParseDuration`, it accepts multiple fractional scalars, so
"now+1.5days-3.21hours" is evaluated properly.

## Documentation

In addition to this handy README.md file, documentation is available
in godoc format at
[![GoDoc](https://godoc.org/github.com/karrick/tparse?status.svg)](https://godoc.org/github.com/karrick/tparse).

## Examples

### ParseNow

`ParseNow` can parse time values that are relative to the current
time, by specifying a string starting with "now", a '+' or '-' byte,
followed by a time duration.

```Go
    package main

    import (
        "fmt"
        "os"
        "time"
        "github.com/karrick/tparse"
    )

    func main() {
        actual, err := tparse.ParseNow(time.RFC3339, "now+1d-3w4mo+7y6h4m")
        if err != nil {
            fmt.Fprintf(os.Stderr, "error: %s\n", err)
            os.Exit(1)
        }
        fmt.Printf("time is: %s\n", actual)
    }
```

### ParseWithMap

`ParseWithMap` can parse time values that use a base time other than "now".

```Go
    package main

    import (
        "fmt"
        "os"
        "time"
        "github.com/karrick/tparse"
    )

    func main() {
        m := make(map[string]time.Time)
        m["end"] = time.Now()

        start, err := tparse.ParseWithMap(time.RFC3339, "end-12h", m)
        if err != nil {
            fmt.Fprintf(os.Stderr, "error: %s\n", err)
            os.Exit(1)
        }

        fmt.Printf("start: %s; end: %s\n", start, end)
    }
```

### AddDuration

`AddDuration` is used to compute the value of a duration string and
add it to a known time. This function is used by the other library
functions to parse all duration strings.

The following tokens may be used to specify the respective unit of
time:

 * Nanosecond: ns
 * Microsecond: us, µs (U+00B5 = micro symbol), μs (U+03BC = Greek letter mu)
 * Millisecond: ms
 * Second: s, sec, second, seconds
 * Minute: m, min, minute, minutes
 * Hour: h, hr, hour, hours
 * Day: d, day, days
 * Week: w, wk, week, weeks
 * Month: mo, mon, month, months
 * Year: y, yr, year, years

```Go
    package main

    import (
        "fmt"
        "os"
        "time"

        "github.com/karrick/tparse"
    )

    func main() {
        now := time.Now()
        another, err := tparse.AddDuration(now, "+1d3w4mo-7y6h4m")
        if err != nil {
            fmt.Fprintf(os.Stderr, "error: %s\n", err)
            os.Exit(1)
        }

        fmt.Printf("time is: %s\n", another)
    }
```

### AbsoluteDuration

When you would rather have the `time.Duration` representation of a
duration string, there is a function for that, but with a
caveat.

First, not every month has 30 days, and therefore Go does not have a
`time.Duration` type constant to represent one month. 

When I add one month to February 3, do I get March 3 or March 4?
Depends on what the year is and whether or not that year is a leap
year.

Is one month always 30 days? Is one month 31 days, or 28, or 29? I did
not want to have to answer this question, so I defaulted to saying the
length of one month depends on which month and year, and I allowed the
Go standard library to add duration concretely to a given moment in
time.

Consider the below two examples of calling `AbsoluteDuration` with the
same duration string, but different base times.

```Go
func ExampleAbsoluteDuration() {
    t1 := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

    d1, err := AbsoluteDuration(t1, "1.5month")
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(d1)

    t2 := time.Date(2020, time.February, 10, 23, 0, 0, 0, time.UTC)

    d2, err := AbsoluteDuration(t2, "1.5month")
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(d2)
    // Output:
    // 1080h0m0s
    // 1056h0m0s
}
```

## Benchmark against goparsetime

```Bash
GO111MODULE=on go test -bench=. -tags goparsetime
```
