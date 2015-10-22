# tparse

Parse will return the time corresponding to the layout and value.  It also parses floating point
epoch values, and values of "now", "now+DURATION", and "now-DURATION".

In addition to the duration abbreviations recognized by time.ParseDuration, it recognizes the
following abbreviations:

  year: y
  month: mo, mon, mth, mn
  week: w
  day: d

## Example

```Go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/karrick/tparse"
)

func main() {
	actual, err := tparse.Parse(time.RFC3339, "now+1d3w4mo7y6h4m")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("time is: %s\n", actual)
}
```
