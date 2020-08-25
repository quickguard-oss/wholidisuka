package holidays

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/quickguard-oss/wholidisuka/cache"

	"gopkg.in/yaml.v2"
)

const (
	holidayJpUrl      = "https://raw.githubusercontent.com/holiday-jp/holiday_jp/master/holidays.yml"
	holidayJpCacheKey = "holidays.yml"
)

type publicHolidays struct {
	cacheDir          string
	expire            time.Duration
	publicHolydaysUrl string
	silent            bool
}

func newPublicHolidays(cacheDir string, expire time.Duration) *publicHolidays {
	return &publicHolidays{
		cacheDir:          cacheDir,
		expire:            expire,
		publicHolydaysUrl: holidayJpUrl,
		silent:            false,
	}
}

func (p *publicHolidays) setPublicHolydaysUrl(v string) {
	p.publicHolydaysUrl = v
}

func (p *publicHolidays) setSilent(v bool) {
	p.silent = v
}

func (p *publicHolidays) lookupJapaneseHolidays(t time.Time) (ListMatch, error) {
	list, err := p.readHolidayJpYml()

	if err != nil {
		return Error, err
	}

	_, ok := lookup(list, t)

	if ok {
		return Holiday, nil
	}

	return Undefined, nil
}

func (p *publicHolidays) readHolidayJpYml() (*map[string]interface{}, error) {
	var list map[string]interface{}

	raw, err := p.readHolidayJpYmlCache()

	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(raw, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

func (p *publicHolidays) readHolidayJpYmlCache() ([]byte, error) {
	var raw []byte

	if c, err := cache.New(p.cacheDir, p.silent); err != nil {
		fmt.Fprintln(os.Stderr, err)

		raw, err = p.httpGetHolidayJpYml()

		if err != nil {
			return nil, err
		}
	} else {
		raw, err = c.Get(holidayJpCacheKey, p.expire)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if raw == nil {
			raw, err = p.httpGetHolidayJpYml()

			if err != nil {
				return nil, err
			}

			if err = c.Set(holidayJpCacheKey, raw); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}

	return raw, nil
}

func (p *publicHolidays) httpGetHolidayJpYml() ([]byte, error) {
	resp, err := http.Get(p.publicHolydaysUrl)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
