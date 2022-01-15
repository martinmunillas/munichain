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

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
