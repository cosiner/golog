package log

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type multifileWriter struct {
	level   Level
	logdir  string
	bufsize int

	day   int
	files []buffedFile
	mu    sync.Mutex
}

func Multifile(logLevel Level, logdir string, bufsize int) (Writer, error) {
	err := os.MkdirAll(logdir, logDirPerm)
	if err != nil {
		return nil, err
	}

	w := &multifileWriter{
		level:   logLevel,
		logdir:  logdir,
		bufsize: bufsize,
		day:     -1,
		files:   make([]buffedFile, levelMax+1),
	}
	w.checkDaily()

	return w, nil
}

func (w *multifileWriter) Write(level Level, bytes []byte) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.checkDaily()
	l := w.level
	for ; l <= level && err == nil; l++ {
		_, err = w.files[l].Write(bytes)
	}

	return
}

func (w *multifileWriter) checkDaily() {
	now := time.Now()
	datetime := now.Format(logFileNameFmt)
	if d := now.Day(); d != w.day {
		w.day = d
		for l := w.level; l <= levelMax; l++ {
			w.files[l].init(w.logfileName(l, datetime), w.bufsize)
		}
	}
}

func (w *multifileWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for l := w.level; l <= levelMax; l++ {
		w.files[l].Flush()
	}
}

func (w *multifileWriter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for l := w.level; l <= levelMax; l++ {
		w.files[l].Close()
	}
}

func (w *multifileWriter) logfileName(level Level, datetime string) string {
	return filepath.Join(w.logdir, level.String()+"."+datetime+".log")
}
