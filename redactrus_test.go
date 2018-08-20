package redactrus

import (
	"regexp"
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

type fireTest struct {
	name           string
	redactionList  []string
	logEntry       *logrus.Entry
	logFields      logrus.Fields
	logDescription string
	expectErr      bool // set to true when an error is expected from Fire()
	description    string
}

func TestFire(t *testing.T) {
	tests := []fireTest{
		{
			name:          "invalidRegex",
			redactionList: []string{"\\"},
			logEntry:      &logrus.Entry{},
			expectErr:     true,
			description:   "Fire() was expected to return an error but did not.",
		},
		{
			name:          "redactWithKey",
			redactionList: []string{"secretKey"},
			logEntry: &logrus.Entry{
				Data: logrus.Fields{
					"secretKey": "secret!",
				},
			},
			expectErr:   false,
			description: "secretKey was expected to be redacted but was not.",
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			e := test.logEntry
			h = &Hook{RedactionList: test.redactionList}
			err := h.Fire(e)

			if test.expectErr {
				assert.NotNil(t, err, test.description)
				return
			}

			// Test all key:val pairs in the logrus Entry. Any value in which
			// the key matches a string in the redaction list should be redacted.
			for k, v := range e.Data {
				for _, s := range test.redactionList {
					re := regexp.MustCompile(s)
					if re.MatchString(k) {
						assert.Equal(t, "[REDACTED]", v, test.description)
					}
				}
			}

		}
		t.Run(test.name, fn)
	}

}
