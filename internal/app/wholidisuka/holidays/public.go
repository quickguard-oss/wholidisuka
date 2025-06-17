package holidays

import (
	"io"
	"net/http"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/quickguard-oss/wholidisuka/internal/app/wholidisuka/cache"
	"github.com/quickguard-oss/wholidisuka/internal/app/wholidisuka/log"
)

const cacheKey = "holidays.yml"

type publicHolidays struct {
	dataUrl     string
	cacheDir    string
	cacheExpire time.Duration
	silent      bool
}

func newPublicHolidays(dataUrl string, cacheDir string, cacheExpire time.Duration, silent bool) publicHolidays {
	return publicHolidays{
		dataUrl:     dataUrl,
		cacheDir:    cacheDir,
		cacheExpire: cacheExpire,
		silent:      silent,
	}
}

func (p publicHolidays) lookup(t time.Time) (LookupResult, error) {
	list, err := p.readData()

	if err != nil {
		return Error, err
	}

	_, ok := list[t.Format("2006-01-02")]

	if !ok {
		return Undefined, nil
	}

	return Holiday, nil
}

func (p publicHolidays) readData() (map[string]any, error) {
	var list map[string]any

	raw, err := p.readCache()

	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(raw, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func (p publicHolidays) readCache() ([]byte, error) {
	c, err := cache.New(p.cacheDir, p.silent)

	if err != nil {
		log.Error(err)

		raw, err := p.getDataUrl()

		if err != nil {
			return nil, err
		}

		return raw, nil
	}

	raw, err := c.Get(cacheKey, p.cacheExpire)

	if err != nil {
		log.Error(err)
	}

	if raw == nil {
		raw, err = p.getDataUrl()

		if err != nil {
			return nil, err
		}

		if err = c.Set(cacheKey, raw); err != nil {
			log.Error(err)
		}
	}

	return raw, nil
}

func (p publicHolidays) getDataUrl() ([]byte, error) {
	r, err := http.Get(p.dataUrl)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
