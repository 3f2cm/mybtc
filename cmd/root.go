/*
Package cmd provides the implementations of the main command and each subcommands
*/
package cmd

import (
	"os"

	"github.com/3f2cm/mybtc/cli"
	"github.com/spf13/cobra"
)

// NewRootCmd will create the root command of mybtc.
func NewRootCmd(env *cli.Env, args []string) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "mybtc",
		Short: "A sample implementation of a CLI for Bitcoin with btcsuite",
		Long: `mybtc is a sample toy implementation of a CLI to manipulate
Bitcoin wallets, addresses, transactions and so on
to understand how Bitcoin and btcsuite packages work.

cf. https://github.com/btcsuite/
    `,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
	rootCmd.SetArgs(args)
	rootCmd.SetIn(env.Stdin)
	rootCmd.SetOut(env.Stdout)
	rootCmd.SetErr(env.Stderr)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(newWIFCmd(env))
	rootCmd.AddCommand(newTxCmd())

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(env *cli.Env, args []string) {
	rootCmd := NewRootCmd(env, args)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

/*
func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mybtc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd := NewRootCmd(cl)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
*/
