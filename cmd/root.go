package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	romsFolder        string
	backgroundsFolder string
	foregroundsFolder string

	rootCmd = &cobra.Command{
		Use:   "azpiri",
		Short: "A launching images generator for RetroPie",
		Long: `Azpiri is a CLI application that will create launching images 
		for games stored in a roms folder, using PNG images
		with the same name as the ROM in the specified artwork folders`,
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
	rootCmd.PersistentFlags().StringVarP(&romsFolder, "roms", "r", ".", "Roms folder to scan")
	romsFolder = filepath.Clean(romsFolder)
	rootCmd.PersistentFlags().StringVarP(&backgroundsFolder, "backgrounds", "b", romsFolder+"/media/screenshots/", "Background images folder to scan")
	rootCmd.PersistentFlags().StringVarP(&foregroundsFolder, "foregrounds", "f", romsFolder+"/media/marquees/", "Foreground images folder to scan")

	backgroundsFolder = filepath.Clean(backgroundsFolder)
	foregroundsFolder = filepath.Clean(foregroundsFolder)
}
