// Package probe performs lightweight TCP/UDP connectivity checks
// against a port to confirm it is truly accepting connections.
package probe

import (
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single probe attempt.
type Result struct {
	Port     int
	Protocol string
	Reachable bool
	Latency  time.Duration
	Err      error
}

// Prober performs connectivity probes.
type Prober struct {
	timeout time.Duration
}

// New returns a Prober with the given dial timeout.
func New(timeout time.Duration) *Prober {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Prober{timeout: timeout}
}

// Probe attempts to connect to the given host:port using protocol.
// Only "tcp" and "tcp6" are supported; all others return an error result.
func (p *Prober) Probe(host string, port int, protocol string) Result {
	switch protocol {
	case "tcp", "tcp4", "tcp6":
	default:
		return Result{
			Port:     port,
			Protocol: protocol,
			Err:      fmt.Errorf("probe: unsupported protocol %q", protocol),
		}
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, err := net.DialTimeout(protocol, addr, p.timeout)
	latency := time.Since(start)

	if err != nil {
		return Result{Port: port, Protocol: protocol, Reachable: false, Latency: latency, Err: err}
	}
	_ = conn.Close()
	return Result{Port: port, Protocol: protocol, Reachable: true, Latency: latency}
}

// ProbeAll probes a list of port/protocol pairs on host concurrently.
func (p *Prober) ProbeAll(host string, targets []Target) []Result {
	results := make([]Result, len(targets))
	ch := make(chan indexed, len(targets))
	for i, t := range targets {
		go func(idx int, t Target) {
			ch <- indexed{idx, p.Probe(host, t.Port, t.Protocol)}
		}(i, t)
	}
	for range targets {
		v := <-ch
		results[v.i] = v.r
	}
	return results
}

// Target is a port+protocol pair to probe.
type Target struct {
	Port     int
	Protocol string
}

type indexed struct {
	i int
	r Result
}
