// cmd/add.go
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/priyeshcodes/smart-task-cli/internal/task"
	"github.com/spf13/cobra"
)

var (
	title       string
	description string
	priority    string
	deadline    string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Run: func(cmd *cobra.Command, args []string) {
		if title == "" {
			log.Fatal("Title is required")
		}

		// Parse priority
		var p task.PriorityLevel
		switch priority {
		case "low":
			p = task.Low
		case "medium":
			p = task.Medium
		case "high":
			p = task.High
		default:
			log.Fatalf("Invalid priority: %s", priority)
		}

		// Parse deadline
		var dl *time.Time
		if deadline != "" {
			d, err := time.Parse("2006-01-02", deadline)
			if err != nil {
				log.Fatalf("Invalid deadline format. Use YYYY-MM-DD")
			}
			dl = &d
		}

		// Create Task
		newTask := &task.Task{
			Title:       title,
			Description: description,
			Priority:    p,
			Deadline:    dl,
			Completed:   false,
		}

		// Init store
		store, err := task.NewTaskStore("tasks.db")
		if err != nil {
			log.Fatal("Failed to open task store:", err)
		}
		defer store.Close()

		// Save task
		err = store.SaveTask(newTask)
		if err != nil {
			log.Fatal("Failed to save task:", err)
		}

		fmt.Println("✅ Task added successfully!")
	},
}

func init() {
	// Register flags
	addCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the task")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Markdown description")
	addCmd.Flags().StringVarP(&priority, "priority", "p", "low", "Priority: low | medium | high")
	addCmd.Flags().StringVar(&deadline, "deadline", "", "Deadline (YYYY-MM-DD)")

	rootCmd.AddCommand(addCmd)
}
