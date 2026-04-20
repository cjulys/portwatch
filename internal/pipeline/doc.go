// Package pipeline provides a single-call abstraction that connects the
// scanner, filter, classifier, throttle and alerter packages into one
// coherent scan-to-alert processing step.
//
// Typical usage:
//
//	p := pipeline.New(pipeline.Config{
//		Scanner:  s,
//		Filter:   f,
//		Classify: c,
//		Throttle: t,
//		Alerter:  a,
//	})
//
//	var prev []scanner.Port
//	for {
//		prev, result, err = p.Run(ctx, prev)
//	}
package pipeline
