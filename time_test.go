package main

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	var (
		s    string
		ts   int64
		p, x time.Time
		err  error
	)

	// Local time
	now := time.Now()

	times := map[string]time.Time{
		"5m":     now.Add(time.Duration(time.Minute * 5)),
		"-0h":    now.Add(-time.Duration(time.Hour * 0)),
		"-48h5m": now.Add(-time.Duration(time.Hour*48 + time.Minute*5)),

		// UTC
		"2013-04-10":                     time.Date(2013, 4, 10, 0, 0, 0, 0, time.UTC),
		"April 4, 2013":                  time.Date(2013, 4, 4, 0, 0, 0, 0, time.UTC),
		"Apr 04, 2013":                   time.Date(2013, 4, 4, 0, 0, 0, 0, time.UTC),
		"47065363200000000":              time.Date(1492, 6, 11, 0, 0, 0, 0, time.UTC),
		"02-01-2006":                     time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		"02-01-2006 2:04 PM":             time.Date(2006, 1, 2, 14, 4, 0, 0, time.UTC),
		"02-01-2006 2:04 PM -0700":       time.Date(2006, 1, 2, 21, 4, 0, 0, time.UTC),
		"02-01-2006 2:04 PM -07:00":      time.Date(2006, 1, 2, 21, 4, 0, 0, time.UTC),
		"2 January 2006":                 time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		"2 January 2006 3:04 PM":         time.Date(2006, 1, 2, 15, 4, 0, 0, time.UTC),
		"2 January 2006 3:04 PM -0700":   time.Date(2006, 1, 2, 22, 4, 0, 0, time.UTC),
		"2 January 2006 3:04 PM -07:00":  time.Date(2006, 1, 2, 22, 4, 0, 0, time.UTC),
		"2006-01-02":                     time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		"2006-01-02 3:04 PM":             time.Date(2006, 1, 2, 15, 4, 0, 0, time.UTC),
		"2006-01-02 3:04 PM -0700":       time.Date(2006, 1, 2, 22, 4, 0, 0, time.UTC),
		"2006-01-02 3:04 PM -07:00":      time.Date(2006, 1, 2, 22, 4, 0, 0, time.UTC),
		"January 2, 2006":                time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		"January 2, 2006 3:04 PM":        time.Date(2006, 1, 2, 15, 4, 0, 0, time.UTC),
		"January 2, 2006 3:04 PM -0700":  time.Date(2006, 1, 2, 22, 4, 0, 0, time.UTC),
		"January 2, 2006 3:04 PM -07:00": time.Date(2006, 1, 2, 22, 4, 0, 0, time.UTC),
		"Jan 2, 2006":                    time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		"Jan 2, 2006, 3:04 PM":           time.Date(2006, 1, 2, 15, 4, 0, 0, time.UTC),
		"Jan 2, 2006 3:04 PM -0700":      time.Date(2006, 1, 2, 22, 4, 0, 0, time.UTC),
		"Jan 2, 2006 3:04 PM -07:00":     time.Date(2006, 1, 2, 22, 4, 0, 0, time.UTC),
	}

	// Duration to truncate for comparison.
	td := time.Second

	for s, x = range times {
		ts, err = ParseTimeString(s)
		p = time.Unix(ts, 0)

		if err != nil {
			t.Errorf("time: failed to parse %s as time", s)
		} else {
			p = p.Truncate(td)
			x = p.Truncate(td)

			if !p.Equal(x) {
				t.Errorf("time: expected %s, got %s", x, p)
			}
		}
	}
}

func BenchmarkParseTimeString__Time(b *testing.B) {
	t := "April 4, 2013"

	for i := 0; i < b.N; i++ {
		ParseTimeString(t)
	}
}

func BenchmarkParseTimeString__Duration(b *testing.B) {
	t := "-48h32m"

	for i := 0; i < b.N; i++ {
		ParseTimeString(t)
	}
}

func BenchmarkParseTimeString__Timestamp(b *testing.B) {
	t := "63592300800000000"

	for i := 0; i < b.N; i++ {
		ParseTimeString(t)
	}
}
