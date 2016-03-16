package tparse

import (
	"testing"
	"time"

	"github.com/etdub/goparsetime"
)

func BenchmarkGoParseTime(b *testing.B) {
	var t time.Time
	var err error

	for i := 0; i < b.N; i++ {
		t, err = goparsetime.Parsetime(benchmarkString)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}
