package main

import (
	"github.com/andrewpmartinez/grid/dump"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
}

func init() {
	parseCmd := &cobra.Command{
		Use:  "parse <file>",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			absPath, err := filepath.Abs(args[0])
			if err != nil {
				logrus.Error("could not find the absolute path for file [%s]", args[0])
				os.Exit(1)
			}

			outDump, err := dump.ParseFile(absPath, logrus.New())

			if err != nil {
				logrus.Error("error parsing file [%s]: %v", absPath, err)
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
