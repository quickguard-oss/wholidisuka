package main

import (
	"github.com/quickguard-oss/wholidisuka/internal/app/wholidisuka/holidays"
)

type returnCode int

const (
	AsHoliday returnCode = iota
	AsBusinessDay
	AsError
)

func run() (returnCode, error) {
	flags, err := parseFlags()

	if err != nil {
		return AsError, err
	}

	date, err := flags.targetDate()

	if err != nil {
		return AsError, err
	}

	h := holidays.New(
		flags.customFile,
		flags.regulars,
		flags.dataUrl, flags.cacheDir, flags.cacheExpire, flags.silent,
	)

	// 1. Check custom holidays
	result, err := h.LookupCustomHolidays(date)

	if isResolved(result) {
		return returnCodeOf(result), err
	}

	// 2. Check regular holidays
	result, err = h.LookupRegularHolidays(date)

	if isResolved(result) {
		return returnCodeOf(result), err
	}

	// 3. Check public holidays
	result, err = h.LookupPublicHolidays(date)

	if isResolved(result) {
		return returnCodeOf(result), err
	}

	return AsBusinessDay, nil
}

func isResolved(result holidays.LookupResult) bool {
	return result != holidays.Undefined
}

func returnCodeOf(result holidays.LookupResult) returnCode {
	switch result {
	case holidays.Holiday:
		return AsHoliday
	case holidays.BusinessDay:
		return AsBusinessDay
	default:
		return AsError
	}
}
