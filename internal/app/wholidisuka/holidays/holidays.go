package holidays

import (
	"time"
)

type holidays struct {
	custom  customHolidays
	regular regularHolidays
	public  publicHolidays
}

type LookupResult int

const (
	Holiday LookupResult = iota
	BusinessDay
	Undefined
	Error
)

func New(
	customFile string,
	regulars []string,
	dataUrl string, cacheDir string, cacheExpire time.Duration, silent bool,
) holidays {
	return holidays{
		custom:  newCustomHolidays(customFile),
		regular: newRegularHolidays(regulars),
		public:  newPublicHolidays(dataUrl, cacheDir, cacheExpire, silent),
	}
}

func (h holidays) LookupCustomHolidays(t time.Time) (LookupResult, error) {
	return h.custom.lookup(t)
}

func (h holidays) LookupRegularHolidays(t time.Time) (LookupResult, error) {
	return h.regular.lookup(t)
}

func (h holidays) LookupPublicHolidays(t time.Time) (LookupResult, error) {
	return h.public.lookup(t)
}
