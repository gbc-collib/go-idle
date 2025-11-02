package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/gbc-collib/go-idle/internal/game"
)

// RenderHeader renders the game title and subtitle.
func RenderHeader(title, subtitle string) string {
	titleRendered := titleStyle.Render(title)
	subtitleRendered := subtitleStyle.Render(subtitle)
	return lipgloss.JoinVertical(lipgloss.Center, titleRendered, subtitleRendered)
}

// RenderResourcePanel renders the resources and stats panel.
func RenderResourcePanel(state game.GameState, uptime time.Duration) string {
	resourcesText := "RESOURCES\n━━━━━━━━━━━━━━━━━━━━━\n\n"
	resourcesText += ResourcesToDisplayText(state.Resources)

	stats := StateToStats(state)
	statsText := fmt.Sprintf(`

PROJECT STATS
━━━━━━━━━━━━━━━━━━━━━

Lines Written: %v
Features Built: %v
Users: %v
Dev Time: %s
Bug Count: %v`,
		stats["linesWritten"],
		stats["featuresBuilt"],
		stats["users"],
		FormatUptime(uptime),
		stats["bugCount"],
	)

	content := resourcesText + statsText
	return resourcePanelStyle.Render(content)
}

// RenderBuildingsPanel renders the buildings table panel.
func RenderBuildingsPanel(buildingsTable table.Model) string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		MarginBottom(1).
		Render("DEVELOPMENT TOOLS")

	tableView := buildingsTable.View()
	content := lipgloss.JoinVertical(lipgloss.Left, header, tableView)
	return buildingTableStyle.Render(content)
}

// RenderStatusBar renders the status bar with controls and optional timer.
func RenderStatusBar(controls string, timer *game.Timer, frozen bool) string {
	statusText := controls
	
	if frozen {
		statusText += "  " + lipgloss.NewStyle().Foreground(errorColor).Render("⚠ IDE FROZEN")
	}

	if timer != nil {
		elapsed := timer.OriginalTime - timer.RemainingTime
		progressPercent := float64(elapsed) / float64(timer.OriginalTime)
		prog := progress.New(progress.WithDefaultGradient())
		progressBar := prog.ViewAs(progressPercent)
		return statusStyle.Render(statusText + "\n" + progressBar)
	}

	return statusStyle.Render(statusText)
}

// RenderHelp renders the help panel.
func RenderHelp() string {
	helpText := `DEVELOPER REFERENCE

Manual Development:
  SPACE     - Write code manually (+1 lines)

Navigation:
  H         - Toggle this help panel
  Q / ^C    - Quit IDE

Coming Soon:
  > code      - Write specific features
  > deploy    - Push to production
  > hire      - Expand your dev team
  > optimize  - Refactor and improve code
  > market    - Promote your idle game`

	return helpStyle.Render(helpText)
}

