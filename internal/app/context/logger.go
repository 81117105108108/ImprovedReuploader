package context

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/kartFr/Asset-Reuploader/internal/color"
)

type logger struct {
	History *bytes.Buffer
	writer  io.Writer
	mu      sync.Mutex
}

func newLogger() *logger {
	var b bytes.Buffer
	w := io.MultiWriter(
		&b,
		color.Output,
	)

	return &logger{
		History: &b,
		writer:  w,
	}
}

func (l *logger) Error(a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.Error.Fprintln(l.writer, a...)
}

func (l *logger) Info(a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.Info.Fprintln(l.writer, a...)
}

func (l *logger) Println(a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.writer, a...)
}

func (l *logger) Success(a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.Success.Fprintln(l.writer, a...)
}

func (l *logger) Warn(a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.Warn.Fprintln(l.writer, a...)
}

// HistoryString returns the history buffer under lock.
func (l *logger) HistoryString() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.History.String()
}

// Historyf appends to history only (used while paused, no console duplicate).
func (l *logger) Historyf(c *color.Color, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	c.Fprintln(l.History, message)
}
