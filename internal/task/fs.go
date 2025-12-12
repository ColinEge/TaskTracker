package task

import (
	"encoding/json"
	"errors"
	"os"
)

var ErrNotExist = os.ErrNotExist

func save(savePath string, tasks []Task) error {
	js, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	if err := os.WriteFile(savePath, js, 0644); err != nil {
		return err
	}
	return nil
}

func load(savePath string) ([]Task, error) {
	bytes, err := os.ReadFile(savePath)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(bytes, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// uses os.Remove but fails silently if path is not found
func deleteFile(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
