package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/quickguard-oss/wholidisuka/internal/app/wholidisuka/log"
)

type arrayFlag []string

type flags struct {
	regulars    arrayFlag
	dataUrl     string
	cacheDir    string
	cacheExpire time.Duration
	silent      bool
	customFile  string
	showVersion bool
	showHelp    bool
}

const dataUrlDefault = "https://raw.githubusercontent.com/holiday-jp/holiday_jp/master/holidays.yml"

var flagOutput io.Writer = os.Stderr

func parseFlags() (flags, error) {
	setupParser()

	f := newFlags()

	f.defineFlags()

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return f, err
	}

	if f.showVersion {
		fmt.Println(version)

		os.Exit(0)
	}

	if f.showHelp {
		flag.Usage()

		os.Exit(0)
	}

	return f, nil
}

func newFlags() flags {
	return flags{}
}

func setupParser() {
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)

	flag.CommandLine.SetOutput(flagOutput)

	flag.Usage = func() {
		header := fmt.Sprintf("wholidisuka (v%s)\n", version)

		usage := `
Wholidisuka is a command line tool that checks whether today is a holiday in Japan.

Usage:
  wholidisuka [-r day_of_week -r ...] [-o file] [-d directory] [-e duration] [date]

  Use the <date> argument in the format 'YYYY-MM-DD' to check a specific date, otherwise today's
  date will be used by default.

Exit codes:
  0: holiday
  1: business day
  2: error occurred

Options:
`

		fmt.Fprint(
			flag.CommandLine.Output(),
			header+usage,
		)

		flag.PrintDefaults()
	}
}

func (f *flags) defineFlags() {
	flag.Var(
		&f.regulars,
		"r",
		""+
			"Set regular holidays.\n"+
			"You can set multiple days by passing this option multiple times.\n"+
			"The valid values for <`day_of_week`> are 'sun', 'mon', 'tue', 'wed', 'thu', 'fri', and 'sat'.\n"+
			"\n"+
			"Example: -r 'sat' -r 'sun'",
	)

	flag.StringVar(
		&f.dataUrl,
		"data-url",
		dataUrlDefault,
		""+
			"Specify a `URL` for the dataset.\n"+
			"\n",
	)

	flag.StringVar(
		&f.cacheDir,
		"d",
		defaultCacheDir(),
		""+
			"Set the `directory` where the cache is stored.\n"+
			"\n",
	)

	flag.DurationVar(
		&f.cacheExpire,
		"e",
		time.Hour*24*180,
		""+
			"Set the expiration time for the cache. \n"+
			"The `duration` must be a string in the format of golang 'time.ParseDuration()'.\n"+
			"\n"+
			"Example: -e '72h'\n"+
			"\n",
	)

	flag.BoolVar(
		&f.silent,
		"silent",
		false,
		"Silent mode.",
	)

	flag.StringVar(
		&f.customFile,
		"o",
		"",
		""+
			"Use your own calendar.\n"+
			"The `file` must be in YAML format and compatible with holiday-jp.",
	)

	flag.BoolVar(
		&f.showVersion,
		"version",
		false,
		"Print version information and exit.",
	)

	flag.BoolVar(
		&f.showHelp,
		"h",
		false,
		"Print help message and exit.",
	)

}

func defaultCacheDir() string {
	dir := ".wholidisuka/cache"

	home, err := os.UserHomeDir()

	if err != nil {
		log.Error(err)

		return dir
	}

	return filepath.Join(home, dir)
}

func (a *arrayFlag) Set(value string) error {
	*a = append(*a, value)

	return nil
}

func (a *arrayFlag) String() string {
	return strings.Join(*a, ", ")
}

func (f flags) targetDate() (time.Time, error) {
	if len(flag.Args()) == 0 {
		return time.Now(), nil
	} else {
		return time.Parse("2006-01-02", flag.Arg(0))
	}
}
