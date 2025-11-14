package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gbc-collib/go-idle/internal/game"
)

type shopItem struct {
	//What we render
	title string
	//What gamestate/engine know building as
	buildingId string
	//Desc we show idk yet
	description string
}

func (i shopItem) Title() string { return i.title }

func (i shopItem) Description() string { return i.description }

func (i shopItem) FilterValue() string { return i.title }

func (i shopItem) BuildingId() string { return i.buildingId }

func (m GameModel) handleShopKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		selectedItem := m.shopList.SelectedItem()
		if item, ok := selectedItem.(*shopItem); ok {
			m.game.QueueInput(&game.PurchaseBuildingCommand{
				BuildingType: item.buildingId,
				Amount:       1,
			})
		}
	}
	return m, nil
}

func (m GameModel) renderShop() string {
	return m.shopList.View() // Render shop list
}
