package holidays

import (
	"time"
)

type holidays struct {
	regularHolidays *regularHolidays
	publicHolidays  *publicHolidays
	customCalendar  *customCalendar
}

type option func(*holidays)

type ListMatch int

const (
	BusinessDay ListMatch = iota
	Holiday
	Undefined
	Error
)

func New(regulars []string, cacheDir string, expire time.Duration, customFile string, opts ...option) *holidays {
	h := &holidays{
		regularHolidays: newRegularHolidays(regulars),
		publicHolidays:  newPublicHolidays(cacheDir, expire),
		customCalendar:  newCustomCalendar(customFile),
	}

	for _, optFunc := range opts {
		optFunc(h)
	}

	return h
}

func HolidaysPublicHolidaysURL(v string) option {
	return func(h *holidays) {
		h.publicHolidays.setPublicHolydaysUrl(v)
	}
}

func HolidaysSilent(v bool) option {
	return func(h *holidays) {
		h.publicHolidays.setSilent(v)
	}
}

func lookup(list *map[string]interface{}, t time.Time) (interface{}, bool) {
	v, ok := (*list)[t.Format("2006-01-02")]

	return v, ok
}

func (h *holidays) IsRegularHoliday(w time.Weekday) bool {
	return h.regularHolidays.isRegularHoliday(w)
}

func (h *holidays) LookupJapaneseHolidays(t time.Time) (ListMatch, error) {
	return h.publicHolidays.lookupJapaneseHolidays(t)
}

func (h *holidays) LookupCustomCalendar(t time.Time) (ListMatch, error) {
	return h.customCalendar.lookupCustomCalendar(t)
}
