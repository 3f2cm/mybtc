package cmd

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/3f2cm/mybtc/tx"
	"github.com/spf13/cobra"
)

// newTxCmd generates command for tx subcommand.
func newTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "tx manipulates transactions",
		Long:  `tx command handles Bitcoin transactions like generating transactions`,
	}

	// register subcommands
	txCmd.AddCommand(newTxGenerateCmd())

	return txCmd
}

func newTxGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "converts given WIFs to TestNet3 Address",
		Long:  `receives WIFs from STDIN and converts them to addresses`,
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("couldn't read the input: %w", err)
			}

			t, err := tx.Generate(input)
			if err != nil {
				return fmt.Errorf("couldn't generate signed transaction from input: %w", err)
			}

			h := hex.EncodeToString(t)
			cmd.Println(h)

			return nil
		},
		SilenceUsage: true,
	}

	return generateCmd
}
