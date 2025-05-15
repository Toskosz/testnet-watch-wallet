package watcher

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "testnet-watch-wallet",
	Short: "A CLI tool for watching testnet wallets",
	Long: `A CLI tool that helps you monitor and interact with testnet wallets.
It provides various commands to check balances, transactions, and other wallet-related operations.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.testnet-watch-wallet.yaml)")
} 