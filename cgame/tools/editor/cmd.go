package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "editor",
	Long: "editor is the cgame map editor",
}

func init() {
	rootCmd.AddCommand(newCmd)
}

var (
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "creates a new map",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			newMap()
		},
	}
	newMapW, newMapH int
)

func init() {
	newCmd.Flags().IntVar(&newMapW, "w", 150, "new map's width in characters")
	newCmd.Flags().IntVar(&newMapH, "h", 30, "new map's height in characters")
}

func newMap() {
	(&editor{}).main(editorModeNew)
}
