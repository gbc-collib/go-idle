package ui

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gbc-collib/go-idle/internal/game"
	"time"
)

// Styles
var (
	// Colors
	primaryColor   = lipgloss.Color("#00FF41") // Matrix green
	secondaryColor = lipgloss.Color("#008F11") // Darker green
	accentColor    = lipgloss.Color("#FFFF00") // Yellow for highlights
	errorColor     = lipgloss.Color("#FF0000") // Red for warnings

	// Header styles
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

	// Resource panel style
	resourcePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				MarginRight(2).
				Width(25)

	// Building table style
	buildingTableStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				Width(50)

	// Status bar style
	statusStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			MarginTop(1)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			MarginTop(1).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor)
)

type GameModel struct {
	game           *game.Game
	buildingsTable table.Model
	// UI-only state
	showHelp   bool
	lastUpdate time.Time
}

func NewGameModel() (*GameModel, error) {
	g, err := game.NewGame()
	if err != nil {
		return nil, err
	}
	g.Start() // Start the game engine

	// Initialize buildings table
	buildingsTable := createBuildingsTable()

	return &GameModel{
		game:           g,
		buildingsTable: buildingsTable,
		showHelp:       false,
		lastUpdate:     time.Now(),
	}, nil
}

func createBuildingsTable() table.Model {
	columns := []table.Column{
		{Title: "Building", Width: 15},
		{Title: "Count", Width: 8},
		{Title: "Production", Width: 12},
		{Title: "Cost", Width: 10},
	}

	// Start with empty rows - will be populated from GameState
	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		BorderBottom(true).
		Bold(false).
		Foreground(primaryColor)
	s.Selected = s.Selected.
		Foreground(accentColor).
		Bold(false)

	t.SetStyles(s)
	return t
}

type tickMsg time.Time

func (m GameModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.game.Update()
		m.lastUpdate = time.Time(msg)
		// Update buildings table with current game state
		m.updateBuildingsTable()
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.game.Stop()
			return m, tea.Quit
		case " ":
			m.game.QueueInput(game.ManualCode)
		case "h":
			m.showHelp = !m.showHelp
		}
	}
	return m, nil
}

func (m *GameModel) updateBuildingsTable() {
	state := m.game.GetState()

	var rows []table.Row

	if len(state.Buildings) == 0 {
		// Show empty state message
		rows = []table.Row{
			{"No dev tools yet...", "", "", ""},
			{"", "", "", ""},
			{"Start with basic", "", "", ""},
			{"code editor!", "", "", ""},
		}
	} else {
		// Populate from actual game state
		for _, building := range state.Buildings {
			// Format numbers nicely
			count := strconv.Itoa(int(building.Count))
			production := fmt.Sprintf("%.1f/s", building.ProductionRate)
			cost := m.formatCurrency(building.Cost)

			// Format building name nicely
			displayName := strings.Title(strings.ReplaceAll(building.Name, "_", " "))

			rows = append(rows, table.Row{
				displayName,
				count,
				production,
				cost,
			})
		}
	}

	m.buildingsTable.SetRows(rows)
}

func (m GameModel) formatCurrency(amount float64) string {
	if amount >= 1000000 {
		return fmt.Sprintf("%.1fM", amount/1000000)
	} else if amount >= 1000 {
		return fmt.Sprintf("%.1fK", amount/1000)
	}
	return fmt.Sprintf("%.0f", amount)
}

func (m GameModel) View() string {
	// Header
	header := m.renderHeader()

	// Main content area
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderResourcePanel(),
		m.renderBuildingsPanel(),
	)

	// Status bar
	status := m.renderStatusBar()

	// Help panel (if toggled)
	help := ""
	if m.showHelp {
		help = "\n" + m.renderHelp()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		mainContent,
		status,
		help,
	)
}

func (m GameModel) renderHeader() string {
	title := titleStyle.Render("⚡ IDLE GAME STUDIO ⚡")
	subtitle := subtitleStyle.Render("Incremental Game Development Engine v1.0.0")
	return lipgloss.JoinVertical(lipgloss.Center, title, subtitle)
}

func (m GameModel) renderResourcePanel() string {
	state := m.game.GetState()

	// Build resources section from GameState
	resourcesText := "RESOURCES\n━━━━━━━━━━━━━━━━━━━━━\n\n"
	if len(state.Resources) == 0 {
		resourcesText += "No resources yet...\n"
	} else {
		for resourceName, amount := range state.Resources {
			// Format resource names nicely
			displayName := strings.Title(strings.ReplaceAll(resourceName, "_", " "))
			resourcesText += fmt.Sprintf("%s: %.1f\n", displayName, amount)
		}
	}

	// Add stats section with game dev flavor
	statsText := fmt.Sprintf(`

PROJECT STATS
━━━━━━━━━━━━━━━━━━━━━

Lines Written: 0
Features Built: %s
Users: 0
Dev Time: %s
Bug Count: ∞`,
		state.Resources["features"],
		m.formatUptime(),
	)

	content := resourcesText + statsText
	return resourcePanelStyle.Render(content)
}

func (m GameModel) renderBuildingsPanel() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		MarginBottom(1).
		Render("DEVELOPMENT TOOLS")

	tableView := m.buildingsTable.View()

	content := lipgloss.JoinVertical(lipgloss.Left, header, tableView)
	return buildingTableStyle.Render(content)
}

func (m GameModel) renderStatusBar() string {
    controls := "CONTROLS: [SPACE] Write Code  [H] Help  [Q] Quit"
    if time.Since(m.lastUpdate) > time.Second {
        controls += "  " + lipgloss.NewStyle().Foreground(errorColor).Render("⚠ IDE FROZEN")
    }

    state := m.game.GetState()
    manualCodeTimer, wasFound := state.GetTimer("manual_code_")
    if wasFound {
        // Calculate progress (0.0 to 1.0)
        elapsed := manualCodeTimer.OriginalTime - manualCodeTimer.RemainingTime
        progressPercent := float64(elapsed) / float64(manualCodeTimer.OriginalTime)

        // Create progress bar
        prog := progress.New(progress.WithDefaultGradient())
        progressBar := prog.ViewAs(progressPercent)

        // Combine controls and progress bar
        return statusStyle.Render(controls + "\n" + progressBar)
    }

    return statusStyle.Render(controls)
}

func (m GameModel) renderHelp() string {
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

func (m GameModel) formatUptime() string {
	uptime := time.Since(m.lastUpdate)
	if uptime < time.Minute {
		return fmt.Sprintf("%.0fs", uptime.Seconds())
	}
	return fmt.Sprintf("%.0fm", uptime.Minutes())
}
