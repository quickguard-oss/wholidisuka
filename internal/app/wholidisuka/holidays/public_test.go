package holidays

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func buildHolidays(t *testing.T, cacheDir string, cacheExpire time.Duration, response string) publicHolidays {
	t.Helper()

	s := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, response)
			},
		),
	)

	t.Cleanup(func() {
		s.Close()
	})

	return newPublicHolidays(s.URL, cacheDir, cacheExpire, true)
}

func Test_publicHolidays_readCache(t *testing.T) {
	assert := assert.New(t)

	cacheDir := filepath.Join(os.TempDir(), "wholidisuka")

	t.Cleanup(func() {
		if err := os.RemoveAll(cacheDir); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("cache not exist", func(t *testing.T) {
		h := buildHolidays(t, cacheDir, time.Hour, "---\n2023-09-19: 'クイック (9.19) ガードの日'")

		got, err := h.readCache()

		assert.Equal("---\n2023-09-19: 'クイック (9.19) ガードの日'", string(got))
		assert.Nil(err)
	})

	t.Run("cache alive", func(t *testing.T) {
		h := buildHolidays(t, cacheDir, time.Hour, "---\n2023-01-01: 元日")

		got, err := h.readCache()

		assert.Equal("---\n2023-09-19: 'クイック (9.19) ガードの日'", string(got))
		assert.Nil(err)
	})

	t.Run("cache expired", func(t *testing.T) {
		h := buildHolidays(t, cacheDir, time.Second*0, "---\n2023-02-11: 建国記念の日")

		got, err := h.readCache()

		assert.Equal("---\n2023-02-11: 建国記念の日", string(got))
		assert.Nil(err)
	})
}
