package bot

import "strings"

type AccessLevel int

const (
	AccessUser       AccessLevel = 0
	AccessSubscriber AccessLevel = 1
	AccessVIP        AccessLevel = 2
	AccessMod        AccessLevel = 3
	AccessStreamer    AccessLevel = 4
)

type User struct {
	ID          string
	Name        string
	DisplayName string
	Level       AccessLevel
}

func MentionUser(name string) string {
	clean := strings.TrimPrefix(name, "@")
	return "@" + clean
}
