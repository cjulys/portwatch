// Package pipeline wires scanner output through filter, diff, classify,
// throttle and notifier in a single reusable processing step.
package pipeline

import (
	"context"
	"io"
	"os"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/classify"
	"portwatch/internal/filter"
	"portwatch/internal/scanner"
	"portwatch/internal/throttle"
)

// Result is the outcome of one pipeline execution.
type Result struct {
	ScannedAt time.Time
	OpenCount int
	AlertCount int
}

// Pipeline orchestrates a single scan-to-alert cycle.
type Pipeline struct {
	scanner   *scanner.Scanner
	filter    *filter.Filter
	classify  *classify.Classifier
	throttle  *throttle.Throttle
	alerter   *alert.Alerter
	fallback  io.Writer
}

// Config holds the dependencies needed to build a Pipeline.
type Config struct {
	Scanner  *scanner.Scanner
	Filter   *filter.Filter
	Classify *classify.Classifier
	Throttle *throttle.Throttle
	Alerter  *alert.Alerter
	Fallback io.Writer
}

// New constructs a Pipeline from the provided Config.
func New(cfg Config) *Pipeline {
	fw := cfg.Fallback
	if fw == nil {
		fw = os.Stderr
	}
	return &Pipeline{
		scanner:  cfg.Scanner,
		filter:   cfg.Filter,
		classify: cfg.Classify,
		throttle: cfg.Throttle,
		alerter:  cfg.Alerter,
		fallback: fw,
	}
}

// Run executes one scan cycle: scan → filter → classify → throttle → alert.
// It returns a Result summarising what happened and any scan-level error.
func (p *Pipeline) Run(ctx context.Context, prev []scanner.Port) ([]scanner.Port, Result, error) {
	res := Result{ScannedAt: time.Now().UTC()}

	current, err := p.scanner.Scan(ctx)
	if err != nil {
		return prev, res, err
	}

	visible := p.filter.Apply(current)
	res.OpenCount = len(visible)

	batch := p.classify.Batch(visible)
	for _, cr := range batch {
		key := cr.Port.Proto + ":"
		if !p.throttle.Allow(key) {
			continue
		}
		p.alerter.Notify([]scanner.Port{cr.Port})
		res.AlertCount++
	}

	return visible, res, nil
}
