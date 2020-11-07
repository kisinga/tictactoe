package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type xoro bool

const (
	o xoro = false
	x xoro = true
)

type playerTurn bool

type gridItem struct {
	played *playerTurn
	value  xoro
	widget.Icon
	status *gameStatus
}
type nextPlayerTurn struct {
	turnCount int
	playerTurn
	xoro
}
type gameStatus struct {
	window         *fyne.Window
	nextPlayerTurn nextPlayerTurn
	grid           [9]gridItem
}

func main() {
	myApp := app.New()
	gameWindow := myApp.NewWindow("Tic Tac Toe")
	status := gameStatus{
		window:         &gameWindow,
		nextPlayerTurn: nextPlayerTurn{},
		grid:           [9]gridItem{},
	}
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
	for i := 0; i < 9; i++ {
		container.Add(createItem(&status))
	}
	gameWindow.SetContent(container)
	gameWindow.ShowAndRun()
}
func (b *gridItem) MinSize() fyne.Size {
	return fyne.NewSize(128, 128)
}
func createItem(status *gameStatus) *gridItem {
	item := &gridItem{
		played: nil,
		value:  o,
		status: status,
	}
	item.SetResource(theme.ViewFullScreenIcon())
	item.ExtendBaseWidget(item)
	return item
}

func (b *gridItem) Tapped(ev *fyne.PointEvent) {
	if b.played != nil {
		return
	}
	currentPlay := !b.status.nextPlayerTurn.xoro
	b.status.nextPlayerTurn = nextPlayerTurn{
		turnCount:  b.status.nextPlayerTurn.turnCount + 1,
		playerTurn: !b.status.nextPlayerTurn.playerTurn,
		xoro:       currentPlay,
	}
	if currentPlay == x {
		b.SetResource(theme.CancelIcon())
	} else {
		b.SetResource(theme.ConfirmIcon())
	}
}
