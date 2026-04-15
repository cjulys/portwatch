package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/watcher"
)

func main() {
	cfgPath := flag.String("config", "", "path to config file (JSON)")
	statePath := flag.String("state", "/tmp/portwatch_state.json", "path to state file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: failed to load config: %v\n", err)
		os.Exit(1)
	}

	if errs := config.Validate(cfg); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "portwatch: config error: %v\n", e)
		}
		os.Exit(1)
	}

	sc := scanner.New(cfg)
	st := state.New(*statePath)
	al := alert.New(cfg, os.Stdout)
	w := watcher.New(cfg, sc, st, al)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("portwatch: starting (interval=%ds, state=%s)", cfg.IntervalSeconds, *statePath)
	if err := w.Run(ctx); err != nil && err != context.Canceled {
		log.Printf("portwatch: exited with error: %v", err)
	}
	log.Println("portwatch: stopped")
}
