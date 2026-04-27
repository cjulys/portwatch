// Package sieve provides a lightweight probabilistic duplicate-suppression
// filter for portwatch scan events.
//
// The Sieve uses a fixed-size bit array hashed with FNV-32a. It offers
// O(1) insert and lookup at the cost of a configurable false-positive rate.
//
// Typical use: create one Sieve per scan cycle, call TestAndSet for each
// diff key, and Reset at the start of every new cycle.
package sieve
