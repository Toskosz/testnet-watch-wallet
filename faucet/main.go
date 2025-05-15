package faucet

import (
	"fmt"
	"os"

	"github.com/Toskosz/testnet-watch-wallet/faucet/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}