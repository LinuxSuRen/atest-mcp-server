package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "atest-mcp-server",
	}
	rootCmd.AddCommand(newServerCommand())
	return rootCmd
}
