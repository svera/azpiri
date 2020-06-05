package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	folder string

	rootCmd = &cobra.Command{
		Use:   "cobra",
		Short: "A launching images generator for RetroPie",
		Long: `XXX is a CLI application that will create launching images 
		for games stored in the specified roms folder, using their artwork if available`,
		Run: func(cmd *cobra.Command, args []string) {
			process()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&folder, "folder", "f", ".", "Roms folder to scan")
}
