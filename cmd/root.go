package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configPath string

var root = &cobra.Command{
	Use: "redirecter",
	Short: `Go Module Redirecter is an HTTP server that redirects Go modules from custom domains
to public version control servers (like GitHub) to provide nice import paths for
developers and redirects to documentation for the same import paths.`,
}

func init() {
	root.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Config file (default is $HOME/.redirecter.yaml)")

	root.AddCommand(serve)
}

// Execute starts the program
func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
