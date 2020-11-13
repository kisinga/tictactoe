package main

import (
	"fmt"
	"io/ioutil"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
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

//move is any valid instance of a turn by a player
type move struct {
	playerTurn
	xoro
}

const (
	player1 playerTurn = false
	player2 playerTurn = true
)

type gridItem struct {
	//played keeps a pointer to the move and the player who made it
	played *move
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
	firstPlayer    playerTurn
	nextPlayerTurn nextPlayerTurn
	grid           [9]*gridItem
	container      *fyne.Container
	xasset         *fyne.StaticResource
	oasset         *fyne.StaticResource
}

func main() {
	myApp := app.New()
	gameWindow := myApp.NewWindow("Tic Tac Toe")
	firstPlayer := player2
	game := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
	xasset, err := ioutil.ReadFile("./assets/x.png")
	if err != nil {
		fmt.Errorf("cant fetch x asset")
	}
	oasset, err := ioutil.ReadFile("./assets/o.png")
	if err != nil {
		fmt.Errorf("cant fetch o asset")
	}
	status := gameStatus{
		window: &gameWindow,
		nextPlayerTurn: nextPlayerTurn{
			turnCount:  0,
			playerTurn: firstPlayer,
			xoro:       x,
		},
		firstPlayer: firstPlayer,
		grid:        [9]*gridItem{},
		container:   game,
		xasset: &fyne.StaticResource{
			StaticName:    "x.png",
			StaticContent: xasset,
		},
		oasset: &fyne.StaticResource{
			StaticName:    "o.png",
			StaticContent: oasset,
		},
	}
	status.resetGame()
	startButton := widget.NewButtonWithIcon("Reset Game", theme.ViewRefreshIcon(), func() {
		status.resetGame()
	})
	gameWindow.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), startButton, game))
	gameWindow.ShowAndRun()
}

//resetGame resets all gridItem to their default values
func (status gameStatus) resetGame() {
	initialised := status.container.Objects
	for i := 0; i < 9; i++ {
		item := createItem(&status)
		status.grid[i] = item
		if initialised != nil {
			status.container.Objects[i] = item
			continue
		}
		status.container.Add(item)
	}
}

//MinSize implements the default fyne definition to provide min icon size
func (b *gridItem) MinSize() fyne.Size {
	return fyne.NewSize(128, 128)
}

//createItem abstracts the creation of a grid item
func createItem(status *gameStatus) *gridItem {
	item := &gridItem{
		played: nil,
		status: status,
	}
	item.SetResource(theme.ViewFullScreenIcon())
	item.ExtendBaseWidget(item)
	return item
}

//Tapped implements default fyne definition
//This function is called when the gridItem is "clicked"
func (b *gridItem) Tapped(ev *fyne.PointEvent) {
	if b.played != nil {
		return
	}
	playerMove := &move{
		playerTurn: b.status.nextPlayerTurn.playerTurn,
		xoro:       !b.status.nextPlayerTurn.xoro,
	}
	b.played = playerMove

	b.status.nextPlayerTurn = nextPlayerTurn{
		turnCount:  b.status.nextPlayerTurn.turnCount + 1,
		playerTurn: !b.status.nextPlayerTurn.playerTurn,
		xoro:       !b.status.nextPlayerTurn.xoro,
	}
	if playerMove.xoro == x {
		b.SetResource(b.status.xasset)
		// b.SetResource(theme.CancelIcon())
	} else {
		b.SetResource(b.status.oasset)

		// b.SetResource(theme.ConfirmIcon())
	}
	b.status.checkStatus()
}

//checkStatus checks the game for any winner.
//It should be called whenever a gridItem is "clicked"
func (status *gameStatus) checkStatus() {
	var winner *playerTurn
	if status.nextPlayerTurn.turnCount == 9 {
		status.announceWinner(nil)
	}
	// Horizontal row checks
	for i := 0; i < 9; i += 3 {
		if status.grid[i].played != nil && status.grid[i+1].played != nil && status.grid[i+2].played != nil {
			if win := status.grid[i].played.xoro == status.grid[i+1].played.xoro && status.grid[i].played.xoro == status.grid[i+2].played.xoro; win {
				winner = &status.grid[i].played.playerTurn
				continue
			}
		}
	}
	if winner != nil {
		status.announceWinner(winner)
		return
	}
	// Vertical row checks
	for i := 0; i < 3; i += 3 {
		if status.grid[i].played != nil && status.grid[i+3].played != nil && status.grid[i+6].played != nil {
			if win := status.grid[i].played.xoro == status.grid[i+3].played.xoro && status.grid[i].played.xoro == status.grid[i+6].played.xoro; win {
				winner = &status.grid[i].played.playerTurn
				continue
			}
		}
	}
	if winner != nil {
		status.announceWinner(winner)
		return
	}
	// Diagonal row checks

	if status.grid[0].played != nil && status.grid[4].played != nil && status.grid[8].played != nil {
		if win := status.grid[0].played.xoro == status.grid[4].played.xoro && status.grid[0].played.xoro == status.grid[8].played.xoro; win {
			winner = &status.grid[0].played.playerTurn
		}
	}
	if winner != nil {
		status.announceWinner(winner)
		return
	}
	if status.grid[2].played != nil && status.grid[4].played != nil && status.grid[6].played != nil {
		if win := status.grid[2].played.xoro == status.grid[4].played.xoro && status.grid[2].played.xoro == status.grid[6].played.xoro; win {
			winner = &status.grid[2].played.playerTurn
		}
	}
	if winner != nil {
		status.announceWinner(winner)
		return
	}
}

//announceWinner displays who won and calls resets the game as well
func (status gameStatus) announceWinner(turn *playerTurn) {
	if turn != nil {
		if *turn {
			dialog.ShowInformation("Player 1 has won!", "Congratulations to player 1 for winning.", *status.window)
			return
		} else {
			dialog.ShowInformation("Player 2 has won!", "Congratulations to player 2 for winning.", *status.window)
			return
		}
	} else {
		dialog.ShowInformation("Draw", "Great game, try again", *status.window)
	}
	status.resetGame()
}
