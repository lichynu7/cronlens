// Package conflict provides utilities for detecting scheduling conflicts
// between multiple cron expressions.
package conflict

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// Conflict represents a scheduling overlap between two cron entries.
type Conflict struct {
	EntryA    string
	EntryB    string
	NextShared time.Time
}

// String returns a human-readable description of the conflict.
func (c Conflict) String() string {
	return fmt.Sprintf("conflict: %q and %q both run at %s",
		c.EntryA, c.EntryB, c.NextShared.Format(time.RFC3339))
}

// Detect finds all pairs of cron expressions that share at least one
// execution time within the next `window` duration from `from`.
func Detect(expressions []string, from time.Time, window time.Duration) ([]Conflict, error) {
	type schedEntry struct {
		expr  string
		times []time.Time
	}

	entries := make([]schedEntry, 0, len(expressions))
	for _, expr := range expressions {
		sched, err := cron.ParseStandard(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid expression %q: %w", expr, err)
		}
		times := collectTimes(sched, from, from.Add(window))
		entries = append(entries, schedEntry{expr: expr, times: times})
	}

	var conflicts []Conflict
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if t, ok := firstShared(entries[i].times, entries[j].times); ok {
				conflicts = append(conflicts, Conflict{
					EntryA:     entries[i].expr,
					EntryB:     entries[j].expr,
					NextShared: t,
				})
			}
		}
	}
	return conflicts, nil
}

// collectTimes returns all scheduled times in [from, until).
func collectTimes(sched cron.Schedule, from, until time.Time) []time.Time {
	var times []time.Time
	t := sched.Next(from.Add(-time.Second))
	for !t.IsZero() && t.Before(until) {
		times = append(times, t.Truncate(time.Minute))
		t = sched.Next(t)
	}
	return times
}

// firstShared returns the earliest time present in both slices (minute precision).
func firstShared(a, b []time.Time) (time.Time, bool) {
	set := make(map[time.Time]struct{}, len(a))
	for _, t := range a {
		set[t] = struct{}{}
	}
	for _, t := range b {
		if _, ok := set[t]; ok {
			return t, true
		}
	}
	return time.Time{}, false
}
