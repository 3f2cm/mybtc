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
}

// Out contains necessary info to establish transaction message's TxOut items.
type Out struct {
	Addr  string `json:"addr"`
	Value int64  `json:"value"`
}

// Input expresses an input to generate command that build transaction message with signatures.
type Input struct {
	Ins  []In     `json:"in"`
	Outs []Out    `json:"outs"`
	WIFs []string `json:"wifs"`
}

// wifDB stores WIFs with keys of their public key hashes.
type wifDB map[string]*btcutil.WIF

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

	// Construct msgTx.TxOut
	if err := addOutToTx(msgTx, input.Outs); err != nil {
		return nil, fmt.Errorf("couldn't construct TxOut for msgTx: %w", err)
	}

	// Construct msgTx.TxIn
	if err := addInToTx(msgTx, input.Ins); err != nil {
		return nil, fmt.Errorf("couldn't construct TxIn for msgTx: %w", err)
	}

	// Decode WIFs in the input
	wdb, err := decodeWIFs(input.WIFs)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode WIFs in the input: %w", err)
	}

	// Add signatures to msgTx.TxIn
	if err := updateSignatures(msgTx, input.Ins, wdb); err != nil {
		return nil, fmt.Errorf("couldn't sign msgTx: %w", err)
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

	return nil
}

func decodeWIFs(wifs []string) (wifDB, error) {
	w := make(wifDB)

	for _, encWIF := range wifs {
		wif, err := btcutil.DecodeWIF(encWIF)
		if err != nil {
			return nil, fmt.Errorf("couldn't decode given WIF: %w", err)
		}

		pubKey := wif.SerializePubKey()
		pubKeyHash := btcutil.Hash160(pubKey)
		pubKeyHashEnc := hex.EncodeToString(pubKeyHash)
		w[pubKeyHashEnc] = wif
	}

	return w, nil
}

func updateSignatures(t *wire.MsgTx, ins []In, wdb wifDB) error {
	// Add a signature to each TxIn in msgTx.TxIn
	for i, txin := range ins {
		// Parse the public key script in the input
		prevPubKeyScriptBytes, err := hex.DecodeString(txin.ScriptPubKey)
		if err != nil {
			return fmt.Errorf("couldn't decode the previous public key script: %w", err)
		}

		// Extract the PubKey hash from the public key script
		prevPubKeyScript, err := txscript.ParsePkScript(prevPubKeyScriptBytes)
		if err != nil {
			return fmt.Errorf("couldn't parse the given public key script: %w", err)
		}

		pubKeyHash := extractPubKeyHash(prevPubKeyScript.Script())
		if pubKeyHash == nil {
			return fmt.Errorf("couldn't generate an address from the public key script: %w", err)
		}

		// Find the WIF corresponding to the above PubKey hash from the WIF list in the input
		wif, ok := wdb[hex.EncodeToString(pubKeyHash)]
		if !ok {
			return fmt.Errorf("couldn't find WIF corresponding to the public key script %s in given WIFs", txin.ScriptPubKey)
		}

		// Construct a signature
		signature, err := txscript.SignatureScript(t, i, prevPubKeyScriptBytes, txscript.SigHashAll, wif.PrivKey, false)
		if err != nil {
			return fmt.Errorf("couldn't generate a signature for the tx: %w", err)
		}

		// Add the constructed signature to the corresponding TxIn
		t.TxIn[i].SignatureScript = signature
	}

	return nil
}

// stolen from https://github.com/btcsuite/btcd/blob/1d77730e9a92eafadb9f1ef86d2522c36ca4db38/txscript/standard.go#L154
func extractPubKeyHash(script []byte) []byte {
	if len(script) == 25 &&
		script[0] == txscript.OP_DUP && script[1] == txscript.OP_HASH160 && script[2] == txscript.OP_DATA_20 &&
		script[23] == txscript.OP_EQUALVERIFY && script[24] == txscript.OP_CHECKSIG {
		return script[3:23]
	}

	return nil
}
