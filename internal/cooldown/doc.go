// Package cooldown suppresses repeated port-change events that fire
// continuously within a configurable quiet window. It complements
// internal/throttle by resetting the timer on every new event rather than
// allowing a fixed burst rate, making it suitable for flapping ports that
// should only generate an alert once they have settled.
package cooldown
