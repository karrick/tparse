// +build goparsetime

package tparse_test

import (
	"testing"
	"time"

	"github.com/etdub/goparsetime"
)

func BenchmarkParseNowMinusDurationGoParseTime(b *testing.B) {
	var t time.Time
	var err error

	for i := 0; i < b.N; i++ {
		t, err = goparsetime.Parsetime(benchmarkNowMinusDuration)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}
