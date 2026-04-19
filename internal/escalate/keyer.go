package escalate

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// PortKey returns a stable string key for a scanned port.
// The key format is "<protocol>:<number>", e.g. "tcp:443".
func PortKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Number)
}

// PortKeyFromParts returns a stable string key given a protocol and port number
// directly, without requiring a scanner.Port value.
func PortKeyFromParts(protocol string, number int) string {
	return fmt.Sprintf("%s:%d", protocol, number)
}
