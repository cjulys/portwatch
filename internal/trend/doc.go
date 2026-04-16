// Package trend provides sliding-window trend analysis over port-scan history.
//
// Given a [history.History] store and a time window, [Analyzer.Analyze] returns
// whether the number of open ports is rising, falling, or stable, together with
// the net delta and the number of history samples examined.
package trend
