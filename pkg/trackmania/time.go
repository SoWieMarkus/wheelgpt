package trackmania

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Time int64

const ExampleTimeFormat = "hh:mm:ss.xxx (hours and minutes are optional)"

var timeRegex = regexp.MustCompile(`^(?:(\d{1,5}):)?(?:(\d{1,5}):)?(\d{1,5})\.(\d{1,3})$`)

func ParseTime(timeStr string) (Time, error) {
	if !timeRegex.MatchString(timeStr) {
		return 0, fmt.Errorf("invalid time format %q, expected %s", timeStr, ExampleTimeFormat)
	}

	parts := strings.Split(timeStr, ".")
	hhmmss := strings.Split(parts[0], ":")

	millisecondString := parts[1]
	for len(millisecondString) < 3 {
		millisecondString += "0"
	}

	milliseconds, err := strconv.ParseInt(millisecondString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse milliseconds: %w", err)
	}

	toInt64 := func(value string) (int64, error) {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	var hours int64
	var minutes int64
	var seconds int64

	seconds, err = toInt64(hhmmss[len(hhmmss)-1])
	if err != nil {
		return 0, fmt.Errorf("parse seconds: %w", err)
	}

	if len(hhmmss) >= 2 {
		minutes, err = toInt64(hhmmss[len(hhmmss)-2])
		if err != nil {
			return 0, fmt.Errorf("parse minutes: %w", err)
		}
	}

	if len(hhmmss) >= 3 {
		hours, err = toInt64(hhmmss[len(hhmmss)-3])
		if err != nil {
			return 0, fmt.Errorf("parse hours: %w", err)
		}
	}

	totalMilliseconds :=
		hours*3_600_000 +
			minutes*60_000 +
			seconds*1_000 +
			milliseconds

	return Time(totalMilliseconds), nil
}

func (t Time) String() string {
	value := int64(t)
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}

	milliseconds := value % 1000
	seconds := (value / 1000) % 60
	minutes := (value / 60_000) % 60
	hours := value / 3_600_000

	if hours > 0 {
		return fmt.Sprintf("%s%d:%02d:%02d.%03d", sign, hours, minutes, seconds, milliseconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%s%d:%02d.%03d", sign, minutes, seconds, milliseconds)
	}

	return fmt.Sprintf("%s%d.%03d", sign, seconds, milliseconds)
}
