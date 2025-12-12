package task

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	// use same time for all tests to account for file creation and reading time
	testTime := time.Now()
	timeBytes, err := testTime.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name                   string
		preExistingFileContent string
		newTask                Task
		expectedID             int64
		expectedFileContent    string
		falsy                  bool
		preventCleanup         bool
	}{
		{
			name:                   "tasksStartWithID1",
			preExistingFileContent: ``,
			newTask:                Task{Description: "Test The add function", Status: StatusInProgress},
			expectedID:             1,
			expectedFileContent:    `[{"id":1,"description":"Test The add function","status":1,"createdAt":` + string(timeBytes) + `}]`,
		},
		{
			name:                   "tasksAppendAsID2",
			preExistingFileContent: `[{"id":1,"description":"Test The add function","status":1,"createdAt":` + string(timeBytes) + `}]`,
			newTask:                Task{Description: "Test The add function appends stuff", Status: StatusTodo},
			expectedID:             2,
			expectedFileContent:    `[{"id":1,"description":"Test The add function","status":1,"createdAt":` + string(timeBytes) + `},{"id":2,"description":"Test The add function appends stuff","status":0,"createdAt":` + string(timeBytes) + `}]`,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			fileName := "test-" + tst.name + ".json"

			// Cleanup files when done
			t.Cleanup(func() {
				if err := deleteFile(fileName); err != nil {
					log.Default().Print(err)
				}
			})

			// Create needed pre-test files
			if tst.preExistingFileContent != "" {
				if err := os.WriteFile(fileName, []byte(tst.preExistingFileContent), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Add a task
			svc := NewTaskService(WithSavePath(fileName), WithTimeFunction(func() time.Time { return testTime }))
			id, err := svc.Add(tst.newTask)
			if err != nil {
				t.Error(err)
			}
			if id != tst.expectedID {
				t.Errorf("%s expected an ID of %d but got %d", tst.name, tst.expectedID, id)
			}

			// Check if the file state is as expected
			bytes, err := os.ReadFile(fileName)
			if err != nil {
				t.Fatal(err)
			}
			if string(bytes) != tst.expectedFileContent {
				t.Errorf("%s expected a file content of %s but got %s", tst.name, tst.expectedFileContent, string(bytes))
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	// use same time for all tests to account for file creation and reading time
	testTime := time.Now()
	timeBytes, err := testTime.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name                   string
		preExistingFileContent string
		newTask                Task
		expectedFileContent    string
		expectedError          error
		falsy                  bool
		preventCleanup         bool
	}{
		{
			name:                   "tasksUpdateDescriptionAndStatus",
			preExistingFileContent: `[{"id":1,"description":"Test The update function","status":1,"createdAt":` + string(timeBytes) + `}]`,
			newTask:                Task{Id: 1, Description: "Test The update function works right", Status: StatusDone},
			expectedFileContent:    fmt.Sprintf(`[{"id":1,"description":"Test The update function works right","status":2,"createdAt":%s,"updatedAt":%s}]`, string(timeBytes), string(timeBytes)),
		},
		{
			name:                   "tasksErrorWithInvalidID16",
			preExistingFileContent: `[{"id":1,"description":"Test The update function fails","status":1,"createdAt":` + string(timeBytes) + `}]`,
			newTask:                Task{Id: 16, Description: "Test The update function fails properly", Status: StatusDone},
			expectedError:          ErrNotFound,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			fileName := "test-" + tst.name + ".json"

			// Cleanup files when done
			t.Cleanup(func() {
				if err := deleteFile(fileName); err != nil {
					log.Default().Print(err)
				}
			})

			// Create needed pre-test files
			if tst.preExistingFileContent != "" {
				if err := os.WriteFile(fileName, []byte(tst.preExistingFileContent), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Add a task
			svc := NewTaskService(WithSavePath(fileName), WithTimeFunction(func() time.Time { return testTime }))
			if err := svc.Update(tst.newTask.Id, tst.newTask); err != nil {
				if errors.Is(err, tst.expectedError) {
					return
				}
				t.Error(err)
			}

			// Check if the file state is as expected
			bytes, err := os.ReadFile(fileName)
			if err != nil {
				t.Fatal(err)
			}
			if string(bytes) != tst.expectedFileContent {
				t.Errorf("%s expected a file content of %s but got %s", tst.name, tst.expectedFileContent, string(bytes))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	// use same time for all tests to account for file creation and reading time
	testTime := time.Now()
	timeBytes, err := testTime.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name                   string
		preExistingFileContent string
		id                     int64
		expectedFileContent    string
		expectedError          error
		falsy                  bool
		preventCleanup         bool
	}{
		{
			name:                   "tasksDeleteFromFile",
			preExistingFileContent: `[{"id":1,"description":"Test The delete function","status":1,"createdAt":` + string(timeBytes) + `}]`,
			id:                     1,
			expectedFileContent:    `[]`,
		},
		{
			name:                   "tasksErrorWithInvalidID16",
			preExistingFileContent: `[{"id":1,"description":"Test The delete function fails","status":1,"createdAt":` + string(timeBytes) + `}]`,
			id:                     16,
			expectedError:          ErrNotFound,
		},
		{
			name:                   "tasksDeleteFromFileAndLeavesTheRest",
			preExistingFileContent: `[{"id":1,"description":"Test The delete function deletes this","status":1,"createdAt":` + string(timeBytes) + `},{"id":2,"description":"Test The delete function leaves this","status":2,"createdAt":` + string(timeBytes) + `}]`,
			id:                     1,
			expectedFileContent:    fmt.Sprintf(`[{"id":2,"description":"Test The delete function leaves this","status":2,"createdAt":%s}]`, string(timeBytes)),
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			fileName := "test-" + tst.name + ".json"

			// Cleanup files when done
			t.Cleanup(func() {
				if err := deleteFile(fileName); err != nil {
					log.Default().Print(err)
				}
			})

			// Create needed pre-test files
			if tst.preExistingFileContent != "" {
				if err := os.WriteFile(fileName, []byte(tst.preExistingFileContent), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Add a task
			svc := NewTaskService(WithSavePath(fileName), WithTimeFunction(func() time.Time { return testTime }))
			if err := svc.Delete(tst.id); err != nil {
				if errors.Is(err, tst.expectedError) {
					return
				}
				t.Error(err)
			}

			// Check if the file state is as expected
			bytes, err := os.ReadFile(fileName)
			if err != nil {
				t.Fatal(err)
			}
			if string(bytes) != tst.expectedFileContent {
				t.Errorf("%s expected a file content of %s but got %s", tst.name, tst.expectedFileContent, string(bytes))
			}
		})
	}
}

func TestMark(t *testing.T) {
	// use same time for all tests to account for file creation and reading time
	testTime := time.Now()
	timeBytes, err := testTime.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name                   string
		preExistingFileContent string
		id                     int64
		newStatus              Status
		expectedFileContent    string
		expectedError          error
		falsy                  bool
		preventCleanup         bool
	}{
		{
			name:                   "tasksMarkUpdatesStatus",
			preExistingFileContent: `[{"id":1,"description":"Test The mark function","status":1,"createdAt":` + string(timeBytes) + `}]`,
			id:                     1,
			newStatus:              2,
			expectedFileContent:    `[{"id":1,"description":"Test The mark function","status":2,"createdAt":` + string(timeBytes) + `}]`,
		},
		{
			name:                   "tasksErrorWithInvalidID16",
			preExistingFileContent: `[{"id":1,"description":"Test The mark function","status":1,"createdAt":` + string(timeBytes) + `}]`,
			id:                     16,
			newStatus:              2,
			expectedError:          ErrNotFound,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			fileName := "test-" + tst.name + ".json"

			// Cleanup files when done
			t.Cleanup(func() {
				if err := deleteFile(fileName); err != nil {
					log.Default().Print(err)
				}
			})

			// Create needed pre-test files
			if tst.preExistingFileContent != "" {
				if err := os.WriteFile(fileName, []byte(tst.preExistingFileContent), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Add a task
			svc := NewTaskService(WithSavePath(fileName), WithTimeFunction(func() time.Time { return testTime }))
			if err := svc.Mark(tst.id, tst.newStatus); err != nil {
				if errors.Is(err, tst.expectedError) {
					return
				}
				t.Error(err)
			}

			// Check if the file state is as expected
			bytes, err := os.ReadFile(fileName)
			if err != nil {
				t.Fatal(err)
			}
			if string(bytes) != tst.expectedFileContent {
				t.Errorf("%s expected a file content of %s but got %s", tst.name, tst.expectedFileContent, string(bytes))
			}
		})
	}
}
