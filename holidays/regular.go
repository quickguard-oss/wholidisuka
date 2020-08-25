package holidays

import (
	"strings"
	"time"
)

type regularHolidays struct {
	regulars []string
}

func newRegularHolidays(regulars []string) *regularHolidays {
	return &regularHolidays{
		regulars: regulars,
	}
}

func (r *regularHolidays) isRegularHoliday(w time.Weekday) bool {
	symbol := strings.ToLower(w.String()[:3])

	for _, regular := range r.regulars {
		if strings.ToLower(regular) == symbol {
			return true
		}
	}

	return false
}
