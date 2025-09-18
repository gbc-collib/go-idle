package game

import (
	"fmt"
	"log/slog"
	"time"
)

type System interface {
	Process(state *GameState, deltaTime time.Duration)
}

type Engine struct {
	systems  []System
	lastTick time.Time
	running  bool
}

func (e *Engine) Update(state *GameState) {
	if !e.running {
		return // Don't process if paused/stopped
	}

	if time.Since(e.lastTick) >= time.Second {
		e.tick(state)
		e.lastTick = time.Now()
		slog.Debug("Game Tick Processed", "lastTick", e.lastTick)
	}
}

func (e *Engine) tick(state *GameState) {
	timeDelta := time.Since(e.lastTick)
	for _, system := range e.systems {
		system.Process(state, timeDelta)
		slog.Debug("Processed System", "systemType", fmt.Sprintf("%T", system))
	}
}

func (e *Engine) AddSystem(system System) {
	slog.Info("Registering System", "system", fmt.Sprintf("%T", system))
	e.systems = append(e.systems, system)

}

func (e *Engine) Start() {
	e.running = true
}

func (e *Engine) Stop() {
	e.running = false
}

func (e *Engine) Pause() {
	e.running = false
}

func NewEngine() *Engine {
	return &Engine{
		systems:  []System{},
		lastTick: time.Now(),
		running:  false,
	}
}
