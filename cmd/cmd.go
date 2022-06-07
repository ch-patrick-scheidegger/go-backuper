package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ch-patrick-scheidegger/go-backuper/backuper"
)

// https://github.com/codeedu/golang-cobra-example
var rootCmd = &cobra.Command{
	Use:   "backuper",
	Short: "Copies all files and sub dirs from the source dir to the destination dir",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Must provide src and dst path")
			return
		}
		src := formatPathParam(args[0])
		dst := formatPathParam(args[1])
		err := validatePathParam(src)
		if err != nil {
			fmt.Println("1. ", err)
			return
		}
		err = validatePathParam(dst)
		if err != nil {
			fmt.Println("2. ", err)
			return
		}

		backuper.BackupDirectory(src, dst)
	},
}

func formatPathParam(parameter string) string {
	parameter = strings.ReplaceAll(parameter, "\\", "/")
	if !strings.HasSuffix(parameter, "/") {
		parameter = parameter + "/"
	}
	return parameter
}

func validatePathParam(parameter string) error {
	if len(parameter) == 0 {
		return errors.New("parameter must be set")
	}
	if _, err := os.Stat(parameter); os.IsNotExist(err) {
		return errors.New("parameter must point to an existing directory")
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// Execute is called by the main function. It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
