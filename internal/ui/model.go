package ui

import (
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

type GameModel struct {
	game           *game.Game
	buildingsTable table.Model
	shopList       list.Model

	CurrentView ViewMode

	showHelp   bool
	startTime  time.Time
	lastUpdate time.Time

	windowWidth  int
	windowHeight int
}

// ---------- Initialization ----------

func NewGameModel() (*GameModel, error) {
	g, err := game.NewGame()
	if err != nil {
		return nil, err
	}
	g.Start()

	now := time.Now()
	model := &GameModel{
		game:           g,
		buildingsTable: createBuildingsTable(),
		shopList:       createShopList(),
		showHelp:       false,
		startTime:      now,
		lastUpdate:     now,
		CurrentView:    ViewMain,
	}

	return model, nil
}

func createBuildingsTable() table.Model {
	columns := []table.Column{
		{Title: "Building", Width: 15},
		{Title: "Count", Width: 8},
		{Title: "Production", Width: 12},
		{Title: "Cost", Width: 10},
	}
	return table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(false),
		table.WithHeight(7),
	)
}

// ---------- BubbleTea Methods ----------

func (m *GameModel) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time

func (m *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		m.game.Update()
		m.lastUpdate = time.Time(msg)
		m.updateBuildingsTable()
		return m, tickCmd()

	case tea.KeyMsg:
		if cmd, handled := m.handleGlobalKeys(msg); handled {
			return m, cmd
		}

		switch m.CurrentView {
		case ViewShop:
			return m.handleShopKeyPress(msg)
		case ViewUpgrades:
			return m.handleUpgradesKeyPress(msg)
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		// Resize UI components dynamically
		m.buildingsTable.SetWidth(msg.Width / 2)
		m.shopList.SetWidth(msg.Width / 3)
	}
	return m, nil
}

// ---------- Key Handling ----------

func (m *GameModel) handleGlobalKeys(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.game.Stop()
		return tea.Quit, true
	case " ":
		m.game.QueueInput(&game.ManualCodeCommand{})
	case "h":
		m.showHelp = !m.showHelp
	case "tab":
		m.cycleViewMode()
	case "1":
		m.CurrentView = ViewMain
	case "2":
		m.CurrentView = ViewShop
	case "3":
		m.CurrentView = ViewUpgrades
	default:
		return nil, false
	}
	return nil, true
}

func (m *GameModel) cycleViewMode() {
	m.CurrentView++
	if m.CurrentView >= LastView {
		m.CurrentView = ViewMain
	}
}

// ---------- Helpers ----------

func (m *GameModel) updateBuildingsTable() {
	state := m.game.GetState()
	rows := BuildingsToTableRows(state.Buildings)
	m.buildingsTable.SetRows(rows)
}

// ---------- View Rendering ----------

func (m *GameModel) View() string {
	state := m.game.GetState()
	uptime := time.Since(m.startTime)

	header := RenderHeader("⚡ IDLE GAME STUDIO ⚡", "Incremental Game Development Engine v1.0.0")

	var mainContent string
	switch m.CurrentView {
	case ViewMain:
		mainContent = lipgloss.JoinVertical(
			lipgloss.Left,
			RenderResourcePanel(state, uptime),
			RenderBuildingsPanel(m.buildingsTable),
		)
	case ViewShop:
		mainContent = m.shopList.View()
	case ViewUpgrades:
		//	mainContent = RenderUpgradesPanel()
	}

	// Timer status
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

// ---------- Stub Methods for Views ----------

func (m *GameModel) handleUpgradesKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// TODO: implement upgrades navigation
	return m, nil
}
