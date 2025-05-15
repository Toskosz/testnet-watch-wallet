package cmd

import (
	"fmt"

	"github.com/Toskosz/testnet-watch-wallet/faucet/db"

	"github.com/spf13/cobra"
)

func newFundCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fund [alias]",
		Short: "Fund a Bitcoin testnet address by its alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias := args[0]

			db, err := db.NewDB(GetDBPath())
			if err != nil {
				return fmt.Errorf("failed to open database: %v", err)
			}
			defer db.Close()

			// Get address from alias
			address, err := db.GetAddress(alias)
			if err != nil {
				return fmt.Errorf("failed to get address: %v", err)
			}
			if address == "" {
				return fmt.Errorf("alias not found: %s", alias)
			}

			client := GetClient()
			// Fund the address with 0.1 BTC
			result, err := client.Call("sendtoaddress", []interface{}{address, 0.1})
			if err != nil {
				return fmt.Errorf("failed to fund address: %v", err)
			}

			fmt.Printf("Successfully funded address %s (%s) with 0.1 BTC\n", alias, address)
			fmt.Printf("Transaction ID: %s\n", string(result))
			return nil
		},
	}
} 