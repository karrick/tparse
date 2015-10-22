# tparse

Parse will return the time corresponding to the layout and value.  It
also parses floating point epoch values, and values of "now",
"now+DURATION", and "now-DURATION". It uses time.ParseDuration to
parse the specified DURATION.

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
	actual, err := tparse.Parse(time.RFC3339, "now+10ms")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("time is: %s\n", actual)
}
```
