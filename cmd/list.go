// cmd/list.go
package cmd

import (
	"fmt"
	"log"
	"sort"
	"strings"
	// "time"

	"github.com/priyeshcodes/smart-task-cli/internal/task"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks with sorting",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := task.NewTaskStore("tasks.db")
		if err != nil {
			log.Fatal("Failed to open task store:", err)
		}
		defer store.Close()

		tasks, err := store.GetAllTasks()
		if err != nil {
			log.Fatal("Failed to fetch tasks:", err)
		}

		// Sort: Incomplete > Complete, High > Medium > Low, Deadline nearest first
		sort.SliceStable(tasks, func(i, j int) bool {
			if tasks[i].Completed != tasks[j].Completed {
				return !tasks[i].Completed // incomplete first
			}
			if tasks[i].Priority != tasks[j].Priority {
				return tasks[i].Priority > tasks[j].Priority // high first
			}
			if tasks[i].Deadline != nil && tasks[j].Deadline != nil {
				return tasks[i].Deadline.Before(*tasks[j].Deadline)
			}
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		})

		if len(tasks) == 0 {
			fmt.Println("📭 No tasks found.")
			return
		}

		for _, t := range tasks {
			status := "❌"
			if t.Completed {
				status = "✅"
			}

			// Deadline display
			var dl string
			if t.Deadline != nil {
				dl = t.Deadline.Format("2006-01-02")
			} else {
				dl = "No deadline"
			}

			fmt.Printf("[%s] %s (%s)\n", status, t.Title, t.Priority.String())
			fmt.Printf("🗓  %s\n", dl)
			if t.Description != "" {
				desc := strings.Split(t.Description, "\n")[0]
				fmt.Printf("📝 %s\n", desc)
			}
			fmt.Printf("🧷 ID: %s\n", t.ID)
			fmt.Println(strings.Repeat("-", 40))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
