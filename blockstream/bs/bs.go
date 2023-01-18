/*
Package bs will be CLI to access to blockstream TestNet3
*/
package bs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// TxStatus expresses the status of the transaction.
type TxStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int64  `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   uint64 `json:"block_time"`
}

// UTXO expresses an unspent transaction output.
type UTXO struct {
	TxID   string   `json:"txid"`
	Idx    uint32   `json:"vout"`
	Status TxStatus `json:"status"`
	Value  uint64   `json:"value"`
}

// GetUTXO retrieves a list of UTXO of the given address.
func GetUTXO(a string) ([]UTXO, error) {
	target := fmt.Sprintf("https://blockstream.info/testnet/api/address/%s/utxo", a)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, target, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("couldn't build an HTTP request to '%s': %w", target, err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't get UTXO about %s: %w", a, err)
	}
	//nolint:errcheck // nothing to do at error
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read body of the response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request for %s failed: code: %d, body: %s", a, res.StatusCode, b)
	}

	us := []UTXO{}
	if err := json.Unmarshal(b, &us); err != nil {
		return nil, fmt.Errorf("could't parse retrieved UTXO output: %w", err)
	}

	return us, nil
}

// Vin expresses vin item in blockstream responses.
type Vin struct {
	TxID       string `json:"txid"`
	Vout       uint32 `json:"vout"`
	PrevOut    Vout   `json:"prevout"`
	IsCoinbase bool   `json:"is_coinbase"`
	Sequence   uint64 `json:"sequence"`
}

// Vout expresses vout item in blockstream responses.
type Vout struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               uint64 `json:"value"`
}

// Tx expresses transactions in blockstream responses.
type Tx struct {
	TxID     string   `json:"txid"`
	Version  int      `json:"version"`
	LockTime int      `json:"locktime"`
	Vin      []Vin    `json:"vin"`
	Vout     []Vout   `json:"vout"`
	Fee      uint64   `json:"fee"`
	Status   TxStatus `json:"status"`
}

// GetTx retrieves transaction data associated with the given transaction ID.
func GetTx(txid string) (*Tx, error) {
	target := fmt.Sprintf("https://blockstream.info/testnet/api/tx/%s", txid)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, target, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("couldn't build an HTTP request to '%s': %w", target, err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("couldn't get Tx of %s: %w", txid, err)
	}

	//nolint:errcheck // nothing to do at error
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read body of the response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request for %s failed: code: %d, body: %s", txid, res.StatusCode, b)
	}

	tx := &Tx{}
	if err := json.Unmarshal(b, tx); err != nil {
		return nil, fmt.Errorf("could't parse retrieved UTXO output: %w", err)
	}

	return tx, nil
}

// VinSummary expresses a summary of Vin.
type VinSummary struct {
	TxID         string `json:"txid"`
	Vout         uint32 `json:"vout"`
	ScriptPubKey string `json:"scriptpubkey"`
	Value        uint64 `json:"value"`
}

// GetUTXOWithScriptPubKey returns a list of summary of UTXO with ScriptPubKey.
func GetUTXOWithScriptPubKey(a string) ([]VinSummary, error) {
	utxos, err := GetUTXO(a)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Initialize the output
	vinInput := []VinSummary{}

	// Fulfill vinInput
	for _, utxo := range utxos {
		tx, err := GetTx(utxo.TxID)
		if err != nil {
			return nil, fmt.Errorf("couldn't get tx of %s: %w", utxo.TxID, err)
		}

		for i, vout := range tx.Vout {
			if vout.ScriptPubKeyAddress == a {
				vinInput = append(vinInput, VinSummary{
					TxID:         tx.TxID,
					Vout:         uint32(i),
					ScriptPubKey: vout.ScriptPubKey,
					Value:        vout.Value,
				})
			}
		}
	}

	return vinInput, nil
}
