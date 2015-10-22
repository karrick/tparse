package tparse

import (
	"testing"
	"time"
)

func TestParseFloatingEpoch(t *testing.T) {
	actual, err := Parse("", "1445535988.5")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	expected := time.Unix(1445535988, fractionToNanos(0.5))
	if actual != expected {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
}

func TestParseFloatingNegativeEpoch(t *testing.T) {
	_, err := Parse("", "-1445535988.5")
	if _, ok := err.(*time.ParseError); err == nil || !ok {
		t.Errorf("Actual: %#v; Expected: %#v", err, "fixme")
	}
}

func TestParseNow(t *testing.T) {
	before := time.Now()
	actual, err := Parse("", "now")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %#v; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowMinusMilliisecond(t *testing.T) {
	before := time.Now()
	time.Sleep(10 * time.Millisecond)
	actual, err := Parse("", "now-10ms")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %#v; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowPlusMilliisecond(t *testing.T) {
	before := time.Now()
	actual, err := Parse("", "now+10ms")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	time.Sleep(10 * time.Millisecond)
	after := time.Now()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %#v; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseLayout(t *testing.T) {
	actual, err := Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	expected := time.Unix(1136214245, 0)
	if !actual.Equal(expected) {
		t.Errorf("Actual: %#v; Expected: %#v", actual.Unix(), expected.Unix())
	}
}

func TestParseNowPlusDay(t *testing.T) {
	before := time.Now().UTC().AddDate(0, 0, 1).Add(time.Hour).Add(time.Minute)
	actual, err := Parse("", "now+1h1d1m")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now().UTC().AddDate(0, 0, 1).Add(time.Hour).Add(time.Minute)
	actual = actual.UTC()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}
