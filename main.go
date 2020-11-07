package main

import (
	"fmt"
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
}

func main() {
	myApp := app.New()
	gameWindow := myApp.NewWindow("Tic Tac Toe")
	firstPlayer := player2
	status := gameStatus{
		window: &gameWindow,
		nextPlayerTurn: nextPlayerTurn{
			turnCount:  0,
			playerTurn: firstPlayer,
			xoro:       x,
		},
		firstPlayer: firstPlayer,
		grid:        [9]*gridItem{},
	}
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
	for i := 0; i < 9; i++ {
		item := createItem(&status)
		status.grid[i] = item
		container.Add(item)
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
	playerMove := &move{
		playerTurn: b.status.nextPlayerTurn.playerTurn,
		xoro:       !b.status.nextPlayerTurn.xoro,
	}
	b.played = playerMove

	b.status.nextPlayerTurn = nextPlayerTurn{
		turnCount:  b.status.nextPlayerTurn.turnCount + 1,
		playerTurn: !playerMove.playerTurn,
		xoro:       playerMove.xoro,
	}
	if playerMove.xoro == x {
		b.SetResource(theme.CancelIcon())
	} else {
		b.SetResource(theme.ConfirmIcon())
	}
	b.status.checkStatus()
}
func (s *gameStatus) checkStatus() {
	var winner *playerTurn
	// Horizontal row checks
	for i := 0; i < 9; i += 3 {
		fmt.Println(s.grid[i].played, i)
		if s.grid[i].played != nil && s.grid[i+1].played != nil && s.grid[i+2].played != nil {
			if win := s.grid[i].played.xoro == s.grid[i+1].played.xoro == s.grid[i+2].played.xoro; win {
				winner = &s.grid[i].played.playerTurn
				continue
			}
		}
	}
	if winner != nil {
		s.announceWinner(*winner)
		return
	}
	// Vertical row checks
	for i := 0; i < 3; i += 3 {
		if s.grid[i].played != nil && s.grid[i+3].played != nil && s.grid[i+6].played != nil {
			if win := s.grid[i].played.xoro == s.grid[i+3].played.xoro == s.grid[i+6].played.xoro; win {
				winner = &s.grid[i].played.playerTurn
				continue
			}
		}
	}
	if winner != nil {
		s.announceWinner(*winner)
		return
	}
	// Diagonal row checks

	if s.grid[0].played != nil && s.grid[4].played != nil && s.grid[8].played != nil {
		if win := s.grid[0].played.xoro == s.grid[4].played.xoro == s.grid[8].played.xoro; win {
			winner = &s.grid[0].played.playerTurn
		}
	}
	if winner != nil {
		s.announceWinner(*winner)
		return
	}
	if s.grid[2].played != nil && s.grid[4].played != nil && s.grid[6].played != nil {
		if win := s.grid[2].played.xoro == s.grid[4].played.xoro == s.grid[6].played.xoro; win {
			winner = &s.grid[2].played.playerTurn
		}
	}
	if winner != nil {
		s.announceWinner(*winner)
		return
	}
}
func (s gameStatus) announceWinner(turn playerTurn) {
	if turn {
		dialog.ShowInformation("Player 1 has won!", "Congratulations to player 1 for winning.", *s.window)
		return
	}
	dialog.ShowInformation("Player 2 has won!", "Congratulations to player 2 for winning.", *s.window)
}
