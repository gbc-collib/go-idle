package main

import (
	"github.com/gbc-collib/go-idle/internal/game"
	"log/slog"
)

func main() {
	_, err := game.NewGame()
	if err != nil {
		slog.Error("Error Encountered Starting Game", "error", err)
	}
}
