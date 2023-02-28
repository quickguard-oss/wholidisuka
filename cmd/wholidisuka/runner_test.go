package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/quickguard-oss/wholidisuka/internal/app/wholidisuka/test/helper"
)

func setupFlag(t *testing.T) {
	t.Helper()

	old := flagOutput

	flagOutput = io.Discard

	t.Cleanup(func() {
		flagOutput = old
	})
}

func Test_run(t *testing.T) {
	assert := assert.New(t)

	type expected struct {
		returnCode returnCode
		isError    bool
	}

	type testCase struct {
		name     string
		args     []string
		expected expected
	}

	testCases := []testCase{
		{
			name: "national holiday = yes",
			args: []string{"2023-01-01"},
			expected: expected{
				returnCode: AsHoliday,
				isError:    false,
			},
		},
		{
			name: "national holiday = no",
			args: []string{"2023-09-19"},
			expected: expected{
				returnCode: AsBusinessDay,
				isError:    false,
			},
		},
		{
			name: "national holiday = yes && override = yes (business day)",
			args: []string{"-o", "./testdata/custom.yml", "2023-01-01"},
			expected: expected{
				returnCode: AsBusinessDay,
				isError:    false,
			},
		},
		{
			name: "national holiday = no  && override = yes (holiday)",
			args: []string{"-o", "./testdata/custom.yml", "2023-09-19"},
			expected: expected{
				returnCode: AsHoliday,
				isError:    false,
			},
		},
		{
			name: "national holiday = yes && override = no",
			args: []string{"-o", "./testdata/custom.yml", "2023-02-11"},
			expected: expected{
				returnCode: AsHoliday,
				isError:    false,
			},
		},
		{
			name: "national holiday = no  && override = no",
			args: []string{"-o", "./testdata/custom.yml", "2023-09-20"},
			expected: expected{
				returnCode: AsBusinessDay,
				isError:    false,
			},
		},
		{
			name: "national holiday = yes && regular holiday = yes",
			args: []string{"-r", "sun", "2023-01-01"},
			expected: expected{
				returnCode: AsHoliday,
				isError:    false,
			},
		},
		{
			name: "national holiday = yes && regular holiday = no",
			args: []string{"-r", "fri", "2023-01-01"},
			expected: expected{
				returnCode: AsHoliday,
				isError:    false,
			},
		},
		{
			name: "national holiday = no  && regular holiday = yes",
			args: []string{"-r", "tue", "2023-09-19"},
			expected: expected{
				returnCode: AsHoliday,
				isError:    false,
			},
		},
		{
			name: "national holiday = no  && regular holiday = no",
			args: []string{"-r", "fri", "2023-09-19"},
			expected: expected{
				returnCode: AsBusinessDay,
				isError:    false,
			},
		},
		{
			name: "national holiday = yes && regular holiday = yes && override = yes (business day)",
			args: []string{"-o", "./testdata/custom.yml", "-r", "sun", "2023-01-01"},
			expected: expected{
				returnCode: AsBusinessDay,
				isError:    false,
			},
		},
		{
			name: "national holiday = yes && regular holiday = no  && override = yes (business day)",
			args: []string{"-o", "./testdata/custom.yml", "-r", "fri", "2023-01-01"},
			expected: expected{
				returnCode: AsBusinessDay,
				isError:    false,
			},
		},
		{
			name: "national holiday = no  && regular holiday = yes && override = yes (holiday)",
			args: []string{"-o", "./testdata/custom.yml", "-r", "tue", "2023-09-19"},
			expected: expected{
				returnCode: AsHoliday,
				isError:    false,
			},
		},
		{
			name: "national holiday = no  && regular holiday = no  && override = yes (holiday)",
			args: []string{"-o", "./testdata/custom.yml", "-r", "fri", "2023-09-19"},
			expected: expected{
				returnCode: AsHoliday,
				isError:    false,
			},
		},
		{
			name: "national holiday = no  && regular holiday = yes && override = no)",
			args: []string{"-o", "./testdata/custom.yml", "-r", "wed", "2023-09-20"},
			expected: expected{
				returnCode: AsHoliday,
				isError:    false,
			},
		},
		{
			name: "national holiday = no  && regular holiday = no  && override = no",
			args: []string{"-o", "./testdata/custom.yml", "-r", "fri", "2023-09-20"},
			expected: expected{
				returnCode: AsBusinessDay,
				isError:    false,
			},
		},
		{
			name: "invalid flag",
			args: []string{"-xxx", "2023-01-01"},
			expected: expected{
				returnCode: AsError,
				isError:    true,
			},
		},
		{
			name: "invalid date",
			args: []string{"0000-00-00"},
			expected: expected{
				returnCode: AsError,
				isError:    true,
			},
		},
	}

	opts := []string{"-d", "./testdata/", "-e", "876000h"}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setupFlag(t)

			helper.SetCommandArgs(t, append(opts, tc.args...))

			got, err := run()

			assert.Equal(tc.expected.returnCode, got)

			if tc.expected.isError {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
			}
		})
	}
}
