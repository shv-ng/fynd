package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shv-ng/fynd/app"
	"github.com/shv-ng/fynd/server"
	"github.com/shv-ng/fynd/store"
	"github.com/shv-ng/fynd/utils"
	"github.com/spf13/cobra"
)

var ctx app.Context

var rootCmd = &cobra.Command{
	Use:   "fynd [query]",
	Short: "Fynd is a fast local search engine for your filesystem.",
	Long: `Fynd lets you search through your local files using full-text indexing, 
offering fast, flexible search with support for filters, and inline query syntax.

You can specify search options using CLI flags or embed them directly into the search query.

Inline Query Format (overrides CLI flags):
  top:N          Show top N results (-1 for all)
  ext:log,txt    Limit results to given file extensions
  (Other text)   Treated as the search query

Examples:

  fynd "top:5;ext:md,txt;search term"
  fynd "error handling" --top=3 --ext=log

Use 'fynd sync' to scan directories, rebuild the index, and keep your data up to date.`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		q := strings.Join(args, " ")
		if cmd.Flag("sync").Changed {
			server.Sync(ctx)
		}
		top, err := strconv.Atoi(cmd.Flag("top").Value.String())
		if err != nil {
			log.Fatalf("Invalid top value: %v\n", err)
		}

		ext := []string{}
		if cmd.Flag("ext").Changed {
			ext = strings.Split(cmd.Flag("ext").Value.String(), ",")
		}
		server.Find(ctx, q, server.QueryOptions{
			Top:   top,
			Ext:   ext,
			Query: []string{},
		})
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get homedir: %v", err)
	}
	s, err := utils.LoadYAMLSettings(filepath.Join(home, ".config", "fynd", "config.yml"))
	if err != nil {
		log.Fatalln(err)
	}

	db, err := store.InitDB(s.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	ctx = app.Context{
		DB:      db,
		Setting: s,
	}

	rootCmd.PersistentFlags().Int("top", 0,
		"Show top N results. Use -1 for all. Can be overridden by 'top:N' in query.")

	rootCmd.PersistentFlags().String("ext", "",
		"Comma-separated file extensions to include (e.g. txt,md). Overridden by 'ext:xyz' in query.")

	rootCmd.PersistentFlags().Bool("sync", false,
		"Shortcut to trigger a sync before searching. Same as 'fynd sync'.")
}
