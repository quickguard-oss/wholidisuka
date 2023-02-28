package helper

import (
	"flag"
	"os"
	"testing"
)

func SetCommandArgs(t *testing.T, args []string) {
	t.Helper()

	t.Cleanup(func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	})

	old := os.Args

	os.Args = append([]string{os.Args[0]}, args...)

	t.Cleanup(func() {
		os.Args = old
	})
}
