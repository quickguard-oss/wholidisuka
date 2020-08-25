package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCacheDir(t *testing.T) {
	defer os.Setenv("HOME", os.Getenv("HOME"))

	os.Setenv("HOME", "/myhome")

	assert.Equal(t, "/myhome/.wholidisuka/cache", defaultCacheDir())
}

func TestTargetDate(t *testing.T) {
	t.Run("No argv", func(t *testing.T) {
		got, _ := targetDate([]string{})

		assert.Less(t, time.Since(got).Seconds(), 3.0)
	})

	t.Run("1 argv", func(t *testing.T) {
		expected := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

		actual, _ := targetDate([]string{"1970-01-01"})

		assert.Equal(t, expected, actual)
	})
}

func TestCmdExecute(t *testing.T) {
	cases := []struct {
		argv     []string
		desc     string
		expected int
	}{
		{
			argv:     []string{"1970-01-01"},
			desc:     "1970-01-01 is a public holiday.",
			expected: 0,
		},
		{
			argv:     []string{"2019-09-19"},
			desc:     "2019-09-19 is not a public holiday.",
			expected: 1,
		},
		{
			argv:     []string{"--regular", "thu,fri", "2019-09-19"},
			desc:     "2019-09-19 (Thu) is a regular holiday (Thu, Fri).",
			expected: 0,
		},
		{
			argv:     []string{"--regular", "sat,sun", "2019-09-19"},
			desc:     "2019-09-19 (Thu) is not a regular holiday (Sat, Sun).",
			expected: 1,
		},
		{
			argv:     []string{"--override", "./testdata/custom.yml", "1970-01-01"},
			desc:     "1970-01-01 is not a holiday in my calendar.",
			expected: 1,
		},
		{
			argv:     []string{"--override", "./testdata/custom.yml", "2019-09-19"},
			desc:     "2019-09-19 is a holiday in my calendar.",
			expected: 0,
		},
		{
			argv:     []string{"--regular", "fri", "--override", "./testdata/custom.yml", "2019-09-20"},
			desc:     "Fri is a regular holiday, but 2019-09-20 is a business day in my calendar.",
			expected: 1,
		},
	}

	opts := []string{"--cache-dir", "./testdata", "--cache-expire", "438000h"}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			var actual int

			cmd := newCmd()

			cmd.SetArgs(
				append(opts, c.argv...),
			)

			err := cmd.Execute()

			switch err {
			case nil:
				actual = 0
			case unmatchError:
				actual = 1
			default:
				actual = 2
			}

			assert.Equal(t, c.expected, actual)
		})
	}
}
