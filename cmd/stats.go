// cmd/stats.go
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/priyeshcodes/smart-task-cli/internal/task"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View task statistics",
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

		totalTasks := len(tasks)
		completed := 0
		pending := 0
		for _, t := range tasks {
			if t.Completed {
				completed++
			} else {
				pending++
			}
		}

		// Calculate Task completion rate (percentage)
		completionRate := float64(completed) / float64(totalTasks) * 100

		// Display task statistics
		fmt.Printf("📊 Task Statistics\n")
		fmt.Printf("--------------------\n")
		fmt.Printf("Total tasks: %d\n", totalTasks)
		fmt.Printf("Completed: %d\n", completed)
		fmt.Printf("Pending: %d\n", pending)
		fmt.Printf("Completion rate: %.2f%%\n", completionRate)

		// Optional: Show weekly task completion trend
		showWeeklyTrend(tasks)
	},
}

func showWeeklyTrend(tasks []task.Task) {
	// Filter tasks completed within the last week
	lastWeek := time.Now().AddDate(0, 0, -7)
	var completedThisWeek int
	for _, t := range tasks {
		if t.Completed && t.CompletedAt.After(lastWeek) {
			completedThisWeek++
		}
	}

	fmt.Printf("\n📅 Tasks completed this week: %d\n", completedThisWeek)
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
