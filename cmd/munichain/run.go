package main

import (
	"fmt"
	"os"

	"github.com/martinmunillas/munichain/node"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Runs the node.",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			fmt.Println("Launching munichain node and its HTTP API...")

			err := node.Run(dataDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	addDefaultRequiredFlags(runCmd)

	return runCmd
}
