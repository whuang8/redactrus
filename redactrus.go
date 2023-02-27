package redactrus

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/sirupsen/logrus"
)

// Hook is a logrus hook for redacting information from logs
type Hook struct {
	// Messages with a log level not contained in this array
	// will not be dispatched. If empty, all messages will be dispatched.
	AcceptedLevels []logrus.Level
	RedactionList  []string
}

// Levels returns the user defined AcceptedLevels
// If AcceptedLevels is empty, all logrus levels are returned
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

// Fire redacts values in an log Entry that match
// with keys defined in the RedactionList
func (h *Hook) Fire(e *logrus.Entry) error {
	for _, redactionKey := range h.RedactionList {
		re, err := regexp.Compile(redactionKey)
		if err != nil {
			return err
		}

		// Redact based on key matching in Data fields
		for k, v := range e.Data {
			if re.MatchString(k) {
				e.Data[k] = "[REDACTED]"
				continue
			}

			// Logrus Field values can be nil
			if v == nil {
				continue
			}

			// Redact based on value matching in Data fields
			switch reflect.TypeOf(v).Kind() {
			case reflect.String:
				switch vv := v.(type) {
				case string:
					e.Data[k] = re.ReplaceAllString(vv, "$1[REDACTED]$2")
				default:
					e.Data[k] = re.ReplaceAllString(fmt.Sprint(v), "$1[REDACTED]$2")
				}
				continue
			// prevent nil *fmt.Stringer from reaching handler below
			case reflect.Ptr:
				if reflect.ValueOf(v).IsNil() {
					continue
				}
			}

			// Handle fmt.Stringer type.
			if vv, ok := v.(fmt.Stringer); ok {
				e.Data[k] = re.ReplaceAllString(vv.String(), "$1[REDACTED]$2")
				continue
			}

		}

		// Redact based on text matching in the Message field
		e.Message = re.ReplaceAllString(e.Message, "$1[REDACTED]$2")
	}

	return nil
}
