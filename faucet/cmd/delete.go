package cmd

import (
	"fmt"

	"faucet"

	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [alias]",
		Short: "Delete a stored address by its alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias := args[0]

			db, err := faucet.NewDB(GetDBPath())
			if err != nil {
				return fmt.Errorf("failed to open database: %v", err)
			}
			defer db.Close()

			// Check if alias exists
			address, err := db.GetAddress(alias)
			if err != nil {
				return fmt.Errorf("failed to check alias: %v", err)
			}
			if address == "" {
				return fmt.Errorf("alias not found: %s", alias)
			}

			if err := db.DeleteAddress(alias); err != nil {
				return fmt.Errorf("failed to delete address: %v", err)
			}

			fmt.Printf("Successfully deleted address with alias: %s\n", alias)
			return nil
		},
	}
} 