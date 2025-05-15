package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/thiago/testnet-watch-wallet/watcher/nodeClient"
)

type TransactionSummary struct {
	Address     string    `json:"address"`
	Balance     string    `json:"balance"`
	LastTxHash  string    `json:"last_tx_hash,omitempty"`
	LastTxTime  time.Time `json:"last_tx_time,omitempty"`
	LastTxValue string    `json:"last_tx_value,omitempty"`
}

type Watcher struct {
	db     *sql.DB
	client *nodeClient.Client
}

func NewWatcher() *Watcher {
	// Get configuration from environment variables
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/root/.faucet/faucet.db"
	}

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		rpcURL = "http://localhost:18332"
	}

	rpcUser := os.Getenv("RPC_USER")
	if rpcUser == "" {
		rpcUser = "user"
	}

	rpcPass := os.Getenv("RPC_PASS")
	if rpcPass == "" {
		rpcPass = "pass"
	}

	// Connect to database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Connect to Bitcoin testnet node
	client := nodeClient.NewClient(rpcURL, rpcUser, rpcPass)

	return &Watcher{
		db:     db,
		client: client,
	}
}

func (w *Watcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	// Run initial check immediately
	if err := w.checkAddresses(); err != nil {
		log.Printf("Error in initial address check: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.checkAddresses(); err != nil {
				log.Printf("Error checking addresses: %v", err)
			}
		}
	}
}

func (w *Watcher) checkAddresses() error {
	rows, err := w.db.Query("SELECT address FROM addresses")
	if err != nil {
		return fmt.Errorf("failed to query addresses: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var address string
		if err := rows.Scan(&address); err != nil {
			return fmt.Errorf("failed to scan address: %v", err)
		}

		summary, err := w.getAddressSummary(address)
		if err != nil {
			log.Printf("Error getting summary for %s: %v", address, err)
			continue
		}

		jsonData, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			log.Printf("Error marshaling summary for %s: %v", address, err)
			continue
		}

		fmt.Println(string(jsonData))
	}

	return nil
}

func (w *Watcher) getAddressSummary(address string) (*TransactionSummary, error) {
	// Get address balance
	balanceResult, err := w.client.Call("getreceivedbyaddress", []interface{}{address, 0})
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}

	var balance float64
	if err := json.Unmarshal(balanceResult, &balance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal balance: %v", err)
	}

	// Get address transactions
	txResult, err := w.client.Call("searchrawtransactions", []interface{}{address, 0, 1, true, nil})
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}

	var txs []struct {
		TxID    string  `json:"txid"`
		Time    int64   `json:"time"`
		Vout    []struct {
			Value        float64 `json:"value"`
			ScriptPubKey struct {
				Addresses []string `json:"addresses"`
			} `json:"scriptPubKey"`
		} `json:"vout"`
	}

	if err := json.Unmarshal(txResult, &txs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transactions: %v", err)
	}

	summary := &TransactionSummary{
		Address: address,
		Balance: fmt.Sprintf("%.8f", balance),
	}

	if len(txs) > 0 {
		tx := txs[0]
		summary.LastTxHash = tx.TxID
		summary.LastTxTime = time.Unix(tx.Time, 0)

		// Calculate total output value for this address
		var totalValue float64
		for _, out := range tx.Vout {
			for _, outAddr := range out.ScriptPubKey.Addresses {
				if outAddr == address {
					totalValue += out.Value
				}
			}
		}
		summary.LastTxValue = fmt.Sprintf("%.8f", totalValue)
	}

	return summary, nil
} 