package budget

import "fmt"

// ScanKey returns a budget map key for a named scan target.
// It is used when multiple independent budgets are managed in a map.
func ScanKey(target string) string {
	return fmt.Sprintf("scan:%s", target)
}

// HostKey returns a budget key scoped to a specific host address.
func HostKey(host string) string {
	return fmt.Sprintf("host:%s", host)
}
