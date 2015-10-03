package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	fmt.Println("hello")

	cmdInit := &cobra.Command{
		Use:   "init NAME",
		Short: "Initialize a new data package",
		Run:   initCmd,
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdInit)
	rootCmd.Execute()
}
