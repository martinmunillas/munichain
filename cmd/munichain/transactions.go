package main

import (
	"fmt"
	"os"

	"github.com/martinmunillas/munichain/munichain"
	"github.com/spf13/cobra"
)

const flagFrom = "from"
const flagTo = "to"
const flagValue = "value"

func txCmd() *cobra.Command {
	var txsCmd = &cobra.Command{
		Use:   "tx",
		Short: "Interact with txs (add...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	txsCmd.AddCommand(txAddCmd())

	return txsCmd
}

func txAddCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add",
		Short: "Adds new TX to database.",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)

			tx := munichain.NewTransaction(munichain.Address(from), munichain.Address(to), value)

			state, err := munichain.NewStateFromDisk()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			state.Add(&tx)

			err = state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Println("TX successfully added to the ledger.")
		},
	}
	cmd.Flags().String(flagFrom, "", "From what account to send tokens")
	cmd.MarkFlagRequired(flagFrom)

	cmd.Flags().String(flagTo, "", "To what account to send tokens")
	cmd.MarkFlagRequired(flagTo)

	cmd.Flags().Uint(flagValue, 0, "How many tokens to send")
	cmd.MarkFlagRequired(flagValue)
	return cmd
}
