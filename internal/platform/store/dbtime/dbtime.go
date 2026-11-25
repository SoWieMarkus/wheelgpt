// Package dbtime maps the DATETIME columns of the Prisma era.
//
// Prisma stores DateTime values in SQLite as integer unix milliseconds, while
// a column DEFAULT CURRENT_TIMESTAMP stores text ("2006-01-02 15:04:05").
// Both exist in production databases, so Time scans either and always writes
// unix milliseconds, which the old Prisma client can still read.
package dbtime

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

// Time is a UTC timestamp stored as unix milliseconds.
type Time struct{ time.Time }

// Now returns the current time truncated to milliseconds.
func Now() Time { return From(time.Now()) }

// From converts t to a Time with millisecond precision.
func From(t time.Time) Time { return Time{t.UTC().Truncate(time.Millisecond)} }

var textLayouts = []string{
	"2006-01-02 15:04:05.999999999-07:00",
	"2006-01-02T15:04:05.999999999-07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04",
	"2006-01-02T15:04",
	"2006-01-02",
	time.RFC3339Nano,
}

// Scan implements sql.Scanner.
func (t *Time) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		t.Time = time.Time{}
		return nil
	case int64:
		t.Time = time.UnixMilli(v).UTC()
		return nil
	case float64:
		t.Time = time.UnixMilli(int64(v)).UTC()
		return nil
	case time.Time:
		t.Time = v.UTC()
		return nil
	case []byte:
		return t.parseText(string(v))
	case string:
		return t.parseText(v)
	default:
		return fmt.Errorf("dbtime: cannot scan %T", src)
	}
}

func (t *Time) parseText(s string) error {
	if ms, err := strconv.ParseInt(s, 10, 64); err == nil {
		t.Time = time.UnixMilli(ms).UTC()
		return nil
	}
	for _, layout := range textLayouts {
		if parsed, err := time.Parse(layout, s); err == nil {
			t.Time = parsed.UTC()
			return nil
		}
	}
	return fmt.Errorf("dbtime: unknown time format %q", s)
}

// Value implements driver.Valuer.
func (t Time) Value() (driver.Value, error) {
	return t.UnixMilli(), nil
}
