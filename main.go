package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	am := NewAsynqMate("127.0.0.1:6379")
	defer am.Close()

	queues, err := am.ListAllQueuesNames(ctx)
	if err != nil {
		panic(err)
	}

	for _, q := range queues {
		fmt.Println(q)
	}

	tasks, err := am.SearchForTask(ctx, TaskStatePending, "without")
	if err != nil {
		panic(err)
	}

	fmt.Println(tasks)
}
