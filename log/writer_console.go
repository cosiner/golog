package log

import (
	"os"
	"sync"
)

type consoleWriter struct {
	mu sync.Mutex
}

func Console() Writer {
	return &consoleWriter{}
}

func (w *consoleWriter) Write(level Level, bytes []byte) error {
	w.mu.Lock()
	out := os.Stdout
	if level >= LevelError {
		out = os.Stderr
	}
	_, err := out.Write(bytes)
	w.mu.Unlock()
	return err
}

func (w *consoleWriter) Flush() {}
func (w *consoleWriter) Close() {}
