package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Toskosz/testnet-watch-wallet/faucet/nodeClient"

	"github.com/spf13/cobra"
)

var (
	dbPath string
	client *nodeClient.Client
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
func GetClient() *nodeClient.Client {
	return client
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Get configuration from environment variables
	dbPath = os.Getenv("DB_PATH")
	if dbPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		dbPath = filepath.Join(homeDir, ".faucet", "addresses.db")
	}
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	// Get RPC configuration from environment variables
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://localhost:18332"
	}

	rpcUser := os.Getenv("RPC_USER")
	if rpcUser == "" {
		rpcUser = "rpcuser"
	}

	rpcPass := os.Getenv("RPC_PASS")
	if rpcPass == "" {
		rpcPass = "rpcpass"
	}

	// Initialize Bitcoin RPC client
	client = nodeClient.NewClient(rpcURL, rpcUser, rpcPass)

	// Register commands
	registerCommands()
}

func registerCommands() {
	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newFundCmd())
	rootCmd.AddCommand(newTransferCmd())
	rootCmd.AddCommand(newDeleteCmd())
} 