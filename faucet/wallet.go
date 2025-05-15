package faucet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func generatePublicKey() string {
	secp256k1 := elliptic.P256()

	privKey, err := ecdsa.GenerateKey(secp256k1, rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(privKey.Public())
	if err != nil {
		log.Fatalf("Failed to marshal public key: %v", err)
	}

	return hex.EncodeToString(pubKeyBytes)
}

func hashPublicKey(pubKey string) string {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		log.Fatalf("Failed to decode public key: %v", err)
	}

	sha256Hash := sha256.Sum256(pubKeyBytes)
	ripemd160 := ripemd160.New()
	ripemd160.Write(sha256Hash[:])
	hashedPubKey := ripemd160.Sum(nil)
	return hex.EncodeToString(hashedPubKey)
}

func createBitcoinAddress(hashedPubKey string) string {
	hashedPubKeyBytes, err := hex.DecodeString(hashedPubKey)
	if err != nil {
		log.Fatalf("Failed to decode hashed public key: %v", err)
	}

	versionBuffer := []byte{0x6f}
	payload := append(versionBuffer, hashedPubKeyBytes...)

	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]

	payload = append(payload, checksum...)
	
	address := base58.Encode(payload)

	return address
}

func IsValidBitcoinAddress(address string) bool {
	_, err := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	return err == nil
}

func GenerateNewAddress() string {
	pubKey := generatePublicKey()

	hashedPubKey := hashPublicKey(pubKey)

	address := createBitcoinAddress(hashedPubKey)

	if !IsValidBitcoinAddress(address) {
		log.Fatalf("Invalid Bitcoin address: %s", address)
	}

	return address
}