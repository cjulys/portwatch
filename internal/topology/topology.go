// Package topology builds a lightweight network topology view by grouping
// open ports by host address and protocol family. It is useful for producing
// structured summaries and for detecting when a single host suddenly exposes
// an unusual number of ports.
package topology

import (
	"fmt"
	"sort"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Host represents a single observed host and the ports it currently exposes.
type Host struct {
	Address  string
	Ports    []scanner.Port
	Protocols map[string]int // protocol -> count
}

// PortCount returns the total number of open ports on this host.
func (h *Host) PortCount() int { return len(h.Ports) }

// Topology holds the current per-host view derived from a port scan result.
type Topology struct {
	mu    sync.RWMutex
	hosts map[string]*Host // keyed by address
}

// New returns an empty Topology ready for use.
func New() *Topology {
	return &Topology{
		hosts: make(map[string]*Host),
	}
}

// Update replaces the internal state with a fresh topology built from ports.
// It is safe to call concurrently.
func (t *Topology) Update(ports []scanner.Port) {
	next := build(ports)

	t.mu.Lock()
	t.hosts = next
	t.mu.Unlock()
}

// Hosts returns a stable, sorted slice of all known hosts.
func (t *Topology) Hosts() []*Host {
	t.mu.RLock()
	defer t.mu.RUnlock()

	out := make([]*Host, 0, len(t.hosts))
	for _, h := range t.hosts {
		out = append(out, h)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Address < out[j].Address
	})
	return out
}

// HostCount returns the number of distinct addresses currently tracked.
func (t *Topology) HostCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.hosts)
}

// Get returns the Host for the given address, or nil if not found.
func (t *Topology) Get(address string) *Host {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.hosts[address]
}

// HeavyHitters returns hosts whose open port count exceeds threshold.
func (t *Topology) HeavyHitters(threshold int) []*Host {
	all := t.Hosts()
	out := all[:0]
	for _, h := range all {
		if h.PortCount() > threshold {
			out = append(out, h)
		}
	}
	return out
}

// Summary returns a human-readable one-line description of the topology.
func (t *Topology) Summary() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	total := 0
	for _, h := range t.hosts {
		total += h.PortCount()
	}
	return fmt.Sprintf("%d host(s), %d open port(s)", len(t.hosts), total)
}

// build constructs a fresh address-keyed map from a flat port list.
func build(ports []scanner.Port) map[string]*Host {
	hosts := make(map[string]*Host)
	for _, p := range ports {
		h, ok := hosts[p.Address]
		if !ok {
			h = &Host{
				Address:   p.Address,
				Protocols: make(map[string]int),
			}
			hosts[p.Address] = h
		}
		h.Ports = append(h.Ports, p)
		h.Protocols[p.Protocol]++
	}
	return hosts
}
