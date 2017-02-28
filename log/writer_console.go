package log

import (
	"io"
	"os"
	"sync"
)

type ConsoleWriter struct {
	stdout io.Writer
	stderr io.Writer

	mu sync.Mutex
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (w *ConsoleWriter) Write(level Level, bytes []byte) error {
	w.mu.Lock()
	out := w.stdout
	if level >= LevelError {
		out = w.stderr
	}
	_, err := out.Write(bytes)
	w.mu.Unlock()
	return err
}

func (w *ConsoleWriter) Flush() {}
func (w *ConsoleWriter) Close() {}
