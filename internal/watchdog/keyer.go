package watchdog

import "fmt"

// Key returns a string key identifying a watchdog instance by name,
// useful when multiple watchdogs run concurrently (e.g. per-interface).
func Key(name string) string {
	return fmt.Sprintf("watchdog:%s", name)
}
