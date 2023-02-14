package main

import (
	"os"

	"github.com/quickguard-oss/wholidisuka/internal/app/wholidisuka/log"
)

var version = "0.0.0"

func main() {
	code, err := run()

	if err != nil {
		log.Error(err)
	}

	os.Exit(
		int(code),
	)
}
