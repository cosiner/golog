package golog

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type multiFileWriter struct {
	level Level
	opts  FileLogOptions

	files    []buffedFile
	day      int
	filedate string
	mu       sync.Mutex
}

func MultiFile(logLevel Level, options ...FileLogOptions) (Writer, error) {
	opts := newDefaultFileLogOptions(options...)
	err := os.MkdirAll(opts.LogDir, logDirPerm)
	if err != nil {
		return nil, err
	}

	w := &multiFileWriter{
		level: logLevel,
		opts:  opts,
		day:   -1,
		files: make([]buffedFile, levelMax+1),
	}
	w.checkDaily()

	return w, nil
}

func (w *multiFileWriter) Write(level Level, bytes []byte) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.checkDaily()
	for l := w.level; l <= level && err == nil; l++ {
		_, err = w.files[l].Write(bytes)
	}
	return
}

func (w *multiFileWriter) checkDaily() {
	now := time.Now()
	if d := now.Day(); d != w.day {
		w.day = d
		datetime := now.Format(logFileDateFmt)

		var err error
		for l := w.level; l <= levelMax; l++ {
			e := w.files[l].init(w.logfileName(l, datetime), w.opts.Bufsize)
			if e != nil {
				err = e
			}
		}
		if err == nil {
			w.filedate = datetime
		}
		cleanLogFiles(&w.opts, w.filedate, w.parseLogDate)
	}
}

func (w *multiFileWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for l := w.level; l <= levelMax; l++ {
		w.files[l].Flush()
	}
}

func (w *multiFileWriter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for l := w.level; l <= levelMax; l++ {
		w.files[l].Close()
	}
}

func (w *multiFileWriter) logfileName(level Level, datetime string) string {
	return filepath.Join(w.opts.LogDir, level.String()+"."+datetime+".log")
}

func (w *multiFileWriter) parseLogDate(filename string) (string, bool) {
	filename = filepath.Base(filename)
	secs := strings.Split(filename, ".")
	if len(secs) != 3 || secs[2] != "log" || secs[0] == "" || secs[1] == "" {
		return "", false
	}
	return secs[1], true
}
