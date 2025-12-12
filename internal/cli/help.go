package cli

import (
	"fmt"
)

func Help() {
	fmt.Println("Usage: task-cli <command> [arguments]")
	fmt.Println(`
Commands:
  add <description>  Add a new task`)
}
