package log

import (
	"fmt"
	"os"
)

func Error(err error) {
	fmt.Fprintln(os.Stderr, err)
}
