package holidays

import (
	"strings"
	"time"
)

type regularHolidays struct {
	regulars []string
}

func newRegularHolidays(regulars []string) regularHolidays {
	return regularHolidays{
		regulars: regulars,
	}
}

func (r regularHolidays) lookup(t time.Time) (LookupResult, error) {
	symbol := strings.ToLower(t.Weekday().String()[:3])

	for _, regular := range r.regulars {
		if strings.ToLower(regular) == symbol {
			return Holiday, nil
		}
	}

	return Undefined, nil
}
