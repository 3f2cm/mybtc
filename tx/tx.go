/*
Package tx provides functions handling transactions

- Generate generates a new signed transaction from an input
*/
package tx

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// In contains necessary info to establish transaction message's TxIn items.
type In struct {
	TxID         string `json:"txid"`
	Vout         uint32 `json:"vout"`
	ScriptPubKey string `json:"scriptpubkey"`
	WIF          string `json:"wif"`
}

// Out contains necessary info to establish transaction message's TxOut items.
type Out struct {
	Addr  string `json:"addr"`
	Value int64  `json:"value"`
}

// Input expresses an input to generate command that build transaction message with signatures.
type Input struct {
	Ins  []In  `json:"in"`
	Outs []Out `json:"outs"`
}

var (
	errEmptyTxIn  = errors.New("there are no TxIn items in the input")
	errEmptyTxOut = errors.New("there are no TxOut items in the input")
)

// Generate generates a transaction with signatures from given input.
func Generate(b []byte) ([]byte, error) {
	// Deserialize the given input
	var input Input
	if err := json.Unmarshal(b, &input); err != nil {
		return nil, fmt.Errorf("couldn't parse given input: %w", err)
	}

	if len(input.Ins) == 0 {
		return nil, errEmptyTxIn
	}

	if len(input.Outs) == 0 {
		return nil, errEmptyTxOut
	}

	// Initialize msgTx to construct it and insert a signature into its TxIn
	msgTx := wire.NewMsgTx(wire.TxVersion)

	// Construct msgTx.TxOut and msgTx.TxIn
	if err := addOutToTx(msgTx, input.Outs); err != nil {
		return nil, fmt.Errorf("couldn't construct TxOut for msgTx: %w", err)
	}

	if err := addInToTx(msgTx, input.Ins); err != nil {
		return nil, fmt.Errorf("couldn't construct TxIn for msgTx: %w", err)
	}

	// Serialize and hexdump the msgTx with signs
	var signedTx bytes.Buffer
	if err := msgTx.Serialize(&signedTx); err != nil {
		return nil, fmt.Errorf("couldn't serialize the built signed transaction: %w", err)
	}

	return signedTx.Bytes(), nil
}

func addOutToTx(t *wire.MsgTx, outs []Out) error {
	for _, out := range outs {
		addr, err := btcutil.DecodeAddress(out.Addr, &chaincfg.TestNet3Params)
		if err != nil {
			return fmt.Errorf("couldn't decode given address '%s': %w", out.Addr, err)
		}

		payToAddrScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return fmt.Errorf("couldn't generate a script to pay: %w", err)
		}

		txOut := wire.NewTxOut(out.Value, payToAddrScript)
		t.AddTxOut(txOut)
	}

	return nil
}

func addInToTx(t *wire.MsgTx, ins []In) error {
	for _, txin := range ins {
		txidHash, err := chainhash.NewHashFromStr(txin.TxID)
		if err != nil {
			return fmt.Errorf("couldn't create a hash from txid: %w", err)
		}

		outPoint := wire.NewOutPoint(txidHash, txin.Vout)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		t.AddTxIn(txIn)
	}

	// Add a signature to each TxIn in msgTx.TxIn
	for i, txin := range ins {
		// Decoding WIF to use its private key
		wif, err := btcutil.DecodeWIF(txin.WIF)
		if err != nil {
			return fmt.Errorf("couldn't decode given wif: %w", err)
		}

		// Construct a signature
		prevPubKeyScript, err := hex.DecodeString(txin.ScriptPubKey)
		if err != nil {
			return fmt.Errorf("couldn't decode the previous public key script: %w", err)
		}

		signature, err := txscript.SignatureScript(t, i, prevPubKeyScript, txscript.SigHashAll, wif.PrivKey, false)
		if err != nil {
			return fmt.Errorf("couldn't generate a signature for the tx: %w", err)
		}

		// Add the constructed signature to the corresponding TxIn
		t.TxIn[i].SignatureScript = signature
	}

	return nil
}
