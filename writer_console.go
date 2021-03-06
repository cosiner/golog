package golog

import (
	"os"
	"sync"
)

type ConsoleColorRender func(Level, []byte) []byte

type consoleWriter struct {
	mu sync.Mutex

	renders     []ConsoleColorRender
	renderCount uint8
	renderCurr  uint8
}

func Console(renders ...ConsoleColorRender) Writer {
	w := consoleWriter{
		renderCount: uint8(len(renders)),
		renders:     renders,
	}
	return &w
}

func (w *consoleWriter) Write(level Level, bytes []byte) error {
	w.mu.Lock()
	if w.renderCount > 0 {
		render := w.renders[w.renderCurr%w.renderCount]
		w.renderCurr++
		bytes = render(level, bytes)
	}
	_, err := os.Stderr.Write(bytes)
	w.mu.Unlock()
	return err
}

func (w *consoleWriter) Flush() {}
func (w *consoleWriter) Close() {}
