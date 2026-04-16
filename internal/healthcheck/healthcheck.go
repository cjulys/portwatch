// Package healthcheck exposes a simple HTTP endpoint reporting daemon health.
package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds the current health snapshot.
type Status struct {
	OK        bool      `json:"ok"`
	Uptime    string    `json:"uptime"`
	Scans     uint64    `json:"scans_total"`
	LastScan  time.Time `json:"last_scan"`
	StartedAt time.Time `json:"started_at"`
}

// Server is a lightweight HTTP health server.
type Server struct {
	addr      string
	scans     atomic.Uint64
	lastScan  atomic.Pointer[time.Time]
	startedAt time.Time
	server    *http.Server
}

// New creates a Server that listens on addr (e.g. ":9090").
func New(addr string) *Server {
	s := &Server{addr: addr, startedAt: time.Now().UTC()}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	s.server = &http.Server{Addr: addr, Handler: mux}
	return s
}

// RecordScan updates the scan counter and timestamp.
func (s *Server) RecordScan() {
	s.scans.Add(1)
	now := time.Now().UTC()
	s.lastScan.Store(&now)
}

// Start begins serving in a background goroutine.
func (s *Server) Start() {
	go func() { _ = s.server.ListenAndServe() }()
}

// Stop shuts the HTTP server down.
func (s *Server) Stop() error { return s.server.Close() }

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	uptime := time.Since(s.startedAt).Truncate(time.Second)
	var last time.Time
	if p := s.lastScan.Load(); p != nil {
		last = *p
	}
	status := Status{
		OK:        true,
		Uptime:    fmt.Sprintf("%s", uptime),
		Scans:     s.scans.Load(),
		LastScan:  last,
		StartedAt: s.startedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
