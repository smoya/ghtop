package logx

import (
	"errors"
	"io"
	"io/ioutil"
	stdlog "log"
	"os"
	"time"
)

const (
	// LevelDebug info required while troubleshooting issues.
	LevelDebug Level = iota + 1

	// LevelInfo info required for future post-mortems.
	LevelInfo

	// Do you miss LevelError? Take a look to https://dave.cheney.net/2015/11/05/lets-talk-about-logging.
)

//Entry is the structure used for logging.
type Entry struct {
	message string
	fields  []Field
	level   Level
	time    time.Time
}

// Field is a key/value pair associated to a log.
type Field struct {
	Key   string
	Value interface{}
}

// NewField creates a new field.
func NewField(key string, value interface{}) Field {
	return Field{
		key,
		value,
	}
}

// Level represents the level of logging.
type Level uint8

func (l Level) String() (string, error) {
	switch l {
	case LevelInfo:
		return "info", nil
	case LevelDebug:
		return "debug", nil
	}

	return "", errors.New("invalid log level")
}

// A Logger enables leveled, structured logging. All methods are safe for concurrent use.
type Logger interface {
	Info(string, ...Field)
	Debug(string, ...Field)
}

// Log implements Logger.
type log struct {
	marshaler      Marshaler
	writer         io.Writer
	thresholdLevel Level
}

// NewDiscardAll creates a logger for testing purposes.
func NewDiscardAll() Logger {
	return newLog(
		newNOOPMarshaller(),
		ioutil.Discard,
		LevelInfo,
	)
}

// NewStdOut creates a logger which prints for the standard output.
func NewStdOut(channel, application, environment string, lvl Level) (Logger, error) {
	return newLog(
			newHumanMarshaller(channel, application, environment),
			os.Stdout,
			lvl,
		),
		nil
}

func newLog(marshaler Marshaler, writer io.Writer, thresholdLevel Level) Logger {
	return &log{
		marshaler,
		writer,
		thresholdLevel,
	}
}

// Info logs data with level Info.
func (l *log) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

// Debug data with level Debug.
func (l *log) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

func (l *log) log(currentLevel Level, msg string, fields ...Field) {
	if currentLevel < l.thresholdLevel {
		return
	}

	entry := &Entry{
		message: msg,
		fields:  fields,
		level:   currentLevel,
		time:    time.Now(),
	}
	data, err := l.marshaler.Marshal(entry)

	if err == nil {
		_, err = l.writer.Write(data)
	}
	if err != nil {
		stdlog.Println(err)
	}
}
