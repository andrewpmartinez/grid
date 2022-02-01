package main

import (
	"github.com/andrewpmartinez/grid/dump"
	"github.com/andrewpmartinez/grid/gui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use:  "<file>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		absPath, err := filepath.Abs(path)

		if err != nil {
			logrus.Errorf("invalid path, could not transform to absolute path: %s", path)
			os.Exit(1)
		}

		_, err = os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				logrus.Errorf("invalid file, does not exist: %s", absPath)
			} else {
				logrus.Errorf("unexpected error attempting to check file %s: %v", absPath, err)
			}
			os.Exit(1)
		}

		win := gui.NewDumpWindow()
		win.LoadFile(absPath)
		win.Run()
	},
}

func init() {
	parseCmd := &cobra.Command{
		Use:  "parse <file>",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			absPath, err := filepath.Abs(args[0])
			if err != nil {
				logrus.Errorf("could not find the absolute path for file [%s]", args[0])
				os.Exit(1)
			}

			outDump, err := dump.ParseFile(absPath, logrus.New())

			if err != nil {
				logrus.Errorf("error parsing file [%s]: %v", absPath, err)
				os.Exit(1)
			}

			logrus.Infof("parsed %d go routines", len(outDump.Routines))
		},
	}

	rootCmd.AddCommand(parseCmd)

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
	}
}
