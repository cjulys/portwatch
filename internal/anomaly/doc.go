// Package anomaly provides statistical anomaly detection for open-port counts.
//
// It maintains a rolling window of historical scan counts and computes a
// z-score for each new observation. Events are emitted at three severity
// levels — mild (1–2 σ), moderate (2–3 σ), and severe (> 3 σ) — allowing
// the caller to decide how to react.
//
// Typical usage:
//
//	det := anomaly.New(60, nil) // keep last 60 samples
//	if ev := det.Evaluate(len(ports)); ev != nil {
//	    log.Println(ev)
//	}
package anomaly
