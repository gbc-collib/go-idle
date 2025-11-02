package game

import (
	"fmt"
	"log/slog"
	"time"
)

// Command interface - all commands implement this
type Command interface {
	Execute(*GameState) error
}

// Concrete command types
type ManualClickCommand struct{}
type ManualCodeCommand struct{}
type PurchaseBuildingCommand struct {
	BuildingType string
}

// Execute implementations
func (c *ManualClickCommand) Execute(state *GameState) error {
	if state.Resources == nil {
		state.Resources = make(map[string]float64)
	}
	state.Resources["code"] += 1.0
	return nil
}

func (c *ManualCodeCommand) Execute(state *GameState) error {
	manualCodeId := "manual_code_"
	_, alreadyActive := state.GetTimer(manualCodeId)
	if alreadyActive {
		slog.Info("Coding Already Active.")
		return nil
	}

	timer := Timer{
		ID:            manualCodeId + fmt.Sprint(time.Now().UnixNano()),
		RemainingTime: 5 * time.Second,
		OnComplete: func(gameState *GameState) {
			if gameState.Resources == nil {
				gameState.Resources = make(map[string]float64)
			}
			gameState.Resources["features"] += 1.0
			slog.Info("Feature completed!", "features", gameState.Resources["features"])
		},
		OriginalTime: 5 * time.Second,
	}
	state.ActiveTimers = append(state.ActiveTimers, timer)
	slog.Info("Started coding...", "duration", "5s")
	return nil
}

func (c *PurchaseBuildingCommand) Execute(state *GameState) error {
	building, exists := state.Buildings[c.BuildingType]
	if !exists {
		return fmt.Errorf("building type %s not found", c.BuildingType)
	}

	if state.Resources["code"] < building.Cost {
		return fmt.Errorf("insufficient code: need %.0f, have %.0f",
			building.Cost, state.Resources["code"])
	}

	state.Resources["code"] -= building.Cost
	building.Count++
	building.Cost *= 1.15 // scaling
	state.Buildings[c.BuildingType] = building

	return nil
}

// InputSystem now processes Commands
type InputSystem struct {
	commandQueue []Command
}

func NewInputSystem() *InputSystem {
	return &InputSystem{
		commandQueue: make([]Command, 0),
	}
}

func (is *InputSystem) Process(state *GameState, _ time.Duration) error {
	for _, cmd := range is.commandQueue {
		if err := cmd.Execute(state); err != nil {
			slog.Error("Command execution failed", "error", err, "command", fmt.Sprintf("%T", cmd))
			// Decide: continue or return? Probably continue for games
			continue
		}
	}
	is.commandQueue = is.commandQueue[:0]
	return nil
}

func (is *InputSystem) QueueCommand(cmd Command) {
	is.commandQueue = append(is.commandQueue, cmd)
}

