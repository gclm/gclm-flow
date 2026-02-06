package cli

import (
	"strings"
)

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// formatError formats an error message for display
func formatError(err error) string {
	if err == nil {
		return ""
	}
	return strings.TrimSpace(err.Error())
}
