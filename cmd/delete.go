package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ColinEge/task-cli/internal/task"
)

func handleDelete(svc task.Tasker) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: task-cli delete <id>")
		return
	}

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Usage: task-cli delete <id>")
		return
	}

	err = svc.Delete(int64(id))
	if err != nil {
		fmt.Println(fmt.Errorf("failed delete task %d: %w", id, err))
		return
	}
	fmt.Printf("Task deleted successfully (ID: %d)\n", id)

}
