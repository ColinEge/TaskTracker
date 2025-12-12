package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ColinEge/task-cli/internal/task"
)

func handleList(svc task.Tasker) {
	showHelp := func() {
		fmt.Println("Usage: task-cli list [|todo|in-progress|done]")
	}

	if len(os.Args) < 2 {
		showHelp()
		return
	}

	var status *task.Status = nil

	if len(os.Args) > 2 {
		switch strings.ToLower(os.Args[2]) {
		case "todo":
			s := task.StatusTodo
			status = &s
		case "in-progress":
			s := task.StatusInProgress
			status = &s
		case "done":
			s := task.StatusDone
			status = &s
		}
	}

	list, err := svc.List(status)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to list tasks: %w", err))
		return
	}
	b := strings.Builder{}
	b.WriteString("Tasks:\n\n")

	fmt.Print(formatTasks(list))
}

func formatTasks(tasks []task.Task) string {
	b := strings.Builder{}

	// Work out the length of each column to make tabular format
	longestDesc := 0
	const statusDefaultLength = 4
	const statusInProgressLength = 11
	statusLength := 0
	for _, t := range tasks {
		descLen := len(t.Description)
		if descLen > longestDesc {
			longestDesc = descLen
		}
		if statusLength != statusInProgressLength {
			if t.Status == task.StatusInProgress {
				statusLength = statusInProgressLength
			} else {
				statusLength = statusDefaultLength
			}
		}
	}

	// Now write the tasks out
	for _, t := range tasks {
		b.WriteString(t.Description)
		for i := 0; i < longestDesc+2-len(t.Description); i++ {
			b.WriteRune(' ')
		}
		b.WriteString(t.Status.String())
		b.WriteRune('\n')
	}
	return b.String()
}
