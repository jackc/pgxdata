package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

func main() {
	cmdInit := &cobra.Command{
		Use:   "init NAME",
		Short: "Initialize a new data package",
		Run:   initCmd,
	}

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Print version and exit",
		Run:   versionCmd,
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdInit)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.Execute()
}

func versionCmd(cmd *cobra.Command, args []string) {
	fmt.Println(VERSION)
	os.Exit(0)
}
