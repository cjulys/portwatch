package suppress

import "time"

// MatchResult describes why a suppression rule did or did not apply.
type MatchResult struct {
	Matched  bool
	RuleKey  string
	ExpiresAt time.Time
}

// Match checks whether any active rule in s suppresses the given
// proto/addr/port combination. It checks the exact key first, then
// the wildcard key. The zero MatchResult (Matched == false) is
// returned when no rule applies.
func (s *Suppressor) Match(proto, addr string, port uint16) MatchResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	candidates := []string{
		Key(proto, addr, port),
		WildcardKey(proto, port),
	}

	for _, k := range candidates {
		r, ok := s.rules[k]
		if !ok {
			continue
		}
		if !r.Until.IsZero() && now.After(r.Until) {
			delete(s.rules, k)
			continue
		}
		return MatchResult{Matched: true, RuleKey: k, ExpiresAt: r.Until}
	}

	return MatchResult{}
}
