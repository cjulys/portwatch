// Package portrank assigns a numeric priority rank to ports based on
// well-known service criticality. Higher rank means more critical.
package portrank

import "fmt"

// Rank represents the criticality of a port (0 = unknown, higher = more critical).
type Rank int

const (
	RankUnknown  Rank = 0
	RankLow      Rank = 1
	RankMedium   Rank = 2
	RankHigh     Rank = 3
	RankCritical Rank = 4
)

func (r Rank) String() string {
	switch r {
	case RankCritical:
		return "critical"
	case RankHigh:
		return "high"
	case RankMedium:
		return "medium"
	case RankLow:
		return "low"
	default:
		return "unknown"
	}
}

// Ranker assigns ranks to ports.
type Ranker struct {
	overrides map[string]Rank
}

// builtinRanks maps "proto:port" to a Rank.
var builtinRanks = map[string]Rank{
	"tcp:22":   RankCritical, // SSH
	"tcp:23":   RankCritical, // Telnet
	"tcp:80":   RankHigh,     // HTTP
	"tcp:443":  RankHigh,     // HTTPS
	"tcp:3306": RankHigh,     // MySQL
	"tcp:5432": RankHigh,     // PostgreSQL
	"tcp:6379": RankMedium,   // Redis
	"tcp:8080": RankMedium,   // Alt HTTP
	"udp:53":   RankHigh,     // DNS
	"udp:161":  RankMedium,   // SNMP
}

// New returns a Ranker with optional overrides.
func New(overrides map[string]Rank) *Ranker {
	if overrides == nil {
		overrides = make(map[string]Rank)
	}
	return &Ranker{overrides: overrides}
}

// Get returns the rank for the given protocol and port number.
func (r *Ranker) Get(proto string, port int) Rank {
	k := fmt.Sprintf("%s:%d", proto, port)
	if rank, ok := r.overrides[k]; ok {
		return rank
	}
	if rank, ok := builtinRanks[k]; ok {
		return rank
	}
	return RankUnknown
}

// IsCritical returns true when the rank is RankCritical.
func (r *Ranker) IsCritical(proto string, port int) bool {
	return r.Get(proto, port) == RankCritical
}
