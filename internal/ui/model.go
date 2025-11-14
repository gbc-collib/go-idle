package ui

import (
	"log/slog"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gbc-collib/go-idle/internal/game"
)

type ViewMode int

const (
	ViewMain ViewMode = iota
	ViewShop
	ViewUpgrades
	LastView
)

// Handles intelligently Cycling View modes
func (m *GameModel) cycleViewMode() {
	if m.CurrentView == LastView {
		slog.Info("ViewMode Cycled to", "view", ViewMain)
		m.CurrentView = ViewMain
	}
	m.CurrentView = +1
	slog.Info("ViewMode Cycled to", "view", m.CurrentView+1)

}

// GameModel is the Bubble Tea model for the game UI.
type GameModel struct {
	game           *game.Game
	buildingsTable table.Model
	CurrentView    ViewMode
	shopList       list.Model

	showHelp   bool
	startTime  time.Time
	lastUpdate time.Time
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

func createShopList() list.Model {
	shopItems := make([]shopItem, 2)
	shopItems[0] = shopItem{title: "Vim Plugin", description: "More plugins = better right?", buildingId: "vimPlugin"}
	shopItems[1] = shopItem{title: "Project Managers", description: "Certified Agile Scrum ceremony Grand Druid", buildingId: "projectManager"}
	shopList := list.New(shopItems, list.NewDefaultDelegate(), 0, 0)
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
		switch msg.String() {
		case "q", "ctrl+c":
			m.game.Stop()
			return m, tea.Quit
		case " ":
			m.game.QueueInput(&game.ManualCodeCommand{})
		case "h":
			m.showHelp = !m.showHelp
		case "tab":
			m.cycleViewMode()
		case "1":
			m.CurrentView = 1
		case "2":
			m.CurrentView = 2
		case "3":
			m.CurrentView = 3
		}
		switch m.CurrentView {
		case ViewMain:
			break
		case ViewShop:
			return m.handleShopKeyPress(msg)
		case ViewUpgrades:
			return m, nil
			//return m.updateResearch(msg)
		}
	}
	return m, nil
}

func (m GameModel) handleGlobalKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.game.Stop()
		return m, tea.Quit
	case " ":
		m.game.QueueInput(&game.ManualCodeCommand{})
	case "h":
		m.showHelp = !m.showHelp
	case "tab":
		m.cycleViewMode()
	case "1":
		m.CurrentView = 1
	case "2":
		m.CurrentView = 2
	case "3":
		m.CurrentView = 3
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
	shop := m.renderShop()

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		RenderResourcePanel(state, uptime),
		RenderBuildingsPanel(m.buildingsTable),
		shop,
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
