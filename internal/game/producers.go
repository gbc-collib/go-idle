package game

import (
	"log/slog"
	"time"
)

type Producer interface {
	Process(state *GameState, deltaTime time.Duration) error
}

type ResourceSystem struct{}

func (rs *ResourceSystem) Process(state *GameState, dt time.Duration) error {
	state.Resources["cpu"] += state.Buildings["compiler"].ProductionRate * dt.Seconds()
	state.Resources["memory"] += state.Buildings["memory"].ProductionRate * dt.Seconds()
	slog.Info("Resource System processed", "state", state)
	return nil
}
