package internal

import (
	"fmt"
)

var (
	version   = "unknown"
	revision  = "unknown"
	buildDate = "unknown"
)

func GetVersion() string {
	return version
}

func GetRevision() string {
	return revision
}

func GetBuildDate() string {
	return buildDate
}

func GetSysVersion() string {
	return fmt.Sprintf("%s-%s", version, revision)
}
