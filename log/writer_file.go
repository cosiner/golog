package log

import (
	"bufio"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// log dir permission when create
	logDirPerm     = 0755
	logFileNameFmt = "20060102"
)

type (
	// logBuffer represent a log w for a special level
	logBuffer struct {
		file *os.File
		*bufio.Writer
	}
)

func (buf *logBuffer) init(name string, bufsize int) error {
	fd, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if buf.file != nil {
		buf.Writer.Flush()
		buf.file.Close()

		buf.Writer.Reset(fd)
	} else {
		buf.Writer = bufio.NewWriterSize(fd, bufsize)
	}
	buf.file = fd

	return nil
}

// close close the log buffer
func (buf *logBuffer) Close() {
	if buf.file != nil {
		buf.Flush()
		buf.file.Close()
	}
}

type FileWriter struct {
	level   Level
	logdir  string
	bufsize int
	day     int
	files   []logBuffer

	singleOutput bool
	mu           sync.RWMutex
}

func NewFileWriter(logLevel Level, logdir string, bufsize int, singleOutput bool) (*FileWriter, error) {
	err := os.MkdirAll(logdir, logDirPerm)
	if err != nil {
		return nil, err
	}

	w := &FileWriter{
		level:   logLevel,
		logdir:  logdir,
		bufsize: bufsize,
		day:     -1,
		files:   make([]logBuffer, levelMax+1),

		singleOutput: singleOutput,
	}
	w.checkDaily()

	return w, nil
}

// Write write log to log file, higher level log will simultaneously
// output to all lower level log file
func (w *FileWriter) Write(level Level, bytes []byte) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.checkDaily()
	l := w.level
	if w.singleOutput {
		l = level
	}
	for ; l <= level && err == nil; l++ {
		_, err = w.files[l].Write(bytes)
	}

	return
}

func (w *FileWriter) checkDaily() {
	now := time.Now()
	datetime := now.Format(logFileNameFmt)
	if d := now.Day(); d != w.day {
		w.day = d
		for l := w.level; l <= levelMax; l++ {
			w.files[l].init(w.logfileName(l, datetime), w.bufsize)
		}
	}
}

func (w *FileWriter) Flush() {
	for l := w.level; l <= levelMax; l++ {
		w.files[l].Flush()
	}
}

func (w *FileWriter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for l := w.level; l <= levelMax; l++ {
		w.files[l].Close()
	}
}

func (w *FileWriter) logfileName(level Level, datetime string) string {
	return filepath.Join(w.logdir, level.String()+"."+datetime+".log")
}
