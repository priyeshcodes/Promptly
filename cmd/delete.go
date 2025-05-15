// cmd/delete.go
package cmd

import (
	"fmt"
	"log"

	"github.com/priyeshcodes/smart-task-cli/internal/task"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [task ID]",
	Short: "Delete a task by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]

		store, err := task.NewTaskStore("tasks.db")
		if err != nil {
			log.Fatal("Failed to open task store:", err)
		}
		defer store.Close()

		err = store.DeleteTask(taskID)
		if err != nil {
			log.Fatalf("❌ Failed to delete task: %v", err)
		}

		fmt.Println("🗑️ Task deleted successfully.")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
