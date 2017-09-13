package logx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugLoggerMessage(t *testing.T) {
	marshaler := marshalerMock()
	logField := NewField("user_id", 12345)
	logMessage := "this is a debug message"

	writer := writerMock(t, logMessage)
	logger := newLog(marshaler, writer, LevelDebug)

	logger.Debug(logMessage, logField)
	assert.Len(t, marshaler.calls.Marshal, 1)
}

func TestDoesNotLogForUnexpectedLevel(t *testing.T) {
	marshaler := marshalerMock()
	logField := NewField("user_id", 12345)
	logMessage := "foo message"

	writer := writerMock(t, logMessage)
	logger := newLog(marshaler, writer, LevelInfo)

	logger.Debug(logMessage, logField)
	assert.Len(t, marshaler.calls.Marshal, 0)
}

func marshalerMock() *MarshalerMock {
	marshaler := MarshalerMock{}
	marshaler.MarshalFunc = func(e *Entry) ([]byte, error) {
		return []byte(e.message), nil
	}

	return &marshaler
}

func writerMock(t *testing.T, expectedMessage string) *WriterMock {
	mockWriter := WriterMock{}
	mockWriter.WriteFunc = func(data []byte) (int, error) {
		assert.Equal(t, expectedMessage, string(data))

		return 0, nil
	}

	return &mockWriter
}
