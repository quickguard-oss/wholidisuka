package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/quickguard-oss/wholidisuka/holidays"
	"github.com/spf13/cobra"
)

var (
	regulars   []string
	customFile string
	cacheDir   string
	expire     time.Duration

	unmatchError = errors.New("Today is a business day")

	version = "v0.0.0"
)

func Execute() (int, error) {
	err := newCmd().Execute()

	switch err {
	case nil:
		return 0, nil
	case unmatchError:
		return 1, nil
	default:
		return 2, err
	}
}

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wholidisuka [YYYY-MM-DD (default today)]",
		Short: "Wholidisuka is a command line tool to determine if today is a Japanese holiday.",
		Long: `Wholidisuka is a command line tool to determine if today is a Japanese holiday.
The exit status is 0 if the date is a holiday and 1 otherwise, but it is 2 if an error occurred.

Complete documentation is available at https://github.com/quickguard-oss/wholidisuka.`,
		Args:          cobra.MaximumNArgs(1),
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          run,
	}

	cmd.Flags().StringSliceVarP(&regulars, "regular", "r", []string{}, `regular holidays (in a 3-letter format and comma separated with no spaces; e.g. "sat,sun")`)
	cmd.Flags().StringVarP(&cacheDir, "cache-dir", "d", defaultCacheDir(), "cache directory path")
	cmd.Flags().DurationVarP(&expire, "cache-expire", "e", time.Hour*24*180, "cache expire time")
	cmd.Flags().StringVarP(&customFile, "override", "o", "", "your calendar path")

	return cmd
}

func defaultCacheDir() string {
	dir := ".wholidisuka/cache"

	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		return dir
	}

	return filepath.Join(home, dir)
}

func run(cmd *cobra.Command, args []string) error {
	date, err := targetDate(args)

	if err != nil {
		return err
	}

	h := holidays.New(regulars, cacheDir, expire, customFile)

	// 1. Check custom calendar
	match, err := h.LookupCustomCalendar(date)

	if err != nil {
		return err
	}

	switch match {
	case holidays.BusinessDay:
		return unmatchError
	case holidays.Holiday:
		return nil
	}

	// 2. Check regular holidays
	if h.IsRegularHoliday(date.Weekday()) {
		return nil
	}

	// 3. Check holiday_jp
	match, err = h.LookupJapaneseHolidays(date)

	if err != nil {
		return err
	}

	if match == holidays.Holiday {
		return nil
	}

	return unmatchError
}

func targetDate(args []string) (time.Time, error) {
	if len(args) == 0 {
		return time.Now(), nil
	} else {
		return time.Parse("2006-01-02", args[0])
	}
}
