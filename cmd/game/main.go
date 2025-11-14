
package main

import (
	"log/slog"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbc-collib/go-idle/internal/ui"
	"github.com/gbc-collib/go-idle/internal/game"
)

func setupLogging() {
    logFile, err := os.OpenFile("game.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }

    logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
    slog.SetDefault(logger)
}

func main() {
	setupLogging()
	// Create the Bubble Tea model
	model, err := ui.NewGameModel()
	if err != nil {
		slog.Error("Error creating game model", "error", err)
		os.Exit(1)
	}
	err = game.InitGameData();
	if err != nil {
		slog.Error("Could not init Game data", "error", err)
		os.Exit(1)
	}

	// Create the Bubble Tea program
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Start the UI
	if _, err := program.Run(); err != nil {
		slog.Error("Error running program", "error", err)
		os.Exit(1)
	}
}


