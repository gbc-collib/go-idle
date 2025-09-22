package game

import (
	"encoding/json"
	"os"

	"github.com/gbc-collib/go-idle/internal/game"
)

const SAVE_FILE = "./save.json"

func load_from_file() *game.GameState, error {
	data, err := os.ReadFile(SAVE_FILE)
	if err != nil {
		return nil, err
	}
	var gameState game.GameState
	err = json.Unmarshal(data, gameState)
		if err != nil {
		return nil, err
	}
	return gameState, nil
}

func (gs GameState) save_to_file() error {

}
