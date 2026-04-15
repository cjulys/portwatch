// Package baseline provides functionality for capturing and comparing a
// known-good snapshot of open ports (the "baseline") against the current
// scan results.
//
// Typical usage:
//
//	b := baseline.New("/var/lib/portwatch/baseline.json")
//	_ = b.Load()
//
//	// On first run, record the current ports as the baseline.
//	if b.IsEmpty() {
//		_ = b.Set(currentPorts)
//	}
//
//	// On subsequent runs, check for deviations.
//	violations := baseline.Check(b.Get(), currentPorts)
//	for _, v := range violations {
//		fmt.Printf("violation: port %d (%s) — %s\n", v.Port.Number, v.Port.Protocol, v.Reason)
//	}
package baseline
