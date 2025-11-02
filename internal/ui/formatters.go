package ui

import (
	"fmt"
	"strings"
	"time"
)

// FormatCurrency formats large numbers with K/M suffixes.
func FormatCurrency(amount float64) string {
	if amount >= 1000000 {
		return fmt.Sprintf("%.1fM", amount/1000000)
	} else if amount >= 1000 {
		return fmt.Sprintf("%.1fK", amount/1000)
	}
	return fmt.Sprintf("%.0f", amount)
}

// FormatUptime formats a duration into human-readable uptime.
func FormatUptime(duration time.Duration) string {
	if duration < time.Minute {
		return fmt.Sprintf("%.0fs", duration.Seconds())
	}
	return fmt.Sprintf("%.0fm", duration.Minutes())
}

// FormatResourceName converts snake_case to Title Case.
func FormatResourceName(name string) string {
	return strings.Title(strings.ReplaceAll(name, "_", " "))
}

// FormatBuildingName converts snake_case to Title Case.
func FormatBuildingName(name string) string {
	return strings.Title(strings.ReplaceAll(name, "_", " "))
}

