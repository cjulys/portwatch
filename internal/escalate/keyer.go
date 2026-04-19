package escalate

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// PortKey returns a stable string key for a scanned port.
func PortKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Number)
}
