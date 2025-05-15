// cmd/complete.go
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/priyeshcodes/smart-task-cli/internal/task"
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete [task ID]",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]

		store, err := task.NewTaskStore("tasks.db")
		if err != nil {
			log.Fatal("Failed to open task store:", err)
		}
		defer store.Close()

		// Fetch task
		t, err := store.GetTaskByID(taskID)
		if err != nil {
			log.Fatalf("Failed to retrieve task: %v", err)
		}

		if !checkDependencies(store, t) {
			fmt.Printf("❌ Task '%s' cannot be completed. The following tasks need to be completed first:\n", t.Title)
			for _, depID := range t.DependsOn {
				depTask, _ := store.GetTaskByID(depID)
				if depTask != nil && !depTask.Completed {
					fmt.Printf(" - Task '%s' (ID: %s)\n", depTask.Title, depTask.ID)
				}
			}
			return
		}

		// Already completed?
		if t.Completed {
			fmt.Println("✔️ Task already marked as completed.")
			return
		}

		// Update and save
		t.Completed = true
		t.CompletedAt = time.Now()
		if err := store.SaveTask(t); err != nil {
			log.Fatalf("Failed to update task: %v", err)
		}

		fmt.Printf("✅ Task '%s' marked as completed!\n", t.Title)
	},
}

func checkDependencies(store *task.TaskStore, task *task.Task) bool {
	// Iterate through dependencies and check if they are completed
	for _, depID := range task.DependsOn {
		depTask, err := store.GetTaskByID(depID)
		if err != nil || depTask == nil || !depTask.Completed {
			// If any dependency is not completed, return false
			return false
		}
	}
	// If all dependencies are completed, return true
	return true
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
