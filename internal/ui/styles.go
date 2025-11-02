package ui

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	primaryColor   = lipgloss.Color("#00FF41") // Matrix green
	secondaryColor = lipgloss.Color("#008F11") // Darker green
	accentColor    = lipgloss.Color("#FFFF00") // Yellow for highlights
	errorColor     = lipgloss.Color("#FF0000") // Red for warnings
)

// Header styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 2).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true).
			MarginBottom(2)
)

// Panel styles
var (
	resourcePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				MarginRight(2).
				Width(25)

	buildingTableStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				Width(50)
)

// Status and help styles
var (
	statusStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			MarginTop(1).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor)
)
