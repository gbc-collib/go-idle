package ui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/gbc-collib/go-idle/internal/game"
)

// BuildingsToTableRows converts game buildings map to table rows for display.
func BuildingsToTableRows(buildings map[string]game.Building) []table.Row {
	if len(buildings) == 0 {
		return []table.Row{
			{"No dev tools yet...", "", "", ""},
			{"", "", "", ""},
			{"Start with basic", "", "", ""},
			{"code editor!", "", "", ""},
		}
	}

	rows := []table.Row{}
	for _, building := range buildings {
		count := strconv.Itoa(int(building.Count))
		production := fmt.Sprintf("%.1f/s", building.ProductionRate)
		cost := FormatCurrency(building.Cost)
		displayName := FormatBuildingName(building.Name)

		rows = append(rows, table.Row{
			displayName,
			count,
			production,
			cost,
		})
	}
	return rows
}

// ResourcesToDisplayText converts resources map to formatted display text.
func ResourcesToDisplayText(resources map[string]float64) string {
	if len(resources) == 0 {
		return "No resources yet...\n"
	}

	text := ""
	for resourceName, amount := range resources {
		displayName := FormatResourceName(resourceName)
		text += fmt.Sprintf("%s: %.1f\n", displayName, amount)
	}
	return text
}

// StateToStats extracts display statistics from game state.
func StateToStats(state game.GameState) map[string]interface{} {
	stats := make(map[string]interface{})
	stats["linesWritten"] = 0 // TODO: Track in GameState
	stats["featuresBuilt"] = state.Resources["features"]
	stats["users"] = 0 // TODO: Implement users system
	stats["bugCount"] = "∞"
	return stats
}

