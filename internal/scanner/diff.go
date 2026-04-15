package scanner

// Diff computes the changes between two port snapshots.
type Diff struct {
	Opened []PortState
	Closed []PortState
}

// HasChanges returns true if there are any opened or closed ports.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Compare returns a Diff between a previous and current set of PortStates.
func Compare(previous, current []PortState) Diff {
	prevMap := toMap(previous)
	currMap := toMap(current)

	var diff Diff

	for key, ps := range currMap {
		if _, exists := prevMap[key]; !exists {
			diff.Opened = append(diff.Opened, ps)
		}
	}

	for key, ps := range prevMap {
		if _, exists := currMap[key]; !exists {
			diff.Closed = append(diff.Closed, ps)
		}
	}

	return diff
}

func toMap(states []PortState) map[string]PortState {
	m := make(map[string]PortState, len(states))
	for _, s := range states {
		key := s.Protocol + s.Address + string(rune(s.Port))
		m[key] = s
	}
	return m
}
