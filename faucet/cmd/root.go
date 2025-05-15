package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"faucet"

	"github.com/spf13/cobra"
)

var (
	dbPath string
	client *faucet.Client
)

var rootCmd = &cobra.Command{
	Use:   "faucet",
	Short: "A Bitcoin testnet faucet CLI",
	Long: `A CLI tool for managing Bitcoin testnet addresses and transactions.
It provides commands to generate addresses, fund them, and transfer funds between addresses.`,
}

// GetDBPath returns the path to the database file
func GetDBPath() string {
	return dbPath
}

// GetClient returns the Bitcoin RPC client
func GetClient() *faucet.Client {
	return client
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	dbPath = filepath.Join(homeDir, ".faucet", "addresses.db")
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	// Initialize Bitcoin RPC client
	client = faucet.NewClient("http://localhost:18332", "rpcuser", "rpcpass")

	// Register commands
	registerCommands()
}

func registerCommands() {
	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newFundCmd())
	rootCmd.AddCommand(newTransferCmd())
	rootCmd.AddCommand(newDeleteCmd())
} 