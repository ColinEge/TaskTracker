# Task Tracker CLI Tool

Task tracker CLI is a simple todo Go program using requirements from `roadmap.sh`. Tasks can be viewed, added, removed, updated and deleted from the CLI and stored as JSON for persistence.

Project Requirements from [roadmap.sh](https://roadmap.sh/projects/task-tracker)

## Building
```bash
# Clone the repository
git clone https://github.com/ColinEge/TaskTracker.git

# Navigate to the project directory
cd TaskTracker

# Build the application for mac/linux
go build -o .\.build\task-cli .\cmd\

# OR Build the application for windows
go build -o .\.build\task-cli.exe .\cmd\

# Navigate to the build directory
cd .\.build

# Run the application
./task-cli
```

## Example Usage
```bash
# Adding a new task
task-cli add "Buy groceries"
# Output: Task added successfully (ID: 1)

# Updating and deleting tasks
task-cli update 1 "Buy groceries and cook dinner"
task-cli delete 1

# Marking a task as in progress or done
task-cli mark-in-progress 1
task-cli mark-done 1

# Listing all tasks
task-cli list

# Listing tasks by status
task-cli list done
task-cli list todo
task-cli list in-progress
```