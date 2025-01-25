package main

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	ctx := context.Background()

	am := NewAsynqMate("127.0.0.1:6379")
	defer am.Close()

	list := tview.NewList()
	res, err := am.SearchForTask(ctx, TaskStatePending, "without")
	if err != nil {
		panic(err)
	}

	items, ok := res.([]interface{})
	if !ok {
		panic("unexpected type")
	}

	for _, item := range items {
		list.AddItem(fmt.Sprintf("%v", item), "", 0, nil)
	}

	app := tview.NewApplication()
	frame := tview.NewFrame(list).
		SetBorders(1, 2, 2, 2, 4, 4).
		AddText("Header left", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("Header middle", true, tview.AlignCenter, tcell.ColorWhite).
		AddText("Header right", true, tview.AlignRight, tcell.ColorWhite).
		AddText("Header second middle", true, tview.AlignCenter, tcell.ColorRed).
		AddText("Footer middle", false, tview.AlignCenter, tcell.ColorGreen).
		AddText("Footer second middle", false, tview.AlignCenter, tcell.ColorGreen)
	if err := app.SetRoot(frame, true).SetFocus(frame).Run(); err != nil {
		panic(err)
	}
}
