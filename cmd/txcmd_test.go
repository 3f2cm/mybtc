package cmd_test

import (
	"bytes"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/3f2cm/mybtc/cli"
	"github.com/3f2cm/mybtc/cmd"
)

func Test_newTxGenerateCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		inputFile  string
		wantTxFile string
		err        bool
	}{
		{
			name:       "sample 1",
			args:       []string{"tx", "generate"},
			inputFile:  "sample_1_input.json",
			wantTxFile: "sample_1_tx.txt",
			err:        false,
		},
		{
			name:       "sample 2",
			args:       []string{"tx", "generate"},
			inputFile:  "sample_2_input.json",
			wantTxFile: "sample_2_tx.txt",
			err:        false,
		},
	}
	for _, tt := range tests {
		input, err := os.ReadFile(path.Join("test_data", tt.inputFile))
		if err != nil {
			t.Fatalf("couldn't read the input file %s: %s", tt.inputFile, err)
		}

		want, err := os.ReadFile(path.Join("test_data", tt.wantTxFile))
		if err != nil {
			t.Fatalf("couldn't read the input file %s: %s", tt.inputFile, err)
		}

		stdin := strings.NewReader(string(input))
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		rootCmd := cmd.NewRootCmd(&cli.Env{
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
		}, tt.args)

		t.Run(tt.name, func(t *testing.T) {
			err := rootCmd.Execute()

			if stdout.String() != string(want) {
				t.Errorf("mybtc tx generate returned %s, want %s", stdout, want)
			}
			if (err != nil) != tt.err {
				t.Errorf("command failed unexpectedly: %s", err)
			}
		})
	}
}
