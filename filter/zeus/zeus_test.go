package zeus

import (
	"../../utils/logevent"
	"testing"
	"time"
)

func Test_DefaultFilterConfig(t *testing.T) {
	DefaultFilterConfig()
}

func Test_InitHandler(t *testing.T) {
	// InitHandler(&make(map[string]interface{}))
}

func Test_Event(t *testing.T) {
	le := logevent.LogEvent{
		Timestamp: time.Now(),
		Message:   "message",
		Tags:      []string{"frg", "grbhrt"},
		Extra:     make(map[string]interface{}),
	}
	config := DefaultFilterConfig()
	config.Event(le)
}
