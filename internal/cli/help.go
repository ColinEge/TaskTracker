package cli

import (
	"fmt"
)

func Help() {
	fmt.Println("Usage: task-cli <command> [arguments]")
	fmt.Println(`
Commands:
  add <description>          Add a new task
  update <id> <description>  Update a task
  delete <id>                Delete a task`)
}
