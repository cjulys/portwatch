package escalate_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/escalate"
)

func TestRepeatEscalationCycle(t *testing.T) {
	e := escalate.New(time.Minute, 2, 3)
	levels := make([]escalate.Level, 4)
	for i := range levels {
		levels[i] = e.Evaluate("port:80")
	}
	if levels[0] != escalate.Info {
		t.Errorf("step 0: want Info got %s", levels[0])
	}
	if levels[1] != escalate.Warning {
		t.Errorf("step 1: want Warning got %s", levels[1])
	}
	if levels[2] != escalate.Critical {
		t.Errorf("step 2: want Critical got %s", levels[2])
	}
	if levels[3] != escalate.Critical {
		t.Errorf("step 3: want Critical got %s", levels[3])
	}
}

func TestFlushThenReEscalate(t *testing.T) {
	e := escalate.New(time.Minute, 2, 3)
	e.Evaluate("port:443")
	e.Evaluate("port:443")
	e.Flush()
	l := e.Evaluate("port:443")
	if l != escalate.Info {
		t.Fatalf("want Info after flush, got %s", l)
	}
}
