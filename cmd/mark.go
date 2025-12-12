package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ColinEge/task-cli/internal/task"
)

func handleMark(svc task.Tasker, status task.Status) {
	showHelp := func() {
		switch status {
		case task.StatusInProgress:
			fmt.Println("Usage: task-cli mark-in-progress <id>")
		case task.StatusDone:
			fmt.Println("Usage: task-cli mark-in-done <id>")
		}
	}

	if len(os.Args) < 3 {
		showHelp()
		return
	}

	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		showHelp()
		return
	}

	err = svc.Mark(int64(id), status)
	statusToString := func() string {
		switch status {
		case task.StatusInProgress:
			return "in-progress"
		case task.StatusDone:
			return "done"
		}
		return ""
	}
	if err != nil {
		fmt.Println(fmt.Errorf("failed to mark task as %s: %w", statusToString(), err))
		return
	}
	fmt.Printf("Task marked as %s successfully (ID: %d)\n", statusToString(), id)
}
