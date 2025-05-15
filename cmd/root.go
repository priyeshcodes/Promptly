// cmd/root.go
package cmd

import (
	//"fmt"
	//"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Smart CLI Task Manager",
	Long: `Task is a smart CLI-based task management tool.
It supports priorities, deadlines, dependencies, markdown descriptions, and more.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
