package cmd

import (
	"fmt"

	"github.com/Toskosz/testnet-watch-wallet/faucet/db"

	"github.com/spf13/cobra"
)

func newTransferCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [from_alias] [to_alias]",
		Short: "Transfer BTC between addresses using their aliases",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromAlias := args[0]
			toAlias := args[1]

			db, err := db.NewDB(GetDBPath())
			if err != nil {
				return fmt.Errorf("failed to open database: %v", err)
			}
			defer db.Close()

			// Get addresses from aliases
			fromAddr, err := db.GetAddress(fromAlias)
			if err != nil {
				return fmt.Errorf("failed to get 'from' address: %v", err)
			}
			if fromAddr == "" {
				return fmt.Errorf("alias not found: %s", fromAlias)
			}

			toAddr, err := db.GetAddress(toAlias)
			if err != nil {
				return fmt.Errorf("failed to get 'to' address: %v", err)
			}
			if toAddr == "" {
				return fmt.Errorf("alias not found: %s", toAlias)
			}

			client := GetClient()
			// Transfer 0.05 BTC
			result, err := client.Call("sendtoaddress", []interface{}{toAddr, 0.05})
			if err != nil {
				return fmt.Errorf("failed to transfer funds: %v", err)
			}

			fmt.Printf("Successfully transferred 0.05 BTC from %s (%s) to %s (%s)\n", 
				fromAlias, fromAddr, toAlias, toAddr)
			fmt.Printf("Transaction ID: %s\n", string(result))
			return nil
		},
	}
} 