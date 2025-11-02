package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gbc-collib/go-idle/internal/game"
)

// GameModel is the Bubble Tea model for the game UI.
type GameModel struct {
	game           *game.Game
	buildingsTable table.Model
	showHelp       bool
	startTime      time.Time
	lastUpdate     time.Time
}

// NewGameModel creates and initializes a new game model.
func NewGameModel() (*GameModel, error) {
	g, err := game.NewGame()
	if err != nil {
		return nil, err
	}
	g.Start()

	now := time.Now()
	return &GameModel{
		game:           g,
		buildingsTable: createBuildingsTable(),
		showHelp:       false,
		startTime:      now,
		lastUpdate:     now,
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

// Init initializes the model and starts the tick loop.
func (m GameModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles incoming messages and updates the model.
func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.game.Update()
		m.lastUpdate = time.Time(msg)
		m.updateBuildingsTable()
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	return m, nil
}

func (m GameModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.game.Stop()
		return m, tea.Quit
	case " ":
		m.game.QueueInput(&game.ManualCodeCommand{})
	case "h":
		m.showHelp = !m.showHelp
	}
	return m, nil
}

func (m *GameModel) updateBuildingsTable() {
	state := m.game.GetState()
	rows := BuildingsToTableRows(state.Buildings)
	m.buildingsTable.SetRows(rows)
}

// View renders the UI.
func (m GameModel) View() string {
	state := m.game.GetState()
	uptime := time.Since(m.startTime)

	header := RenderHeader("⚡ IDLE GAME STUDIO ⚡", "Incremental Game Development Engine v1.0.0")

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		RenderResourcePanel(state, uptime),
		RenderBuildingsPanel(m.buildingsTable),
	)

	// Get active timer if exists
	var timer *game.Timer
	if t, found := state.GetTimer("manual_code_"); found {
		timer = &t
	}

	frozen := time.Since(m.lastUpdate) > time.Second
	status := RenderStatusBar("CONTROLS: [SPACE] Write Code  [H] Help  [Q] Quit", timer, frozen)

	help := ""
	if m.showHelp {
		help = "\n" + RenderHelp()
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, mainContent, status, help)
}
