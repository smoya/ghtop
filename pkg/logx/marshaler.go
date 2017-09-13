package logx

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()
}

// logInfo  ...
type logInfo struct {
	channel     string
	application string
	environment string
	hostname    string
}

// Marshaler marshals
type Marshaler interface {
	// Marshal returns the info encoded to be readable by humans
	Marshal(entry *Entry) ([]byte, error)
}

type noopMarshaller struct{}

func newNOOPMarshaller() Marshaler {
	return &noopMarshaller{}

}

// Marshal  ...
func (n *noopMarshaller) Marshal(entry *Entry) ([]byte, error) {
	return nil, nil
}

type humanMarshaller struct {
	logInfo
}

// newHumanMarshaller ...
func newHumanMarshaller(channel, application, environment string) Marshaler {
	return &humanMarshaller{
		logInfo{
			channel:     channel,
			application: application,
			environment: environment,
			hostname:    hostname,
		},
	}
}

// Marshal returns the info encoded to be readable by humans.
func (m *humanMarshaller) Marshal(entry *Entry) ([]byte, error) {
	separator := ", "
	var buffer bytes.Buffer

	lvl, _ := entry.level.String()
	buffer.WriteString(fmt.Sprintf("[%v]%v", entry.time.Format(time.RFC3339), separator))
	buffer.WriteString(fmt.Sprintf("%v.%v%v", m.application, m.environment, separator))
	buffer.WriteString(fmt.Sprintf("%v.%v%v", m.channel, lvl, separator))
	buffer.WriteString(fmt.Sprintf("%v%v", entry.message, separator))
	buffer.WriteString("[")
	// rest
	for i, field := range entry.fields {
		buffer.WriteString(fmt.Sprintf("%v:%v", field.Key, field.Value))
		if i != len(entry.fields)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]\n")

	return buffer.Bytes(), nil
}
