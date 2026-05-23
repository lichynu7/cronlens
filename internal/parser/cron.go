package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// CronExpression represents a parsed cron schedule.
type CronExpression struct {
	Raw     string
	Minute  string
	Hour    string
	Day     string
	Month   string
	Weekday string
}

// Parse parses a cron expression string into a CronExpression.
// Supports standard 5-field cron format: minute hour day month weekday
func Parse(expr string) (*CronExpression, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, fmt.Errorf("empty cron expression")
	}

	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression %q: expected 5 fields, got %d", expr, len(fields))
	}

	c := &CronExpression{
		Raw:     expr,
		Minute:  fields[0],
		Hour:    fields[1],
		Day:     fields[2],
		Month:   fields[3],
		Weekday: fields[4],
	}

	if err := validateField(c.Minute, 0, 59, "minute"); err != nil {
		return nil, err
	}
	if err := validateField(c.Hour, 0, 23, "hour"); err != nil {
		return nil, err
	}
	if err := validateField(c.Day, 1, 31, "day"); err != nil {
		return nil, err
	}
	if err := validateField(c.Month, 1, 12, "month"); err != nil {
		return nil, err
	}
	if err := validateField(c.Weekday, 0, 7, "weekday"); err != nil {
		return nil, err
	}

	return c, nil
}

// validateField checks that a cron field is syntactically valid within [min, max].
func validateField(field string, min, max int, name string) error {
	if field == "*" {
		return nil
	}
	// Handle step values like */5 or 1-5/2
	parts := strings.Split(field, "/")
	if len(parts) > 2 {
		return fmt.Errorf("invalid %s field %q: multiple slashes", name, field)
	}
	if len(parts) == 2 {
		step, err := strconv.Atoi(parts[1])
		if err != nil || step <= 0 {
			return fmt.Errorf("invalid %s field %q: bad step value", name, field)
		}
		if parts[0] == "*" {
			return nil
		}
		field = parts[0]
	}
	// Handle ranges like 1-5
	rangeParts := strings.Split(field, "-")
	if len(rangeParts) == 2 {
		lo, err1 := strconv.Atoi(rangeParts[0])
		hi, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil || lo > hi || lo < min || hi > max {
			return fmt.Errorf("invalid %s field %q: bad range", name, field)
		}
		return nil
	}
	// Handle lists like 1,2,3
	for _, part := range strings.Split(field, ",") {
		v, err := strconv.Atoi(part)
		if err != nil || v < min || v > max {
			return fmt.Errorf("invalid %s field %q: value %q out of range [%d,%d]", name, field, part, min, max)
		}
	}
	return nil
}
