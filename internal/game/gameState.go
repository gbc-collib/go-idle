package game

type GameState struct {
	Resources map[string]float64 `json:"resources"`
	Buildings map[string]float64 `json:"buildings"`
}
