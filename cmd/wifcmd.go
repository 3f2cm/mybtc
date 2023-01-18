package cmd

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/3f2cm/mybtc/cli"
	"github.com/3f2cm/mybtc/wif"
	"github.com/spf13/cobra"
)

// wifCmd represents the wif command.
func newWIFCmd(env *cli.Env) *cobra.Command {
	wifCmd := &cobra.Command{
		Use:   "wif",
		Short: "wif manipulates WIF",
		Long: `wif command handles WIF (Wallet Import Format) like
generating WIF, parsing the contents, extract its pubkey address, and so on.`,
	}

	// register subcommands
	wifCmd.AddCommand(newWIFGenerateCmd(env))
	wifCmd.AddCommand(newWIFAddressCmd())

	return wifCmd
}

func newWIFGenerateCmd(env *cli.Env) *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generates a new WIF",
		Long:  `generates a new WIF by generating a new private key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			w, err := wif.New(env.Rand)
			if err != nil {
				return fmt.Errorf("couldn't generate a WIF: %w", err)
			}

			cmd.Println(w)

			return nil
		},
		SilenceUsage: true,
	}

	return generateCmd
}

func newWIFAddressCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "address",
		Short: "converts given WIFs to TestNet3 Address",
		Long:  `receives WIFs from STDIN and converts them to addresses`,
		RunE: func(cmd *cobra.Command, args []string) error {
			buf := bufio.NewReader(cmd.InOrStdin())

			i := uint64(0)
			eof := false
			for {
				i++

				s, err := buf.ReadString('\n')
				if err == io.EOF {
					eof = true
				} else if err != nil {
					return fmt.Errorf("read error happened: %w", err)
				}

				s = strings.TrimSpace(s)
				if len(s) > 0 {
					a, err := wif.ExtractAddr(s)
					if err != nil {
						return fmt.Errorf("couldn't extract an address from line %d: %w", i, err)
					}

					cmd.Println(a)
				}

				if eof {
					break
				}
			}

			return nil
		},
		SilenceUsage: true,
	}

	return generateCmd
}
