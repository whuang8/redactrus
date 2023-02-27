package redactrus

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var h = &Hook{}

type levelsTest struct {
	name           string
	acceptedLevels []logrus.Level
	expected       []logrus.Level
	description    string
}

func TestLevels(t *testing.T) {
	tests := []levelsTest{
		{
			name:           "undefinedAcceptedLevels",
			acceptedLevels: []logrus.Level{},
			expected:       logrus.AllLevels,
			description:    "All logrus levels expected, but did not recieve them",
		},
		{
			name:           "definedAcceptedLevels",
			acceptedLevels: []logrus.Level{logrus.InfoLevel},
			expected:       []logrus.Level{logrus.InfoLevel},
			description:    "Logrus Info level expected, but did not recieve that.",
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			h.AcceptedLevels = test.acceptedLevels
			levels := h.Levels()
			assert.Equal(t, test.expected, levels, test.description)
		}

		t.Run(test.name, fn)
	}
}

type levelThresholdTest struct {
	name        string
	level       logrus.Level
	expected    []logrus.Level
	description string
}

func TestLevelThreshold(t *testing.T) {
	tests := []levelThresholdTest{
		{
			name:        "unknownLogLevel",
			level:       logrus.Level(100),
			expected:    []logrus.Level{},
			description: "An empty Level slice was expected but was not returned",
		},
		{
			name:        "errorLogLevel",
			level:       logrus.ErrorLevel,
			expected:    []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel},
			description: "The panic, fatal, and error levels were expected but were not returned",
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			levels := LevelThreshold(test.level)
			assert.Equal(t, test.expected, levels, test.description)
		}

		t.Run(test.name, fn)
	}
}

func TestInvalidRegex(t *testing.T) {
	e := &logrus.Entry{}
	h = &Hook{RedactionList: []string{"\\"}}
	err := h.Fire(e)

	assert.NotNil(t, err)
}

type EntryDataValuesTest struct {
	name          string
	redactionList []string
	logFields     logrus.Fields
	expected      logrus.Fields
	description   string
}

// Test that any occurence of a redaction pattern
// in the values of the entry's data fields is redacted.
func TestEntryDataValues(t *testing.T) {
	tests := []EntryDataValuesTest{
		{
			name:          "match on key",
			redactionList: []string{"Password"},
			logFields:     logrus.Fields{"Password": "password123!"},
			expected:      logrus.Fields{"Password": "[REDACTED]"},
			description:   "Password value should have been redacted, but was not.",
		},
		{
			name:          "string value",
			redactionList: []string{"William"},
			logFields:     logrus.Fields{"Description": "His name is William"},
			expected:      logrus.Fields{"Description": "His name is [REDACTED]"},
			description:   "William should have been redacted, but was not.",
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			logEntry := &logrus.Entry{
				Data: test.logFields,
			}
			h = &Hook{RedactionList: test.redactionList}
			err := h.Fire(logEntry)

			assert.Nil(t, err)
			assert.Equal(t, test.expected, logEntry.Data)
		}
		t.Run(test.name, fn)
	}
}

// Test that any occurence of a redaction pattern
// in the entry's Message field is redacted.
func TestEntryMessage(t *testing.T) {
	logEntry := &logrus.Entry{
		Message: "Secret Password: password123!",
	}
	h = &Hook{RedactionList: []string{`(Password: ).*`}}
	err := h.Fire(logEntry)

	assert.Nil(t, err)
	assert.Equal(t, "Secret Password: [REDACTED]", logEntry.Message)
}

// Logrus fields can have a nil value so test we handle this edge case
func TestNilField(t *testing.T) {
	logEntry := &logrus.Entry{
		Data: logrus.Fields{"Nil": nil},
	}
	h = &Hook{RedactionList: []string{"foo"}}
	err := h.Fire(logEntry)

	assert.Nil(t, err)
	assert.Equal(t, logrus.Fields{"Nil": nil}, logEntry.Data)
}

type stringerValue struct {
	value string
}

func (v stringerValue) String() string {
	return v.value
}

func TestStringer(t *testing.T) {
	logEntry := &logrus.Entry{
		Data: logrus.Fields{"Stringer": stringerValue{"kind is fmt.Stringer"}},
	}
	h = &Hook{RedactionList: []string{"kind"}}
	err := h.Fire(logEntry)

	assert.Nil(t, err)
	assert.Equal(t, logrus.Fields{"Stringer": "[REDACTED] is fmt.Stringer"}, logEntry.Data)

	var s *stringerValue
	nilStringerEntry := &logrus.Entry{
		Data: logrus.Fields{"Stringer": s},
	}
	err = h.Fire(nilStringerEntry)

	assert.Nil(t, err)
}

type TypedString string

// Logrus fields can have re-typed strings so test we handle this edge case
func TestTypedStringValue(t *testing.T) {
	logEntry := &logrus.Entry{
		Data: logrus.Fields{"TypedString": TypedString("kind is string")},
	}
	h = &Hook{RedactionList: []string{"kind"}}
	err := h.Fire(logEntry)

	assert.Nil(t, err)
	assert.Equal(t, logrus.Fields{"TypedString": "[REDACTED] is string"}, logEntry.Data)
}
