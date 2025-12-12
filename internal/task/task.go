package task

import (
	"errors"
	"fmt"
	"time"
)

type Tasker interface {
	Add(Task) (int64, error)
	Update(id int64, t Task) error
	Delete(id int64) error
	Mark(id int64, status Status) error
	List(status *Status) ([]Task, error)
}

type Status int

const (
	StatusTodo Status = iota
	StatusInProgress
	StatusDone
)

var (
	ErrNotFound = errors.New("task not found")
)

type Task struct {
	Id          int64     `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt,omitzero"`
}

type NowFunc func() time.Time
type TaskService struct {
	savePath string
	now      NowFunc
}

type TaskServiceOption func(svc *TaskService)

func WithSavePath(path string) TaskServiceOption {
	return func(svc *TaskService) {
		svc.savePath = path
	}
}

func WithTimeFunction(f NowFunc) TaskServiceOption {
	return func(svc *TaskService) {
		svc.now = f
	}
}

func NewTaskService(opts ...TaskServiceOption) Tasker {
	svc := TaskService{}
	for _, opt := range opts {
		opt(&svc)
	}
	return svc
}

func (s TaskService) Add(t Task) (int64, error) {
	tasks, err := loadOrCreate(s.savePath)
	if err != nil {
		return 0, err
	}

	// Fill in the tasks blanks
	var maxID int64 = 0
	for _, t := range tasks {
		if t.Id > maxID {
			maxID = t.Id
		}
	}
	t.Id = maxID + 1
	t.CreatedAt = s.now()

	tasks = append(tasks, t)
	err = save(s.savePath, tasks)

	return t.Id, err
}

func (s TaskService) Update(id int64, t Task) error {
	tasks, err := loadOrCreate(s.savePath)
	if err != nil {
		return err
	}

	// Find the task to update
	found := false
	for i, task := range tasks {
		if task.Id != id {
			continue
		}
		found = true
		task.Description = t.Description
		task.UpdatedAt = s.now()
		tasks[i] = task
	}
	if !found {
		return fmt.Errorf("%w: with id %d", ErrNotFound, id)
	}
	return save(s.savePath, tasks)
}

func (s TaskService) Delete(id int64) error {
	tasks, err := loadOrCreate(s.savePath)
	if err != nil {
		return err
	}
	// Find and delete the task
	found := false
	for i, task := range tasks {
		if task.Id != id {
			continue
		}
		found = true
		t2 := append(tasks[:i], tasks[i+1:]...)
		tasks = t2
		break
	}
	if !found {
		return fmt.Errorf("%w: with id %d", ErrNotFound, id)
	}
	return save(s.savePath, tasks)
}

func (s TaskService) Mark(id int64, status Status) error {
	tasks, err := loadOrCreate(s.savePath)
	if err != nil {
		return err
	}
	// Find and update the task status
	found := false
	for i, task := range tasks {
		if task.Id != id {
			continue
		}
		found = true
		task.Status = status
		tasks[i] = task
		break
	}
	if !found {
		return fmt.Errorf("%w: with id %d", ErrNotFound, id)
	}
	return save(s.savePath, tasks)
}

func (s TaskService) List(status *Status) ([]Task, error) {
	tasks, err := loadOrCreate(s.savePath)
	if err != nil {
		return nil, err
	}
	// Return blank or unfiltered list
	if len(tasks) == 0 || status == nil {
		return tasks, nil
	}

	// filter list by status
	var filteredList []Task
	for _, task := range tasks {
		if task.Status != *status {
			continue
		}
		filteredList = append(filteredList, task)
	}
	return filteredList, nil
}

func loadOrCreate(path string) ([]Task, error) {
	tasks, err := load(path)
	if err != nil {
		if !errors.Is(err, ErrFileNotExist) {
			return nil, err
		}
		tasks = []Task{}
	}
	return tasks, nil
}
