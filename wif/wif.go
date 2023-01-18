/*
Package wif handles WIFs

- New generates new WIF
- ExtractAddress extract the address from WIF
*/
package wif

import (
	"crypto/ecdsa"
	"fmt"
	"io"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// New will construct a new WIF object with the given random generator.
func New(r io.Reader) (string, error) {
	key, err := ecdsa.GenerateKey(btcec.S256(), r)
	if err != nil {
		return "", fmt.Errorf("couldn't generate a private key: %w", err)
	}

	pk := secp256k1.PrivKeyFromBytes(key.D.Bytes())

	wif, err := btcutil.NewWIF(pk, &chaincfg.TestNet3Params, false)
	if err != nil {
		return "", fmt.Errorf("couldn't create a wif key: %w", err)
	}

	return wif.String(), nil
}

// ExtractAddr extracts the address from the given WIF.
func ExtractAddr(s string) (string, error) {
	wif, err := btcutil.DecodeWIF(s)
	if err != nil {
		return "", fmt.Errorf("failed to parse given WIF: %w", err)
	}

	pubkey := wif.PrivKey.PubKey().SerializeUncompressed()

	addr, err := btcutil.NewAddressPubKey(pubkey, &chaincfg.TestNet3Params)
	if err != nil {
		return "", fmt.Errorf("failed to generate an address with pubkey %s: %w", pubkey, err)
	}

	return addr.EncodeAddress(), nil
}
