package bot

import "strings"

type ViewerAccessLevel int

const (
	ViewerAccessLevelUser ViewerAccessLevel = iota
	ViewerAccessLevelSubscriber
	ViewerAccessLevelVIP
	ViewerAccessLevelModerator
	ViewerAccessLevelStreamer
)

type Viewer struct {
	ID          string
	Name        string
	AccessLevel ViewerAccessLevel
	DisplayName string
}

func MentionViewer(name string) string {
	clean := strings.TrimPrefix(name, "@")
	return "@" + clean
}
