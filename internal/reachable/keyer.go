package reachable

import (
	"fmt"
	"strings"
)

// PortKey returns the canonical tracker key for a port.
// proto is normalised to lower-case before constructing the key.
func PortKey(proto string, port int) string {
	return fmt.Sprintf("%s:%d", strings.ToLower(proto), port)
}

// SplitKey parses a key produced by PortKey and returns the protocol
// and port number. It returns an error when the key is malformed.
func SplitKey(key string) (proto string, port int, err error) {
	var p int
	_, err = fmt.Sscanf(key, "%s", &key) // ensure non-empty
	parts := strings.SplitN(key, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("reachable: malformed key %q", key)
	}
	proto = parts[0]
	_, err = fmt.Sscanf(parts[1], "%d", &p)
	if err != nil {
		return "", 0, fmt.Errorf("reachable: malformed port in key %q", key)
	}
	return proto, p, nil
}
