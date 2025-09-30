package game

import "strings"

type Building struct {
	Count          int
	Name           string
	ProductionRate float64
	Cost           float64
}
type GameState struct {
	Resources    map[string]float64
	Buildings    map[string]Building
	ActiveTimers []Timer
}

func (gs *GameState) ManualCodeTimerActive() bool {
	hasManualCodeTimer := false

	for _, timer := range gs.ActiveTimers {
		if strings.HasPrefix(timer.ID, "manual_click") {
			hasManualCodeTimer = true
			break
		}
	}
	return hasManualCodeTimer
}

func (gs *GameState) GetTimer(id string) (Timer, bool) {
	for _, timer := range gs.ActiveTimers {
		if strings.HasPrefix(timer.ID, id) {
			return timer, true
		}
	}
	return Timer{}, false
}
