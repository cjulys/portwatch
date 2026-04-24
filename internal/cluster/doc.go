// Package cluster provides port-range grouping for portwatch.
//
// When many consecutive ports change state simultaneously — for example during
// a service restart that binds a block of ephemeral ports — individual per-port
// alerts become noisy. cluster.Group collapses those ports into compact Range
// values (e.g. "tcp/8080-8090") that can be forwarded to alert and reporter
// pipelines without modification.
package cluster
