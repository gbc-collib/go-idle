package game

import (
	"log/slog"
	"time"
)

type Timer struct {
	ID            string
	RemainingTime time.Duration // TimeToFinish - Time Passed
	OnComplete    func(*GameState)
	OriginalTime  time.Duration // Full length for timer
}

type TimerSystem struct {
}

func (ts *TimerSystem) Process(state *GameState, dt time.Duration) error {
	// Process all active timers
	for i := len(state.ActiveTimers) - 1; i >= 0; i-- {
		timer := &state.ActiveTimers[i]
		timer.RemainingTime -= dt
		slog.Info("Timer Processed", "Timer", timer.ID)

		if timer.RemainingTime <= 0 {
		slog.Info("Timer Completed", "Timer", timer.ID)
			// Timer completed - execute callback
			timer.OnComplete(state)

			// Remove completed timer
			state.ActiveTimers = append(state.ActiveTimers[:i], state.ActiveTimers[i+1:]...)
		}
	}
	return nil
}
