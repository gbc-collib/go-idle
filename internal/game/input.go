package game

import (
	"fmt"
	"log/slog"
	"time"
)

type InputHandlerFunc = func(*GameState, time.Duration) error

type InputToken string

const (
	ManualClick InputToken = "manualClick"
	ManualCode  InputToken = "manualCode"
)

type InputSystem struct {
	inputHandlers map[InputToken]InputHandlerFunc
	inputQueue    []InputToken
}

func NewInputSystem() *InputSystem {
	return &InputSystem{
		inputHandlers: make(map[InputToken]InputHandlerFunc), // Add this
		inputQueue:    make([]InputToken, 0),                 // And this
	}
}

func (is *InputSystem) Process(state *GameState, dt time.Duration) error {
	for _, input := range is.inputQueue {
		handler := is.inputHandlers[input]
		if handler == nil {
			slog.Warn("No handler registered", "input", input)
			continue
		}
		err := handler(state, dt)
		if err != nil {
			return err
		}
	}
	is.inputQueue = is.inputQueue[:0]
	return nil
}

func (is *InputSystem) RegisterHandler(input InputToken, handler InputHandlerFunc) {
	is.inputHandlers[input] = handler
	slog.Info("Registered Handle", "input", input, "handler", handler)
}

func (is *InputSystem) QueueInput(input InputToken) {
	is.inputQueue = append(is.inputQueue, input)
}

func HandleManualCode(state *GameState, dt time.Duration) error {
	// Start a 5-second coding timer
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
			// This runs when timer finishes
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
