package log

import (
	"io"
	"os"
	"sync"

	"github.com/cosiner/gohper/terminal/color"
)

type ConsoleWriter struct {
	colors []*color.Renderer

	stdout io.Writer
	stderr io.Writer

	mu sync.Mutex
}

func NewConsoleWriter(mappings map[Level][]color.Code) *ConsoleWriter {
	defaultMappings := []*color.Renderer{
		color.New(color.Highlight, color.FgWhite),   //LEVEL_DEBUG
		color.New(color.Highlight, color.FgGreen),   //LEVEL_INFO
		color.New(color.Highlight, color.FgYellow),  //LEVEL_WARN
		color.New(color.Highlight, color.FgMagenta), //LEVEL_ERROR
		color.New(color.Highlight, color.FgRed),     //LEVEL_PANIC
		color.New(color.Highlight, color.FgRed),     //LEVEL_FATAL
	}

	for l, m := range mappings {
		if levelMin < l && l <= levelMax {
			defaultMappings[l] = color.New(m...)
		}
	}

	return &ConsoleWriter{
		colors: defaultMappings,
		stdout: color.Stdout,
		stderr: color.Stderr,
	}
}

func (w *ConsoleWriter) DisableColor() *ConsoleWriter {
	w.colors = nil
	w.stdout = os.Stdout
	w.stderr = os.Stderr
	return w
}

func (w *ConsoleWriter) Write(level Level, bytes []byte) error {
	w.mu.Lock()
	out := w.stdout
	if level >= LevelError {
		out = w.stderr
	}
	if w.colors == nil {
		_, err := out.Write(bytes)
		w.mu.Unlock()
		return err
	}

	tc := w.colors[level]
	tc.Begin(out)
	_, err := out.Write(bytes)
	tc.End(out)
	w.mu.Unlock()
	return err
}

func (w *ConsoleWriter) Flush() {}
func (w *ConsoleWriter) Close() {}
