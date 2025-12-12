package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ColinEge/task-cli/internal/cli"
	"github.com/ColinEge/task-cli/internal/task"
)

func handleAdd() {
	if len(os.Args) < 2 {
		cli.Help()
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage: task-cli add <description>")
		return
	}

	id, err := task.NewTaskService(task.WithSavePath("tasks.json"), task.WithTimeFunction(time.Now)).Add(task.Task{Description: os.Args[2]})
	if err != nil {
		fmt.Println(fmt.Errorf("failed add task to list: %w", err))
		return
	}
	fmt.Printf("Task added successfully (ID: %d)\n", id)

}
