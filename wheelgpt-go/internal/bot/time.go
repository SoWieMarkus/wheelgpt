package bot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const ExampleFormat = "hh:mm:ss.xxx (hours and minutes are optional)"

var timeRegex = regexp.MustCompile(`^(?:(\d{1,5}):)?(?:(\d{1,5}):)?(\d{1,5})(?:\.(\d{1,3}))$`)

// TMTime stores a Trackmania time as total milliseconds.
type TMTime int64

func ParseTMTime(s string) (TMTime, bool) {
	if !timeRegex.MatchString(s) {
		return 0, false
	}
	parts := strings.SplitN(s, ".", 2)
	msStr := "000"
	if len(parts) == 2 {
		msStr = fmt.Sprintf("%-3s", parts[1])[:3]
	}
	ms, _ := strconv.ParseInt(msStr, 10, 64)

	hms := strings.Split(parts[0], ":")
	// reverse
	for i, j := 0, len(hms)-1; i < j; i, j = i+1, j-1 {
		hms[i], hms[j] = hms[j], hms[i]
	}
	parseInt := func(i int) int64 {
		if i >= len(hms) {
			return 0
		}
		v, _ := strconv.ParseInt(hms[i], 10, 64)
		return v
	}
	total := parseInt(2)*3600000 + parseInt(1)*60000 + parseInt(0)*1000 + ms
	return TMTime(total), true
}

func (t TMTime) String() string {
	ms := int64(t) % 1000
	sec := (int64(t) / 1000) % 60
	min := (int64(t) / 60000) % 60
	hrs := int64(t) / 3600000

	pad := func(n int64, w int) string {
		return fmt.Sprintf("%0*d", w, n)
	}
	switch {
	case hrs > 0:
		return fmt.Sprintf("%d:%s:%s.%s", hrs, pad(min, 2), pad(sec, 2), pad(ms, 3))
	case min > 0:
		return fmt.Sprintf("%d:%s.%s", min, pad(sec, 2), pad(ms, 3))
	default:
		return fmt.Sprintf("%d.%s", sec, pad(ms, 3))
	}
}

func (t TMTime) Diff(other TMTime) int64 {
	return int64(t) - int64(other)
}
