package conflict_test

import (
	"testing"
	"time"

	"github.com/user/cronlens/internal/conflict"
)

func baseTime() time.Time {
	// Monday 2024-01-01 00:00:00 UTC
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestDetect_NoConflict(t *testing.T) {
	expressions := []string{
		"0 * * * *",  // every hour at :00
		"30 * * * *", // every hour at :30
	}
	conflicts, err := conflict.Detect(expressions, baseTime(), 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d: %v", len(conflicts), conflicts)
	}
}

func TestDetect_WithConflict(t *testing.T) {
	expressions := []string{
		"0 * * * *",   // every hour at :00
		"0 */2 * * *", // every 2 hours at :00
	}
	conflicts, err := conflict.Detect(expressions, baseTime(), 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) == 0 {
		t.Error("expected at least one conflict, got none")
	}
}

func TestDetect_InvalidExpression(t *testing.T) {
	expressions := []string{
		"not-a-cron",
		"0 * * * *",
	}
	_, err := conflict.Detect(expressions, baseTime(), time.Hour)
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestDetect_IdenticalExpressions(t *testing.T) {
	expressions := []string{
		"*/5 * * * *",
		"*/5 * * * *",
	}
	conflicts, err := conflict.Detect(expressions, baseTime(), time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) == 0 {
		t.Error("expected conflicts for identical expressions")
	}
}

func TestConflict_String(t *testing.T) {
	c := conflict.Conflict{
		EntryA:     "0 * * * *",
		EntryB:     "0 */2 * * *",
		NextShared: baseTime(),
	}
	s := c.String()
	if s == "" {
		t.Error("expected non-empty string from Conflict.String()")
	}
}
