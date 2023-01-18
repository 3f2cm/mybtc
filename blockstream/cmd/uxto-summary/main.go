/*
Package uxto-summary provides a command that retrieves UXTO info
about given TestNet3 Bitcoin address
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/3f2cm/mybtc/blockstream/bs"
)

func main() {
	if len(os.Args) != 1+1 {
		if _, err := os.Stdout.WriteString(fmt.Sprintf("Usage: %s <Bitcoind address>\n", os.Args[0])); err != nil {
			log.Fatalf("Couldn't write strings to Stdou: %s", err)
		}

		os.Exit(1)
	}

	utxos, err := bs.GetUTXOWithScriptPubKey(os.Args[1])
	if err != nil {
		log.Fatalf("Couldn't get all the info: %s", err)
	}

	b, err := json.Marshal(utxos)
	if err != nil {
		log.Fatalf("Couldn't serialize Vin info: %s", err)
	}

	fmt.Printf("%s\n", b)
}
