package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// canonicalizeMirrorInterval converts Forgejo's duration representation to the
// hour-prefixed format the provider stores in state. Forgejo omits a zero-hour
// prefix, returning "10m0s" for "0h10m0s".
//
// The output format must stay in sync with the mirror_interval regex validator
// in repositoryResource.Schema() — see repository_resource.go. An unparseable
// value is returned unchanged, which restores the "Provider produced
// inconsistent result after apply" error the canonicalization prevents, so the
// fallback is logged.
func canonicalizeMirrorInterval(ctx context.Context, value string) string {
	if value == "" {
		return value
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		tflog.Warn(ctx, "Unable to parse mirror interval, storing API value verbatim", map[string]any{
			"mirror_interval": value,
			"error":           err.Error(),
		})

		return value
	}

	if duration < 0 || duration%time.Second != 0 {
		tflog.Warn(ctx, "Mirror interval is negative or sub-second, storing API value verbatim", map[string]any{
			"mirror_interval": value,
		})

		return value
	}

	seconds := int64(duration / time.Second)
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	return fmt.Sprintf("%dh%dm%ds", hours, minutes, remainingSeconds)
}
