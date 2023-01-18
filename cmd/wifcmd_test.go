package cmd_test

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"

	"github.com/3f2cm/mybtc/cli"
	"github.com/3f2cm/mybtc/cmd"
	"github.com/btcsuite/btcd/btcutil"
)

func Test_newGenerateCmd(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		stdin string
		seed  int64
		wif   string
	}{
		{
			name: "seed = 1",
			args: []string{"wif", "generate"},
			seed: 1,
			wif:  "91kiSsVtKYWBFJePnqk51k9yafKPJBpP52ZDxWc4em2KXtF82B3\n",
		},
		{
			name: "seed = 37",
			args: []string{"wif", "generate"},
			seed: 37,
			wif:  "92K4bgqyq2d2dfoifzqfcjKR2N2vWUd86MmWZp6VHrTBPQiHqUP\n",
		},
	}
	for _, tt := range tests {
		stdin := strings.NewReader(tt.stdin)
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		//nolint:gosec // It's for test, and we need specified seeds
		rand := rand.New(rand.NewSource(tt.seed))

		rootCmd := cmd.NewRootCmd(&cli.Env{
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
			Rand:   rand,
		}, tt.args)

		t.Run(tt.name, func(t *testing.T) {
			if err := rootCmd.Execute(); err != nil {
				t.Errorf("Error happened during execution: %s", err)
			}

			wif := stdout.String()
			if tt.wif != stdout.String() {
				t.Errorf("mybtc wifi generate returned %s, want %s", wif, tt.wif)
			}
		})
	}
}

func Test_newGenerateCmdGenerateRandomly(t *testing.T) {
	attempts := 100000
	args := []string{"wif", "generate"}
	rand := crand.Reader

	t.Run(fmt.Sprintf("generate %d WIFs", attempts), func(t *testing.T) {
		var mu sync.Mutex
		wg := sync.WaitGroup{}
		c := make(chan string)
		defer close(c)
		fatals := []string{}

		go func() {
			for {
				msg := <-c
				if msg == "" {
					break
				}
				fatals = append(fatals, msg)
			}
		}()

		m := make(map[string]struct{})
		for i := 1; i <= attempts; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()
				stdout := &bytes.Buffer{}
				stderr := &bytes.Buffer{}
				rootCmd := cmd.NewRootCmd(&cli.Env{
					Stdout: stdout,
					Stderr: stderr,
					Rand:   rand,
				}, args)
				if err := rootCmd.Execute(); err != nil {
					c <- fmt.Sprintf("Execution failure happened: %s", err)

					return
				}

				wif := strings.TrimSpace(stdout.String())
				if _, err := btcutil.DecodeWIF(wif); err != nil {
					c <- fmt.Sprintf("invalid WIF was generated at the %d-th attempt: %s: %s", i, wif, err)

					return
				}

				mu.Lock()
				m[wif] = struct{}{}
				mu.Unlock()
			}(i)
		}

		wg.Wait()
		c <- ""

		if len(fatals) > 0 {
			t.Fatalf("there were %d fatal errors during WIF generations: %s", len(fatals), strings.Join(fatals[:10], ","))
		}

		l := len(m)
		if l != attempts {
			t.Errorf("mybtc wifi generate generated only %d WIFs with %d attempts", l, attempts)
		}
	})
}

func Test_newAddressCmd(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		stdin  string
		stdout string
		stderr string
		isErr  bool
	}{
		{
			name:   "one line",
			args:   []string{"wif", "address"},
			stdin:  "91kiSsVtKYWBFJePnqk51k9yafKPJBpP52ZDxWc4em2KXtF82B3\n",
			stdout: "mz6u8QbVrdChQovwCb9AirW1Wm99fUJ7ko\n",
			stderr: "",
			isErr:  false,
		},
		{
			name:   "two lines",
			args:   []string{"wif", "address"},
			stdin:  "91kiSsVtKYWBFJePnqk51k9yafKPJBpP52ZDxWc4em2KXtF82B3\n92K4bgqyq2d2dfoifzqfcjKR2N2vWUd86MmWZp6VHrTBPQiHqUP\n",
			stdout: "mz6u8QbVrdChQovwCb9AirW1Wm99fUJ7ko\nmkYBnkquUXqFvCA4NJDDGowjksHfJzQRMf\n",
			stderr: "",
			isErr:  false,
		},
		{
			name:   "null input",
			args:   []string{"wif", "address"},
			stdin:  "",
			stdout: "",
			stderr: "",
			isErr:  false,
		},
		{
			name:   "broken input",
			args:   []string{"wif", "address"},
			stdin:  "aaa",
			stdout: "",
			stderr: "Error: couldn't extract an address from line 1",
			isErr:  true,
		},
		{
			name: "broken input in the middle",
			args: []string{"wif", "address"},
			stdin: "91kiSsVtKYWBFJePnqk51k9yafKPJBpP52ZDxWc4em2KXtF82B3\n" +
				"broken\n92K4bgqyq2d2dfoifzqfcjKR2N2vWUd86MmWZp6VHrTBPQiHqUP",
			stdout: "mz6u8QbVrdChQovwCb9AirW1Wm99fUJ7ko\n",
			stderr: "Error: couldn't extract an address from line 2",
			isErr:  true,
		},
	}
	for _, tt := range tests {
		stdin := strings.NewReader(tt.stdin)
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		rootCmd := cmd.NewRootCmd(&cli.Env{
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
		}, tt.args)

		t.Run(tt.name, func(t *testing.T) {
			err := rootCmd.Execute()

			if tt.stdout != stdout.String() {
				t.Errorf("mybtc wifi generate returned %s, want %s", stdout, tt.stdout)
			}
			if !strings.HasPrefix(stderr.String(), tt.stderr) {
				t.Errorf("mybtc wifi generate returned %s, want %s", stderr, tt.stderr)
			}
			if (err != nil) != tt.isErr {
				t.Errorf("Error happened during execution: %s", err)
			}
		})
	}
}
