package main

import (
	"os"

	"github.com/ColinEge/task-cli/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		cli.Help()
		return
	}

	switch os.Args[1] {
	case "add":
		handleAdd()
	}
}
