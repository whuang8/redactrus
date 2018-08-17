package redactrus

import (
	"github.com/Sirupsen/logrus"
)

// Replacement holds the definitions of what to redact.
// Key: A string or pattern to redact
// Value: The text used to redact the key
type Replacements map[string]string

// Hook is a logrus hook for redacting information from logs
type Hook struct {
	// Messages with a log level not contained in this array
	// will not be dispatched. If empty, all messages will be dispatched.
	AcceptedLevels []logrus.Level
	Replacements   Replacements
}

// Levels ...
func (h *Hook) Levels() []logrus.Level {
	if len(h.AcceptedLevels) == 0 {
		return logrus.AllLevels
	}
	return h.AcceptedLevels
}

// LevelThreshold returns a []logrus.Level including all levels
// above and including the level given. If the provided level does not exit,
// an empty slice is returned.
func LevelThreshold(l logrus.Level) []logrus.Level {
	if l < 0 || int(l) > len(logrus.AllLevels) {
		return []logrus.Level{}
	}
	return logrus.AllLevels[:l+1]
}

// Fire ...
func (h *Hook) Fire(e *logrus.Entry) error {
	e.Message = "[REDACTED]"
	return nil
}
