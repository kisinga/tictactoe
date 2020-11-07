package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"testing"
)

func Test_gridItem_Tapped(t *testing.T) {
	type fields struct {
		played *move
		Icon   widget.Icon
		status *gameStatus
	}
	type args struct {
		ev *fyne.PointEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &gridItem{
				played: tt.fields.played,
				Icon:   tt.fields.Icon,
				status: tt.fields.status,
			}
		})
	}
}
