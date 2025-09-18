package main

import (
	"github.com/gbc-collib/go-idle/internal/game"
	"log/slog"
)

func main() {
	newGame, err := game.NewGame()
	if err != nil {
		slog.Error("Error Encountered Starting Game", "error", err)
	}
	newGame.Start()

	for  {
		newGame.Update()
	}
}
