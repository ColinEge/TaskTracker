package task

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	tests := []struct {
		name           string
		testFileName   string
		preventCleanup bool
		input          []Task
		expected       []Task
		falsy          bool
	}{
		{
			name:         "Should add quotes to string",
			testFileName: "tst-TestSaveFunction.json",
			input: []Task{{
				Id:          1,
				Description: "Test the save function",
				Status:      StatusInProgress,
				CreatedAt:   timeMustParse(time.RFC3339, "2025-12-12T13:13:59Z"),
				UpdatedAt:   timeMustParse(time.RFC3339, "2025-12-12T13:13:59Z"),
			}},
			expected: []Task{{
				Id:          1,
				Description: "Test the save function",
				Status:      StatusInProgress,
				CreatedAt:   timeMustParse(time.RFC3339, "2025-12-12T13:13:59Z"),
				UpdatedAt:   timeMustParse(time.RFC3339, "2025-12-12T13:13:59Z"),
			}},
		},
	}

	t.Cleanup(func() {
		for _, tst := range tests {
			if tst.preventCleanup {
				return
			}
			os.Remove(tst.testFileName)
		}
	})
	for _, tst := range tests {
		// Do the saving
		if err := save(tst.testFileName, tst.input); err != nil {
			t.Error(err)
		}
		// Read and compare what was saved
		actual, err := readSaveFile(tst.testFileName)
		if err != nil {
			t.Error(err)
		}
		var tasks []Task
		if err := json.Unmarshal(actual, &tasks); err != nil {
			t.Error(err)
		}
		if tst.falsy {
			if isTasksSame(tasks, tst.expected) {
				t.Errorf("%s expected anything but %v but got %v", tst.name, tst.expected, actual)
			}
		} else {
			if !isTasksSame(tasks, tst.expected) {
				t.Errorf("%s expected %v but got %v", tst.name, tst.expected, actual)
			}
		}
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name            string
		testFileName    string
		preventCleanup  bool
		testFileContent string
		expected        []Task
		falsy           bool
	}{
		{
			name:            "Should read",
			testFileName:    "tst-TestLoadFunction.json",
			testFileContent: `[{"id": 1, "description": "Test the load function", "status": 1, "createdAt": "2025-12-12T13:13:59Z", "updatedAt": "2025-12-12T13:13:59Z"}]`,
			expected: []Task{
				{
					Id:          1,
					Description: "Test the load function",
					Status:      StatusInProgress,
					CreatedAt:   timeMustParse(time.RFC3339, "2025-12-12T13:13:59Z"),
					UpdatedAt:   timeMustParse(time.RFC3339, "2025-12-12T13:13:59Z"),
				},
			},
		},
	}

	t.Cleanup(func() {
		for _, tst := range tests {
			if tst.preventCleanup {
				return
			}
			os.Remove(tst.testFileName)
		}
	})
	for _, tst := range tests {
		if err := os.WriteFile(tst.testFileName, []byte(tst.testFileContent), 0466); err != nil {
			t.Errorf("failed to prepare test file: %s", err.Error())
		}

		actual, err := load(tst.testFileName)
		if err != nil {
			t.Error(err)
		}
		if tst.falsy {
			if isTasksSame(actual, tst.expected) {
				t.Errorf("%s expected anything but %v but got %v", tst.name, tst.expected, actual)
			}
		} else {
			if !isTasksSame(actual, tst.expected) {
				t.Errorf("%s expected %v but got %v", tst.name, tst.expected, actual)
			}
		}
	}
}

func isTasksSame(t1, t2 []Task) bool {
	if len(t1) != len(t2) {
		return false
	}
	for i := range t1 {
		if !isTaskSame(t1[i], t2[i]) {
			return false
		}
	}
	return true
}

func isTaskSame(t1, t2 Task) bool {
	if t1.Id != t2.Id ||
		t1.Description != t2.Description ||
		t1.Status != t2.Status ||
		!t1.CreatedAt.Equal(t2.CreatedAt) ||
		!t1.UpdatedAt.Equal(t2.UpdatedAt) {
		return false
	}
	return true
}

func readSaveFile(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Parses time with given format or panics
func timeMustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}
