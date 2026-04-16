// Package healthcheck provides a lightweight HTTP server that exposes a
// /health endpoint for the portwatch daemon.
//
// The endpoint returns a JSON payload containing uptime, total scans
// performed, and the timestamp of the most recent scan. It is intended
// for use with external monitoring systems or simple liveness probes.
package healthcheck
