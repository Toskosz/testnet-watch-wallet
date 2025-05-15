package cmd

import (
	"fmt"

	"github.com/Toskosz/testnet-watch-wallet/faucet/db"
	"github.com/Toskosz/testnet-watch-wallet/faucet/wallet"

	"github.com/spf13/cobra"
)

func newGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate [alias]",
		Short: "Generate a new Bitcoin testnet address with an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias := args[0]

			db, err := db.NewDB(GetDBPath())
			if err != nil {
				return fmt.Errorf("failed to open database: %v", err)
			}
			defer db.Close()

			// Check if alias already exists
			address, err := db.GetAddress(alias)
			if err != nil {
				return fmt.Errorf("failed to check alias: %v", err)
			}
			if address != "" {
				return fmt.Errorf("alias already exists: %s", alias)
			}

			newAddress := wallet.GenerateNewAddress()
			if err := db.StoreAddress(newAddress, alias); err != nil {
				return fmt.Errorf("failed to store address: %v", err)
			}

			fmt.Printf("Generated new address: %s with alias: %s\n", address, alias)
			return nil
		},
	}
} 