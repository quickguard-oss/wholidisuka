package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/quickguard-oss/wholidisuka/internal/app/wholidisuka/test/helper"
)

const timeDiffThreshold = time.Second * 3

func buildFlags(t *testing.T, args []string) flags {
	t.Helper()

	helper.SetCommandArgs(t, args)

	f, err := parseFlags()

	if err != nil {
		t.Fatal(err)
	}

	return f
}

func Test_parseFlags(t *testing.T) {
	assert := assert.New(t)

	t.Run("regulars", func(t *testing.T) {
		f := buildFlags(t, []string{"-r", "sun", "-r", "mon"})

		assert.ElementsMatch([]string{"sun", "mon"}, f.regulars)
	})
}

func Test_defaultCacheDir(t *testing.T) {
	old := os.Getenv("HOME")

	os.Setenv("HOME", "/myhome")

	t.Cleanup(func() {
		os.Setenv("HOME", old)
	})

	assert.Equal(t, "/myhome/.wholidisuka/cache", defaultCacheDir())
}

func Test_flags_targetDate(t *testing.T) {
	assert := assert.New(t)

	t.Run("0 argv", func(t *testing.T) {
		f := buildFlags(t, []string{})

		got, err := f.targetDate()

		assert.WithinDuration(time.Now(), got, timeDiffThreshold)
		assert.Nil(err)
	})

	t.Run("1 argv", func(t *testing.T) {
		f := buildFlags(t, []string{"2022-09-19"})

		got, err := f.targetDate()

		assert.Equal(time.Date(2022, time.September, 19, 0, 0, 0, 0, time.UTC), got)
		assert.Nil(err)
	})

	t.Run("2 argv", func(t *testing.T) {
		f := buildFlags(t, []string{"2022-09-19", "1970-01-01"})

		got, err := f.targetDate()

		assert.Equal(time.Date(2022, time.September, 19, 0, 0, 0, 0, time.UTC), got)
		assert.Nil(err)
	})
}
