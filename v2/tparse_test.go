package tparse_test

import (
	"testing"
	"time"

	tparse "github.com/karrick/tparse/v2"
)

const rfc3339 = "2006-01-02T15:04:05Z"

// AddDuration

func TestAddDurationRejectsSignWithoutDigits(t *testing.T) {
	t.Run("negative", func(t *testing.T) {
		_, err := tparse.AddDuration(time.Now(), "-")
		if err == nil {
			t.Errorf("(GOT): %v; (WNT): %v", err, "cannot parse sign without digits")
		}
	})
	t.Run("positive", func(t *testing.T) {
		_, err := tparse.AddDuration(time.Now(), "+")
		if err == nil {
			t.Errorf("(GOT): %v; (WNT): %v", err, "cannot parse sign without digits")
		}
	})
}

func TestAddDurationPositiveFractionalYear(t *testing.T) {
	start, err := tparse.Parse(time.RFC3339, "2003-07-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := tparse.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := tparse.AddDuration(start, "+2.5years")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}

	if actual != expected {
		t.Errorf("Actual: %s; Expected: %s", actual, expected)
	}
}

func TestAddDurationNegativeFractionalYear(t *testing.T) {
	start, err := tparse.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := tparse.Parse(time.RFC3339, "2003-07-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := tparse.AddDuration(start, "-2.5years")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}

	if actual != expected {
		t.Errorf("Actual: %s; Expected: %s", actual, expected)
	}
}

func TestAddDurationPositiveFractionalMonth(t *testing.T) {
	start, err := tparse.Parse(time.RFC3339, "2003-06-01T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := tparse.Parse(time.RFC3339, "2003-08-16T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := tparse.AddDuration(start, "+2.5months")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}

	if actual != expected {
		t.Errorf("Actual: %s; Expected: %s", actual, expected)
	}
}

func TestAddDurationNegativeFractionalMonth(t *testing.T) {
	start, err := tparse.Parse(time.RFC3339, "2003-08-16T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := tparse.Parse(time.RFC3339, "2003-06-01T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := tparse.AddDuration(start, "-2.5months")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}

	if actual != expected {
		t.Errorf("Actual: %s; Expected: %s", actual, expected)
	}
}

func TestAddDurationPositiveFractionalDay(t *testing.T) {
	start, err := tparse.Parse(time.RFC3339, "2003-06-01T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := tparse.Parse(time.RFC3339, "2003-06-04T03:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := tparse.AddDuration(start, "+2.5days")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}

	if actual != expected {
		t.Errorf("Actual: %s; Expected: %s", actual, expected)
	}
}

func TestAddDurationNegativeFractionalDay(t *testing.T) {
	start, err := tparse.Parse(time.RFC3339, "2003-06-04T03:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := tparse.Parse(time.RFC3339, "2003-06-01T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := tparse.AddDuration(start, "-2.5days")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}

	if actual != expected {
		t.Errorf("Actual: %s; Expected: %s", actual, expected)
	}
}

// ParseWithMap

func TestParseWithMapFloatingEpochPositive(t *testing.T) {
	actual, err := tparse.ParseWithMap("", "1445535988.5", nil)
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}

	nanos := int64(0.5 * float64(time.Second/time.Nanosecond))

	expected := time.Unix(1445535988, nanos)
	if actual != expected {
		t.Errorf("Actual: %s; Expected: %s", actual, expected)
	}
}

func TestParseWithMapFloatingEpochNegative(t *testing.T) {
	_, err := tparse.ParseWithMap("", "-1445535988.5", nil)
	if _, ok := err.(*time.ParseError); err == nil || !ok {
		t.Errorf("Actual: %#v; Expected: %s", err, "negative floating point not allowed")
	}
}

func TestParseWithMap(t *testing.T) {
	before := time.Now().UTC()
	dict := map[string]time.Time{
		"start": time.Now().UTC().AddDate(0, 0, -7),
	}
	after := time.Now().UTC()

	actual, err := tparse.ParseWithMap(time.ANSIC, "start+1week", dict)
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}

	actual = actual.UTC()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

// ParseNow

func TestParseNow(t *testing.T) {
	before := time.Now()
	actual, err := tparse.ParseNow("", "now")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowMinusSecond(t *testing.T) {
	before := time.Now().UTC().Add(-2 * time.Second)
	actual, err := tparse.ParseNow("", "now-2second")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now().UTC().Add(-2 * time.Second)
	actual = actual.UTC()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowMinusMillisecond(t *testing.T) {
	before := time.Now()
	time.Sleep(10 * time.Millisecond)
	actual, err := tparse.ParseNow("", "now-10ms")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowPlusMillisecond(t *testing.T) {
	before := time.Now()
	actual, err := tparse.ParseNow("", "now+10ms")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	time.Sleep(10 * time.Millisecond)
	after := time.Now()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowPlusQuarterDay(t *testing.T) {
	before := time.Now().UTC().Add(6 * time.Hour)
	actual, err := tparse.ParseNow("", "now+0.25day")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now().UTC().Add(6 * time.Hour)
	actual = actual.UTC()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

func TestParseNowPlusDay(t *testing.T) {
	before := time.Now().UTC().AddDate(0, 0, 1).Add(time.Hour).Add(time.Minute)
	actual, err := tparse.ParseNow("", "now+1h1d1m")
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
	actual, err := tparse.ParseNow("", "now+1h-1d+1m")
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
	actual, err := tparse.ParseNow("", "now-12hour+34day-56min")
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	after := time.Now().UTC().Add(-time.Hour*12).AddDate(0, 0, 34).Add(-time.Minute * 56)
	actual = actual.UTC()
	if before.After(actual) || actual.After(after) {
		t.Errorf("Actual: %s; Expected between: %s and %s", actual, before, after)
	}
}

// Parse

func TestParseLayout(t *testing.T) {
	actual, err := tparse.Parse(time.RFC3339, rfc3339)
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	expected := time.Unix(1136214245, 0)
	if !actual.Equal(expected) {
		t.Errorf("Actual: %d; Expected: %d", actual.Unix(), expected.Unix())
	}
}
