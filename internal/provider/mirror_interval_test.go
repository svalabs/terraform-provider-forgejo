package provider

import "testing"

func TestCanonicalizeMirrorInterval(t *testing.T) {
	tests := map[string]string{
		"":        "",
		"10m0s":   "0h10m0s",
		"0h10m0s": "0h10m0s",
		"1h10m0s": "1h10m0s",
		"90m0s":   "1h30m0s",
		"invalid": "invalid",
		"-10m0s":  "-10m0s",
		"1.5s":    "1.5s",
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			if actual := canonicalizeMirrorInterval(input); actual != expected {
				t.Fatalf("canonicalizeMirrorInterval(%q) = %q, want %q", input, actual, expected)
			}
		})
	}
}
