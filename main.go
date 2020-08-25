package main

import (
	"fmt"
	"os"

	"github.com/quickguard-oss/wholidisuka/cmd"
)

func main() {
	code, err := cmd.Execute()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(code)
}
