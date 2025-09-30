package game

import (
	"log/slog"
)

type Game struct {
	Engine *Engine
	State  *GameState
}

func NewGame() (*Game, error) {
	slog.Info("Creating new Game")

	state := &GameState{
		Resources: make(map[string]float64),
	}
	slog.Debug("Created GameState", "GameState", state)

	engine := NewEngine()
	slog.Debug("Created Engine", "Engine", engine)

	return &Game{
		Engine: engine,
		State:  state,
	}, nil
}


func (g *Game) Update(){
	g.Engine.Update(g.State)
}

func (g *Game) Start(){
	g.Engine.Start()
}

func (g *Game) Stop(){
	g.Engine.Stop()
}

func (g *Game) QueueInput(token InputToken) {
    if inputSys := GetSystem[*InputSystem](g.Engine); inputSys != nil {
        inputSys.QueueInput(token)
    }
}

func (g *Game) GetState() GameState{
	return *g.State
}
