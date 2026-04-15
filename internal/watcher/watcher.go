package watcher

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Watcher orchestrates periodic port scanning and change detection.
type Watcher struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	store   *state.Store
	alerter *alert.Alerter
}

// New creates a Watcher wired up with the provided dependencies.
func New(cfg *config.Config, sc *scanner.Scanner, st *state.Store, al *alert.Alerter) *Watcher {
	return &Watcher{
		cfg:     cfg,
		scanner: sc,
		store:   st,
		alerter: al,
	}
}

// Run starts the watch loop and blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	interval := time.Duration(w.cfg.IntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if err := w.tick(); err != nil {
		log.Printf("portwatch: initial scan error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(); err != nil {
				log.Printf("portwatch: scan error: %v", err)
			}
		}
	}
}

// tick performs a single scan cycle: scan → diff → alert → persist.
func (w *Watcher) tick() error {
	current, err := w.scanner.Scan()
	if err != nil {
		return err
	}

	prev, err := w.store.Load()
	if err != nil {
		return err
	}

	diffs := scanner.Compare(prev.Ports, current)
	if len(diffs) > 0 {
		w.alerter.Notify(diffs)
	}

	return w.store.Save(current)
}
