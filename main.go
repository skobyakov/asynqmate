package main

import (
	"context"

	"github.com/rivo/tview"
)

func main() {
	ctx := context.Background()

	am := NewAsynqMate("127.0.0.1:6379")
	defer am.Close()

	queues, err := am.ListAllQueuesNames(ctx)
	if err != nil {
		panic(err)
	}

	tasks, err := am.SearchForTasks(ctx, TaskStatePending, "with")
	if err != nil {
		panic(err)
	}

	pages := tview.NewPages()

	queuesList := tview.NewList().ShowSecondaryText(false)
	queuesList.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		pages.SwitchToPage("tasks")
	})

	tasksList := tview.NewList()
	tasksList.SetDoneFunc(func() {
		pages.SwitchToPage("queues")
	})

	pages.SetBorder(true).SetTitle("asynqmate").SetBorderPadding(1, 1, 1, 1)
	pages.AddPage("queues", queuesList, true, true)
	pages.AddPage("tasks", tasksList, true, false)

	for _, q := range queues {
		queuesList.AddItem(q, "", 0, nil)
	}

	for _, t := range tasks {
		tasksList.AddItem(t.ID, t.Msg, 0, nil)
	}

	app := tview.NewApplication()
	app.SetRoot(pages, true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
