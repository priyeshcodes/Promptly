package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/priyeshcodes/smart-task-cli/internal/task"
	"github.com/spf13/cobra"
)

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Show tasks that are due in the next 24 hours",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := task.NewTaskStore("tasks.db")
		if err != nil {
			log.Fatalf("Failed to open task store: %v", err)
		}
		defer store.Close()

		tasks, err := store.GetAllTasks()
		if err != nil {
			log.Fatalf("Failed to retrieve tasks: %v", err)
		}

		now := time.Now()
		soon := now.Add(24 * time.Hour)
		found := false

		fmt.Println("🔔 Upcoming Task Deadlines (Next 24 hours):")
		for _, t := range tasks {
			if !t.Completed && !t.Deadline.IsZero() && t.Deadline.After(now) && t.Deadline.Before(soon) {
				fmt.Printf("• [%s] %s (Due: %s)\n", t.ID, t.Title, t.Deadline.Format("Mon, 02 Jan 2006 15:04"))
				found = true
			}
		}

		if !found {
			fmt.Println("🎉 No tasks due in the next 24 hours.")
		}
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)
}
