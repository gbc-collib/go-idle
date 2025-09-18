package game

import (
	"log/slog"
)

type Game struct {
	engine *Engine
	state  *GameState
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
		engine: engine,
		state:  state,
	}, nil
}


func (g *Game) Update(){
	g.engine.Update(g.state)
}
