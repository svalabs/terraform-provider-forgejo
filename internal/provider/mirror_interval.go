package provider

import (
	"fmt"
	"time"
)

// canonicalizeMirrorInterval converts Forgejo's duration representation to
// the hour-prefixed format accepted by the Terraform resource validator.
// Forgejo omits a zero-hour prefix, returning "10m0s" for "0h10m0s".
func canonicalizeMirrorInterval(value string) string {
	if value == "" {
		return value
	}

	duration, err := time.ParseDuration(value)
	if err != nil || duration < 0 || duration%time.Second != 0 {
		return value
	}

	seconds := int64(duration / time.Second)
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	return fmt.Sprintf("%dh%dm%ds", hours, minutes, remainingSeconds)
}
