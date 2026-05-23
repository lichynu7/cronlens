package schedule_test

import (
	"testing"
	"time"

	"github.com/cronlens/cronlens/internal/schedule"
)

var fixedNow = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestNextN_Basic(t *testing.T) {
	runs, err := schedule.NextN("*/5 * * * *", fixedNow, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 3 {
		t.Fatalf("expected 3 runs, got %d", len(runs))
	}
	// First run should be at 12:05
	want := fixedNow.Add(5 * time.Minute)
	if !runs[0].At.Equal(want) {
		t.Errorf("first run: got %v, want %v", runs[0].At, want)
	}
}

func TestNextN_InvalidExpr(t *testing.T) {
	_, err := schedule.NextN("invalid expr", fixedNow, 1)
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestNextN_ZeroN(t *testing.T) {
	_, err := schedule.NextN("* * * * *", fixedNow, 0)
	if err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
}

func TestNextN_ScheduleField(t *testing.T) {
	expr := "0 9 * * 1"
	runs, err := schedule.NextN(expr, fixedNow, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range runs {
		if r.Schedule != expr {
			t.Errorf("Schedule field: got %q, want %q", r.Schedule, expr)
		}
		if r.In == "" {
			t.Error("In field should not be empty")
		}
	}
}

func TestNextN_InField_Minutes(t *testing.T) {
	runs, err := schedule.NextN("*/10 * * * *", fixedNow, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if runs[0].In == "" {
		t.Error("expected non-empty In field")
	}
	// Next run at 12:10, so In should represent ~10 minutes
	if runs[0].In != "10m" {
		t.Errorf("expected '10m', got %q", runs[0].In)
	}
}
