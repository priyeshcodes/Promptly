// main.go
package main

import (
	"log"
	"github.com/priyeshcodes/smart-task-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
