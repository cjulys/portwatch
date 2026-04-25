// Package normalize provides port normalization utilities that canonicalize
// raw scanner output before it is compared, stored, or alerted on.
package normalize

import (
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Option configures a Normalizer.
type Option func(*Normalizer)

// WithLowerProtocol lowercases the Protocol field of every port.
func WithLowerProtocol() Option {
	return func(n *Normalizer) { n.lowerProto = true }
}

// WithTrimAddress strips leading/trailing whitespace from the Address field.
func WithTrimAddress() Option {
	return func(n *Normalizer) { n.trimAddr = true }
}

// WithDeduplication removes duplicate ports from the slice (same port+protocol).
func WithDeduplication() Option {
	return func(n *Normalizer) { n.dedup = true }
}

// Normalizer applies a configurable set of transformations to a port slice.
type Normalizer struct {
	lowerProto bool
	trimAddr   bool
	dedup      bool
}

// New returns a Normalizer configured with the supplied options.
func New(opts ...Option) *Normalizer {
	n := &Normalizer{}
	for _, o := range opts {
		o(n)
	}
	return n
}

// Apply returns a new slice with all configured transformations applied.
// The original slice is never modified.
func (n *Normalizer) Apply(ports []scanner.Port) []scanner.Port {
	out := make([]scanner.Port, 0, len(ports))
	seen := make(map[string]struct{})

	for _, p := range ports {
		if n.lowerProto {
			p.Protocol = strings.ToLower(p.Protocol)
		}
		if n.trimAddr {
			p.Address = strings.TrimSpace(p.Address)
		}
		if n.dedup {
			k := key(p)
			if _, exists := seen[k]; exists {
				continue
			}
			seen[k] = struct{}{}
		}
		out = append(out, p)
	}
	return out
}

func key(p scanner.Port) string {
	return p.Protocol + "/" + itoa(p.Port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
