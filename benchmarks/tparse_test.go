package benchmarks

import (
	"testing"
	"time"

	"github.com/karrick/tparse"
)

func BenchmarkAddDuration(b *testing.B) {
	var err error
	var t time.Time
	epoch := time.Now().UTC()

	for i := 0; i < b.N; i++ {
		t, err = tparse.AddDuration(epoch, benchmarkDuration)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}

func BenchmarkParseDurationPseudoStandardLibrary(b *testing.B) {
	var d time.Duration
	var err error

	for i := 0; i < b.N; i++ {
		d, err = time.ParseDuration(benchmarkDuration)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = d
}

func BenchmarkAddDurationStandardLibrary(b *testing.B) {
	var d time.Duration
	var err error
	var t time.Time
	epoch := time.Now().UTC()

	for i := 0; i < b.N; i++ {
		d, err = time.ParseDuration(benchmarkDuration)
		if err != nil {
			b.Fatal(err)
		}
		t = epoch.Add(d)
	}
	_ = t
}

//

func BenchmarkParseNowMinusDuration(b *testing.B) {
	var t time.Time
	var err error

	for i := 0; i < b.N; i++ {
		t, err = tparse.ParseNow("", benchmarkNowMinusDuration)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}

func BenchmarkParseWithMapEpoch(b *testing.B) {
	var t time.Time
	var err error
	value := "1458179403.12345"

	for i := 0; i < b.N; i++ {
		t, err = tparse.ParseWithMap(time.ANSIC, value, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}

func BenchmarkParseWithMapKeyedValue(b *testing.B) {
	var t time.Time
	var err error
	value := "end"

	m := make(map[string]time.Time)
	m["end"] = time.Now()

	for i := 0; i < b.N; i++ {
		t, err = tparse.ParseWithMap(time.ANSIC, value, m)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}

func BenchmarkParseWithMapKeyedValueAndDuration(b *testing.B) {
	var t time.Time
	var err error
	value := "end+1hr"

	m := make(map[string]time.Time)
	m["end"] = time.Now()

	for i := 0; i < b.N; i++ {
		t, err = tparse.ParseWithMap(time.ANSIC, value, m)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}

//

func BenchmarkParseRFC3339(b *testing.B) {
	var t time.Time
	var err error

	for i := 0; i < b.N; i++ {
		t, err = tparse.Parse(time.RFC3339, rfc3339)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}

func BenchmarkParseRFC3339StandardLibrary(b *testing.B) {
	var t time.Time
	var err error

	for i := 0; i < b.N; i++ {
		t, err = time.Parse(time.RFC3339, rfc3339)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}

func BenchmarkParseNow(b *testing.B) {
	var t time.Time
	var err error
	value := "now-5s"

	for i := 0; i < b.N; i++ {
		t, err = tparse.ParseNow(time.ANSIC, value)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}

func BenchmarkParseUsingMap(b *testing.B) {
	var t time.Time
	var err error
	value := "end-1mo"

	m := make(map[string]time.Time)
	m["end"] = time.Now()

	for i := 0; i < b.N; i++ {
		t, err = tparse.ParseWithMap(time.ANSIC, value, m)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}
