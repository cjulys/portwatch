package throttle

import (
	"fmt"

	"portwatch/internal/scanner"
)

// PortKey returns a stable string key for a scanned port,
// suitable for use as a Throttle key.
func PortKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Number)
}

// DiffKey returns a key that encodes both the port identity and the
// direction of change ("opened" or "closed"), so that open and close
// events for the same port are throttled independently.
func DiffKey(p scanner.Port, opened bool) string {
	direction := "closed"
	if opened {
		direction = "opened"
	}
	return fmt.Sprintf("%s:%d:%s", p.Protocol, p.Number, direction)
}
