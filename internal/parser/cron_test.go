package parser

import (
	"testing"
)

func TestParse_Valid(t *testing.T) {
	cases := []struct {
		expr string
	}{
		{"* * * * *"},
		{"0 12 * * *"},
		{"*/5 * * * *"},
		{"0 9-17 * * 1-5"},
		{"30 4 1,15 * *"},
		{"0 0 1 1 *"},
		{"0 */2 * * *"},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			c, err := Parse(tc.expr)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.expr, err)
			}
			if c.Raw != tc.expr {
				t.Errorf("expected Raw=%q, got %q", tc.expr, c.Raw)
			}
		})
	}
}

func TestParse_Invalid(t *testing.T) {
	cases := []struct {
		expr    string
		wantErr string
	}{
		{"", "empty cron expression"},
		{"* * *", "expected 5 fields"},
		{"60 * * * *", "minute"},
		{"* 25 * * *", "hour"},
		{"* * 32 * *", "day"},
		{"* * * 13 *", "month"},
		{"* * * * 8", "weekday"},
		{"1/2/3 * * * *", "multiple slashes"},
		{"1-60 * * * *", "minute"},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			_, err := Parse(tc.expr)
			if err == nil {
				t.Fatalf("Parse(%q) expected error containing %q, got nil", tc.expr, tc.wantErr)
			}
		})
	}
}

func TestParse_Fields(t *testing.T) {
	c, err := Parse("30 6 1 3 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Minute != "30" || c.Hour != "6" || c.Day != "1" || c.Month != "3" || c.Weekday != "5" {
		t.Errorf("unexpected fields: %+v", c)
	}
}
