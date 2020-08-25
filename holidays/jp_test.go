package holidays

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const serverResponse = "---\n1970-01-01: 元日"

type Suite struct {
	suite.Suite

	cacheDir string
	server   *httptest.Server
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, serverResponse)
	}))
}

func (s *Suite) TearDownSuite() {
	s.server.Close()
}

func (s *Suite) SetupTest() {
	dir, err := ioutil.TempDir("", "wholidisuka-")

	if err != nil {
		s.Fail("Failed to create cache dir.")
	}

	s.cacheDir = dir
}

func (s *Suite) TearDownTest() {
	if err := os.RemoveAll(s.cacheDir); err != nil {
		s.Fail("Failed to remove cache dir.")
	}
}

func (s *Suite) TestReadHolidayJpYmlCache_NoCache() {
	p := s.newPublicHolidays(time.Hour * 100)

	s.Run("Func return", func() {
		got, _ := p.readHolidayJpYmlCache()

		s.Equal(serverResponse, string(got))
	})

	s.Run("Set cache", func() {
		cacheData, err := s.readCache()

		if err != nil {
			s.Fail("Failed to read cache file.")
		}

		s.Equal(serverResponse, string(cacheData))
	})
}

func (s *Suite) TestReadHolidayJpYmlCache_CacheAlive() {
	p := s.newPublicHolidays(time.Hour * 100)

	cacheContent := []byte("---\n2019-09-19: 'クイック (9.19) ガードの日'")

	if err := s.writeCache(cacheContent); err != nil {
		s.Fail("Failed to write cache file.")
	}

	got, _ := p.readHolidayJpYmlCache()

	s.Equal(cacheContent, got)
}

func (s *Suite) TestReadHolidayJpYmlCache_CacheExpired() {
	p := s.newPublicHolidays(time.Second * 0)

	s.Run("Func return", func() {
		cacheContent := []byte("---\n2019-09-19: 'クイック (9.19) ガードの日'")

		if err := s.writeCache(cacheContent); err != nil {
			s.Fail("Failed to write cache file.")
		}

		got, _ := p.readHolidayJpYmlCache()

		s.Equal(serverResponse, string(got))
	})

	s.Run("Set cache", func() {
		cacheData, err := s.readCache()

		if err != nil {
			s.Fail("Failed to read cache file.")
		}

		s.Equal(serverResponse, string(cacheData))
	})
}

func (s *Suite) newPublicHolidays(expire time.Duration) *publicHolidays {
	p := newPublicHolidays(s.cacheDir, expire)

	p.setPublicHolydaysUrl(s.server.URL)
	p.setSilent(true)

	return p
}

func (s *Suite) readCache() ([]byte, error) {
	return ioutil.ReadFile(s.cacheFilePath())
}

func (s *Suite) writeCache(data []byte) error {
	return ioutil.WriteFile(s.cacheFilePath(), data, os.ModePerm)
}

func (s *Suite) cacheFilePath() string {
	return filepath.Join(s.cacheDir, holidayJpCacheKey)
}
