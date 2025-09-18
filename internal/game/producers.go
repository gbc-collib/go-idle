package game

import (
	"log/slog"
	"time"
)

type Producer interface {
	Process(state *GameState, deltaTime time.Duration)
}

type ResourceSystem struct{}

func (rs *ResourceSystem) Process(state *GameState, dt time.Duration) {
	state.Resources["cpu"] += state.Buildings["compiler"] * dt.Seconds()
	state.Resources["memory"] += state.Buildings["memory"] * dt.Seconds()
	slog.Info("Resource System processed", "state", state)
}
