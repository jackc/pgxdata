package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

var appArgs struct {
	templatesPath string
}

func main() {
	cmdInit := &cobra.Command{
		Use:   "init NAME",
		Short: "Initialize a new data package",
		Run:   initCmd,
	}
	cmdInit.Flags().StringVarP(&appArgs.templatesPath, "templates", "t", "", "Templates directory to read from")

	cmdGenerate := &cobra.Command{
		Use:   "generate",
		Short: "Build",
		Run:   generateCmd,
	}
	cmdGenerate.Flags().StringVarP(&appArgs.templatesPath, "templates", "t", "", "Templates directory to read from")

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Print version and exit",
		Run:   versionCmd,
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdInit)
	rootCmd.AddCommand(cmdGenerate)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.Execute()
}

func versionCmd(cmd *cobra.Command, args []string) {
	fmt.Println(VERSION)
	os.Exit(0)
}
