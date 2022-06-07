package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	_ "github.com/ch-patrick-scheidegger/go-backuper/backuper"
)

var rootCmd = &cobra.Command{
	Use:   "backuper srcdir dstdir",
	Short: "Copies all files and sub dirs from the source dir to the destination dir",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Got params '%v' and '%v'\n", args[0], args[1])
		err := validatePathParam(args[0])
		if err != nil {
			fmt.Println("1. ", err)
			return
		}
		err = validatePathParam(args[1])
		if err != nil {
			fmt.Println("2. ", err)
			return
		}

		//backuper.BackupDirectory(args[0], args[1])
	},
}

func validatePathParam(parameter string) error {
	if len(parameter) == 0 {
		return errors.New("parameter must be set")
	}
	if _, err := os.Stat(parameter); os.IsNotExist(err) {
		return errors.New("parameter must point to an existing directory")
	}
	if !strings.HasSuffix(parameter, "/") {
		return errors.New("parameter must end with a backslash '/'")
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// Execute is called by the main function. It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
