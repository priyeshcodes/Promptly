// internal/task/model.go
package task

import (
	"time"
)

// PriorityLevel defines task priority
type PriorityLevel int

const (
	Low PriorityLevel = iota
	Medium
	High
)

func (p PriorityLevel) String() string {
	return [...]string{"Low", "Medium", "High"}[p]
}

// Task defines the structure of a task
type Task struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"` // Markdown
	Priority     PriorityLevel `json:"priority"`
	Deadline     *time.Time    `json:"deadline,omitempty"`
	CompletedAt  time.Time     `json:"completed_at,omitempty"`
	Completed    bool          `json:"completed"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	DependsOn    []string      `json:"depends_on,omitempty"`
	Dependencies []string      `json:"dependencies"` // IDs of tasks that must be completed first
}
