package holidays

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type customCalendar struct {
	customFile string
}

func newCustomCalendar(customFile string) *customCalendar {
	return &customCalendar{
		customFile: customFile,
	}
}

func (c *customCalendar) lookupCustomCalendar(t time.Time) (ListMatch, error) {
	if c.customFile == "" {
		return Undefined, nil
	}

	list, err := c.readCustomYml()

	if err != nil {
		return Error, err
	}

	v, ok := lookup(list, t)

	if !ok {
		return Undefined, nil
	}

	if v == nil {
		return BusinessDay, nil
	}

	return Holiday, nil
}

func (c *customCalendar) readCustomYml() (*map[string]interface{}, error) {
	var list map[string]interface{}

	raw, err := ioutil.ReadFile(c.customFile)

	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(raw, &list); err != nil {
		return nil, err
	}

	return &list, nil
}
