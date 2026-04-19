// Package retry implements configurable retry logic with exponential backoff
// for use in portwatch when performing operations that may transiently fail,
// such as webhook delivery, probe connections, or file I/O.
//
// Usage:
//
//	r := retry.New(retry.DefaultPolicy())
//	err := r.Do(ctx, func() error {
//		return doSomething()
//	})
package retry
