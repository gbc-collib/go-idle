package game

import "time"

type Upgrade struct {
}

type UpgradeSystem struct {
	upgrades []Upgrade
	luck     int
	cost     int
}

func (us *UpgradeSystem) Process(gs *GameState, dt time.Duration) error {
	return nil

}
