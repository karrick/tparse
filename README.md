# tparse

Parse will return the time corresponding to the layout and value.  It
also parses floating point epoch values, and values of "now",
"now+DURATION", and "now-DURATION".

In addition to the duration abbreviations recognized by
time.ParseDuration, it recognizes various tokens for days, weeks,
months, and years.

## Documentation

In addition to this handy README.md file, documentation is available
in godoc format at
[![GoDoc](https://godoc.org/github.com/karrick/tparse?status.svg)](https://godoc.org/github.com/karrick/tparse).

## Examples

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
    	actual, err := tparse.ParseNow(time.RFC3339, "now+1d3w4mo7y6h4m")
    	if err != nil {
    		fmt.Fprintf(os.Stderr, "error: %s\n", err)
    		os.Exit(1)
        }
    	fmt.Printf("time is: %s\n", actual)
    }
```

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
    	start, err := tparse.ParseNow(time.RFC3339, "now")
    	if err != nil {
    		fmt.Fprintf(os.Stderr, "error: %s\n", err)
    		os.Exit(1)
    	}

        m := make(map[string]time.Time)
        m["start"] = start

    	end, err := tparse.ParseWithMap(time.RFC3339, "start+8h", m)
    	if err != nil {
    		fmt.Fprintf(os.Stderr, "error: %s\n", err)
    		os.Exit(1)
    	}

    	fmt.Printf("start: %s; end: %s\n", start, end)
    }
```
