package GNBS

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	Execute()
}

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "gnbs",
		Short: "GNBS is a scripting language based on C++ and Golang languages",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
