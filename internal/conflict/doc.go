// Package conflict implements schedule conflict detection for cronlens.
//
// A conflict occurs when two or more cron expressions share an execution
// time within a given observation window. This is useful for identifying
// resource contention or unintended overlapping jobs in a crontab.
//
// Basic usage:
//
//	import "github.com/user/cronlens/internal/conflict"
//
//	conflicts, err := conflict.Detect(
//		[]string{"0 * * * *", "0 */2 * * *"},
//		time.Now(),
//		24*time.Hour,
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, c := range conflicts {
//		fmt.Println(c)
//	}
//
// The detection window is configurable; a larger window increases
// detection accuracy for low-frequency schedules at the cost of
// additional computation.
package conflict
