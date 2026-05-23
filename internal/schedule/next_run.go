// Package schedule provides utilities for computing next run times
// and human-readable descriptions of cron schedules.
package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// NextRun holds information about the next scheduled execution.
type NextRun struct {
	// At is the absolute time of the next run.
	At time.Time
	// In is a human-readable duration until the next run.
	In string
	// Schedule is the original cron expression.
	Schedule string
}

// NextN returns the next n run times for the given cron expression,
// computed relative to from.
func NextN(expr string, from time.Time, n int) ([]NextRun, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be greater than zero")
	}

	parser := cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)

	sched, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression %q: %w", expr, err)
	}

	runs := make([]NextRun, 0, n)
	t := from
	for i := 0; i < n; i++ {
		t = sched.Next(t)
		runs = append(runs, NextRun{
			At:       t,
			In:       humanDuration(t.Sub(from)),
			Schedule: expr,
		})
	}
	return runs, nil
}

// humanDuration converts a duration into a short human-readable string.
func humanDuration(d time.Duration) string {
	if d < 0 {
		return "in the past"
	}
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		m := int(d.Minutes())
		s := int(d.Seconds()) % 60
		if s == 0 {
			return fmt.Sprintf("%dm", m)
		}
		return fmt.Sprintf("%dm %ds", m, s)
	case d < 24*time.Hour:
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		if m == 0 {
			return fmt.Sprintf("%dh", h)
		}
		return fmt.Sprintf("%dh %dm", h, m)
	default:
		days := int(d.Hours()) / 24
		h := int(d.Hours()) % 24
		if h == 0 {
			return fmt.Sprintf("%dd", days)
		}
		return fmt.Sprintf("%dd %dh", days, h)
	}
}
