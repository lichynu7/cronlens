// Package parser provides functionality for parsing and validating
// standard 5-field cron expressions used by cronlens.
//
// A cron expression consists of five space-separated fields:
//
//	┌───────────── minute        (0–59)
//	│ ┌─────────── hour          (0–23)
//	│ │ ┌───────── day of month  (1–31)
//	│ │ │ ┌─────── month         (1–12)
//	│ │ │ │ ┌───── day of week   (0–7, where 0 and 7 are Sunday)
//	│ │ │ │ │
//	* * * * *
//
// Supported field syntax:
//   - *         — every unit (wildcard)
//   - n         — exact value
//   - n,m,...   — list of values
//   - n-m       — inclusive range
//   - */step    — every step units
//   - n-m/step  — range with step
//
// Example usage:
//
//	expr, err := parser.Parse("*/5 9-17 * * 1-5")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(expr.Hour) // "9-17"
package parser
