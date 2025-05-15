// cmd/edit.go
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/priyeshcodes/smart-task-cli/internal/task"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [taskID]",
	Short: "Edit a task's title, description, priority, or deadline",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := args[0]

		store, err := task.NewTaskStore("tasks.db")
		if err != nil {
			log.Fatalf("Failed to open task store: %v", err)
		}
		defer store.Close()

		t, err := store.GetTaskByID(taskID)
		if err != nil {
			log.Fatalf("Task not found: %v", err)
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Editing task: %s\n", t.Title)

		// Title
		fmt.Printf("Title [%s]: ", t.Title)
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)
		if title != "" {
			t.Title = title
		}

		// Description
		fmt.Printf("Description [%s]: ", t.Description)
		desc, _ := reader.ReadString('\n')
		desc = strings.TrimSpace(desc)
		if desc != "" {
			t.Description = desc
		}

		// Priority
		fmt.Printf("Priority (1-5) [%d]: ", t.Priority)
		priorityStr, _ := reader.ReadString('\n')
		priorityStr = strings.TrimSpace(priorityStr)
		if priorityStr != "" {
			p, err := strconv.Atoi(priorityStr)
			if err == nil && p >= 1 && p <= 5 {
				t.Priority = task.PriorityLevel(p)
			} else {
				fmt.Println("⚠️ Invalid priority. Keeping existing value.")
			}
		}

		// Deadline
		fmt.Printf("Deadline (YYYY-MM-DD HH:MM) [%s]: ", t.Deadline.Format("2006-01-02 15:04"))
		deadlineStr, _ := reader.ReadString('\n')
		deadlineStr = strings.TrimSpace(deadlineStr)
		if deadlineStr != "" {
			parsed, err := time.Parse("2006-01-02 15:04", deadlineStr)
			if err == nil {
				t.Deadline = &parsed
			} else {
				fmt.Println("⚠️ Invalid deadline format. Keeping existing value.")
			}
		}

		// Save changes
		err = store.UpdateTask(t)
		if err != nil {
			log.Fatalf("Failed to update task: %v", err)
		}

		fmt.Println("✅ Task updated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
