package main

import (
	"fmt"
	"testnet-watch-wallet/src/wallet"
)

func main() {

	address := wallet.GenerateNewAddress()
	fmt.Println(address)

	// cmd.Execute()
}
