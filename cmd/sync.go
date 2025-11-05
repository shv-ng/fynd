package cmd

import (
	"github.com/shv-ng/fynd/server"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Rebuild the search index from the latest filesystem content.",
	Long: `The 'sync' command scans the root directory, processes file content and metadata, 
adds new or modified files to the index, and removes data for deleted files.

Use this to keep your local search index up to date with the actual filesystem state.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Sync(ctx)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
