package main

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func TestDateFilterFlagsWorkWithAliases(t *testing.T) {
	cases := map[string]struct {
		args []string
		want int
	}{
		"none":              {nil, 0},
		"from alias":        {[]string{"--from", "2026-09-01T00:00:00Z"}, 1},
		"from alias with =": {[]string{"--from=2026-09-01T00:00:00Z"}, 1},
		"to alias":          {[]string{"--to", "2026-09-30T00:00:00Z"}, 1},
		"both aliases":      {[]string{"--from", "2026-09-01T00:00:00Z", "--to", "2026-09-30T00:00:00Z"}, 2},
		"full names":        {[]string{"--output_filter_start_date", "2026-09-01T00:00:00Z", "--output_filter_end_date", "2026-09-30T00:00:00Z"}, 2},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := -1
			testApp := &cli.App{
				Flags: dateFilterFlags(),
				Action: func(c *cli.Context) error {
					got = len(dateFilterTransformers(c))
					return nil
				},
			}
			if err := testApp.Run(append([]string{"mappyboi"}, tc.args...)); err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("date filters applied = %d, want %d", got, tc.want)
			}
		})
	}
}
