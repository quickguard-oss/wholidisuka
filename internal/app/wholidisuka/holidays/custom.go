package holidays

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type customHolidays struct {
	customFile string
}

func newCustomHolidays(customFile string) customHolidays {
	return customHolidays{
		customFile: customFile,
	}
}

func (c customHolidays) lookup(t time.Time) (LookupResult, error) {
	if c.customFile == "" {
		return Undefined, nil
	}

	list, err := c.readData()

	if err != nil {
		return Error, err
	}

	v, ok := list[t.Format("2006-01-02")]

	if !ok {
		return Undefined, nil
	}

	if v == nil {
		return BusinessDay, nil
	}

	return Holiday, nil
}

func (c customHolidays) readData() (map[string]any, error) {
	var list map[string]any

	raw, err := os.ReadFile(c.customFile)

	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(raw, &list); err != nil {
		return nil, err
	}

	return list, nil
}
