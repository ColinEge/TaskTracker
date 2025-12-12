package task

import (
	"errors"
	"time"
)

type Tasker interface {
	Add(Task) (int64, error)
}

type Status int

const (
	StatusTodo Status = iota
	StatusInProgress
	StatusDone
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
	// Give task an ID
	tasks, err := load(s.savePath)
	if err != nil {
		if !errors.Is(err, ErrNotExist) {
			return 0, err
		}
		tasks = []Task{}
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
