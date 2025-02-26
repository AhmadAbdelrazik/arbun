package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Level uint8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Log struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Log {
	return &Log{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Log) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Log) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Log) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Log) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := struct {
		Level      string            `json:"level"`
		Message    string            `json:"message"`
		Time       time.Time         `json:"time"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Message:    message,
		Time:       time.Now(),
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	ss, err := json.Marshal(aux)
	if err != nil {
		return 0, err
	}

	l.mu.Lock()

	l.out.Write(append(ss, '\n'))

	l.mu.Unlock()

	return 0, nil
}

func (l *Log) Write(p []byte) (int, error) {
	return l.print(LevelError, string(p), nil)
}
