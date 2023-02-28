package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func buildCache(t *testing.T) cache {
	t.Helper()

	tmpDir := filepath.Join(os.TempDir(), "wholidisuka")

	c, err := New(tmpDir, true)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatal(err)
		}
	})

	return c
}

func Test_Get(t *testing.T) {
	assert := assert.New(t)

	t.Run("alive", func(t *testing.T) {
		c := buildCache(t)

		c.Set("key", []byte("a"))

		got, err := c.Get("key", time.Hour*100)

		assert.Equal([]byte("a"), got)
		assert.Nil(err)
	})

	t.Run("not exist", func(t *testing.T) {
		c := buildCache(t)

		got, err := c.Get("key", time.Hour*100)

		assert.Nil(got)
		assert.Nil(err)
	})

	t.Run("expired", func(t *testing.T) {
		c := buildCache(t)

		c.Set("key", []byte("a"))

		got, err := c.Get("key", time.Second*0)

		assert.Nil(got)
		assert.Nil(err)
	})
}
