package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var cmd = &cobra.Command{
		Use:   "munichain",
		Short: "Munichain CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Munichain CLI")
		},
	}

	cmd.AddCommand(versionCmd)
	cmd.AddCommand(balancesCmd())
	cmd.AddCommand(txCmd())

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}

const flagDataDir = "dataDir"

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(
		flagDataDir,
		"",
		"Data directory",
	)

	cmd.MarkFlagRequired(flagDataDir)
}
