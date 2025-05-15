package cmd

import (
	"fmt"
	"strings"

	"github.com/priyeshcodes/smart-task-cli/internal/task"

	"github.com/spf13/cobra"
)

const dbPath = "tasks.db" // you can move this to a shared config later

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search for tasks by keyword in title or description",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := strings.ToLower(strings.Join(args, " "))

		// Initialize task store
		store, err := task.NewTaskStore(dbPath)
		if err != nil {
			fmt.Println("Failed to open task store:", err)
			return
		}
		defer store.Close() // Optional: if you have a Close() method to close BoltDB

		// Get all tasks
		tasks, err := store.GetAllTasks()
		if err != nil {
			fmt.Println("Error retrieving tasks:", err)
			return
		}

		// Search for keyword in title or description
		found := false
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), keyword) ||
				strings.Contains(strings.ToLower(t.Description), keyword) {
				fmt.Printf("🔍 %s\n   %s\n\n", t.Title, t.Description)
				found = true
			}
		}

		if !found {
			fmt.Println("No matching tasks found.")
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
