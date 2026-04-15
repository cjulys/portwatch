package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// PortState represents the state of a single open port.
type PortState struct {
	Protocol string
	Port     int
	Address  string
}

func (p PortState) String() string {
	return fmt.Sprintf("%s:%d (%s)", p.Address, p.Port, p.Protocol)
}

// Scanner is responsible for discovering open ports on the host.
type Scanner struct {
	Protocols []string
}

// New returns a Scanner configured to check the given protocols.
func New(protocols []string) *Scanner {
	return &Scanner{Protocols: protocols}
}

// Scan returns all currently open ports by attempting connections.
func (s *Scanner) Scan(portRange [2]int) ([]PortState, error) {
	var results []PortState

	for _, proto := range s.Protocols {
		for port := portRange[0]; port <= portRange[1]; port++ {
			addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
			conn, err := net.Dial(proto, addr)
			if err != nil {
				continue
			}
			conn.Close()

			host, portStr, _ := net.SplitHostPort(addr)
			p, _ := strconv.Atoi(portStr)
			results = append(results, PortState{
				Protocol: strings.ToUpper(proto),
				Port:     p,
				Address:  host,
			})
		}
	}

	return results, nil
}
