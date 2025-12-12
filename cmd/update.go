package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ColinEge/task-cli/internal/task"
)

func handleUpdate(svc task.Tasker) {
	if len(os.Args) < 4 {
		fmt.Println("Usage: task-cli update <id> <description>")
		return
	}

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Usage: task-cli update <id> <description>")
		return
	}

	err = svc.Update(int64(id), task.Task{Description: os.Args[3]})
	if err != nil {
		fmt.Println(fmt.Errorf("failed update task %d: %w", id, err))
		return
	}
	fmt.Printf("Task updated successfully (ID: %d)\n", id)

}
