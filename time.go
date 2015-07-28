package main

import (
	"fmt"
	"strconv"
	"time"
)

// TimeLayouts is a list of time layouts that are used when parsing
// a time string.
var timeLayouts = []string{
	"02-01-2006",
	"02-01-2006 3:04 PM",
	"02-01-2006 3:04 PM -0700",
	"02-01-2006 3:04 PM -07:00",
	"_2 January 2006",
	"_2 January 2006 3:04 PM",
	"_2 January 2006 3:04 PM -0700",
	"_2 January 2006 3:04 PM -07:00",
	"2006-01-02",
	"2006-01-02 3:04 PM",
	"2006-01-02 3:04 PM -0700",
	"2006-01-02 3:04 PM -07:00",
	time.RFC1123,
	time.RFC1123Z,
	time.RFC822,
	time.RFC822Z,
	"January _2, 2006",
	"January _2, 2006 3:04 PM",
	"January _2, 2006 3:04 PM -0700",
	"January _2, 2006 3:04 PM -07:00",
	"Jan _2, 2006",
	"Jan _2, 2006, 3:04 PM",
	"Jan _2, 2006 3:04 PM -0700",
	"Jan _2, 2006 3:04 PM -07:00",
	time.RFC3339,
	time.ANSIC,
}

// ParseTimeString parses a string into a Unix timestamp. The string may represent an
// absolute time, duration relative to the current time, or a millisecond resolution
// timestamp. All times are converted to UTC.
func ParseTimeString(s string) (int64, error) {
	var (
		t   time.Time
		d   time.Duration
		err error
	)

	// Duration
	d, err = time.ParseDuration(s)

	if err == nil {
		return time.Now().UTC().Add(d).Unix(), nil
	}

	// Parse time.
	for _, layout := range timeLayouts {
		t, err = time.Parse(layout, s)

		if err == nil {
			return t.UTC().Unix(), nil
		}
	}

	// Timestamp; assume this is UTC time.
	i, err := strconv.ParseInt(s, 10, 64)

	if err == nil {
		return i, nil
	}

	return 0, fmt.Errorf("[time] could not parse %s", s)
}
