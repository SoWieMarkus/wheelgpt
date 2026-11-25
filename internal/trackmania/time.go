// Package trackmania holds the game logic that does not need a database:
// race times, medals, map names and the chat texts built from them.
package trackmania

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TimeFormatHelp explains the accepted guess format.
const TimeFormatHelp = "hh:mm:ss.xxx (hours and minutes are optional)"

var timePattern = regexp.MustCompile(`^(?:(\d{1,5}):)?(?:(\d{1,5}):)?(\d{1,5})\.(\d{1,3})$`)

// Time is a race time in milliseconds.
type Time int64

// ParseTime parses "s.xxx", "m:ss.xxx" or "h:mm:ss.xxx". Milliseconds are
// required; one or two digits are padded ("1.5" is 1.500).
func ParseTime(value string) (Time, bool) {
	if !timePattern.MatchString(value) {
		return 0, false
	}
	whole, fraction, _ := strings.Cut(value, ".")
	ms, _ := strconv.ParseInt((fraction + "000")[:3], 10, 64)

	parts := strings.Split(whole, ":")
	var total int64
	multipliers := []int64{1000, 60_000, 3_600_000}
	for i := range parts {
		n, _ := strconv.ParseInt(parts[len(parts)-1-i], 10, 64)
		total += n * multipliers[i]
	}
	return Time(total + ms), true
}

// Milliseconds returns the time in milliseconds.
func (t Time) Milliseconds() int64 { return int64(t) }

// String formats like the game does: "12.345", "1:02.345" or "1:01:02.345".
func (t Time) String() string {
	v := int64(t)
	if v < 0 {
		v = -v
	}
	ms := v % 1000
	s := (v / 1000) % 60
	m := (v / 60_000) % 60
	h := v / 3_600_000
	switch {
	case h > 0:
		return fmt.Sprintf("%d:%02d:%02d.%03d", h, m, s, ms)
	case m > 0:
		return fmt.Sprintf("%d:%02d.%03d", m, s, ms)
	default:
		return fmt.Sprintf("%d.%03d", s, ms)
	}
}

// Abs returns the absolute value.
func (t Time) Abs() Time {
	if t < 0 {
		return -t
	}
	return t
}
