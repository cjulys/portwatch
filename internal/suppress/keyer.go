package suppress

import "fmt"

// Key returns a canonical string key for a suppression rule
// based on protocol, address, and port number.
func Key(proto, addr string, port uint16) string {
	return fmt.Sprintf("%s:%s:%d", proto, addr, port)
}

// WildcardKey returns a key that matches any address for the given
// protocol and port — used when addr is empty or "*".
func WildcardKey(proto string, port uint16) string {
	return Key(proto, "*", port)
}
