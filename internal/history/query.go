package history

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Since returns all entries recorded on or after the given time.
func (h *History) Since(t time.Time) []Entry {
	var out []Entry
	for _, e := range h.entries {
		if !e.Timestamp.Before(t) {
			out = append(out, e)
		}
	}
	return out
}

// Latest returns the most recent entry, or nil if history is empty.
func (h *History) Latest() *Entry {
	if len(h.entries) == 0 {
		return nil
	}
	e := h.entries[len(h.entries)-1]
	return &e
}

// PortSeen reports whether a port with the given number and protocol
// appeared in any recorded entry.
func (h *History) PortSeen(number int, protocol string) bool {
	for _, e := range h.entries {
		for _, p := range e.Ports {
			if p.Number == number && p.Protocol == protocol {
				return true
			}
		}
	}
	return false
}

// UniquePortsEver returns a deduplicated list of all ports ever observed.
func (h *History) UniquePortsEver() []scanner.Port {
	seen := make(map[string]scanner.Port)
	for _, e := range h.entries {
		for _, p := range e.Ports {
			key := p.Protocol + ":" + itoa(p.Number)
			if _, ok := seen[key]; !ok {
				seen[key] = p
			}
		}
	}
	out := make([]scanner.Port, 0, len(seen))
	for _, p := range seen {
		out = append(out, p)
	}
	return out
}

// itoa converts an int to its decimal string representation
// without importing strconv to keep the package lightweight.
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
