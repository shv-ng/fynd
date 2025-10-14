package cmd

import (
	"github.com/shv-ng/fynd/server"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
