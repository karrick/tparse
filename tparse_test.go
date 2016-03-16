package tparse

import (
	"testing"
	"time"
)

const benchmarkString = "now-2second"

func TestParseFloatingEpoch(t *testing.T) {
	actual, err := Parse("", "1445535988.5")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	expected := time.Unix(1445535988, fractionToNanos(0.5))
	if actual != expected {
		t.Errorf("Actual: %s; Expected: %s", actual, expected)
	}
}

func TestParseFloatingNegativeEpoch(t *testing.T) {
	_, err := Parse("", "-1445535988.5")
	if _, ok := err.(*time.ParseError); err == nil || !ok {
		t.Errorf("Actual: %#v; Expected: %s", err, "fixme")
	}
}

func TestParseNow(t *testing.T) {
	before := time.Now()
	actual, err := ParseNow("", "now")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowMinusMilliisecond(t *testing.T) {
	before := time.Now()
	time.Sleep(10 * time.Millisecond)
	actual, err := ParseNow("", "now-10ms")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowPlusMilliisecond(t *testing.T) {
	before := time.Now()
	actual, err := ParseNow("", "now+10ms")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	time.Sleep(10 * time.Millisecond)
	after := time.Now()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseLayout(t *testing.T) {
	actual, err := Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	expected := time.Unix(1136214245, 0)
	if !actual.Equal(expected) {
		t.Errorf("Actual: %d; Expected: %d", actual.Unix(), expected.Unix())
	}
}

func TestParseNowPlusDay(t *testing.T) {
	before := time.Now().UTC().AddDate(0, 0, 1).Add(time.Hour).Add(time.Minute)
	actual, err := ParseNow("", "now+1h1d1m")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now().UTC().AddDate(0, 0, 1).Add(time.Hour).Add(time.Minute)
	actual = actual.UTC()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowPlusAndMinus(t *testing.T) {
	before := time.Now().UTC().Add(time.Hour).AddDate(0, 0, -1).Add(time.Minute)
	actual, err := ParseNow("", "now+1h-1d+1m")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now().UTC().Add(time.Hour).AddDate(0, 0, -1).Add(time.Minute)
	actual = actual.UTC()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowMinusAndPlus(t *testing.T) {
	before := time.Now().UTC().Add(-time.Hour*12).AddDate(0, 0, 34).Add(-time.Minute * 56)
	actual, err := ParseNow("", "now-12h+34day-56m")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now().UTC().Add(-time.Hour*12).AddDate(0, 0, 34).Add(-time.Minute * 56)
	actual = actual.UTC()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseUsingMap(t *testing.T) {
	before := time.Now().UTC()
	dict := map[string]time.Time{
		"start": time.Now().UTC().AddDate(0, 0, -7),
	}
	after := time.Now().UTC()

	actual, err := ParseWithMap(time.ANSIC, "start+1week", dict)
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}

	actual = actual.UTC()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func BenchmarkParseNow(b *testing.B) {
	var t time.Time
	var err error

	for i := 0; i < b.N; i++ {
		t, err = ParseNow(time.ANSIC, benchmarkString)
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
		t, err = ParseWithMap(time.ANSIC, value, m)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = t
}

func TestParseNowMinusSecond(t *testing.T) {
	before := time.Now().UTC().Add(-2 * time.Second)
	actual, err := ParseNow("", "now-2second")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now().UTC().Add(-2 * time.Second)
	actual = actual.UTC()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}
