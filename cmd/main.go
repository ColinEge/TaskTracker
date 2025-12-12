package main

import (
	"os"
	"time"

	"github.com/ColinEge/task-cli/internal/cli"
	"github.com/ColinEge/task-cli/internal/task"
)

func main() {
	if len(os.Args) < 2 {
		cli.Help()
		return
	}

	svc := task.NewTaskService(task.WithSavePath("tasks.json"), task.WithTimeFunction(time.Now))

	switch os.Args[1] {
	case "add":
		handleAdd(svc)
	case "update":
		handleUpdate(svc)
	case "delete":
		handleDelete(svc)
	}
}
