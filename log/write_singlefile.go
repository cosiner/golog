package log

import (
	"bufio"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	logDirPerm     = 0755
	logFileNameFmt = "20060102"
)

type (
	buffedFile struct {
		file *os.File
		*bufio.Writer
	}
)

func (bf *buffedFile) init(name string, bufsize int) error {
	fd, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if bf.file != nil {
		bf.Writer.Flush()
		bf.file.Close()

		bf.Writer.Reset(fd)
	} else {
		bf.Writer = bufio.NewWriterSize(fd, bufsize)
	}
	bf.file = fd

	return nil
}

func (bf *buffedFile) Close() {
	if bf.file != nil {
		bf.Flush()
		bf.file.Close()
	}
}

type singleWriter struct {
	level   Level
	logdir  string
	bufsize int
	day     int
	file    buffedFile

	mu sync.Mutex
}

func Singlefile(logLevel Level, logdir string, bufsize int) (Writer, error) {
	err := os.MkdirAll(logdir, logDirPerm)
	if err != nil {
		return nil, err
	}

	w := &singleWriter{
		level:   logLevel,
		logdir:  logdir,
		bufsize: bufsize,
		day:     -1,
	}
	w.checkDaily()

	return w, nil
}

func (w *singleWriter) Write(level Level, bytes []byte) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.checkDaily()
	_, err = w.file.Write(bytes)
	return
}

func (w *singleWriter) checkDaily() {
	now := time.Now()
	if d := now.Day(); d != w.day {
		w.day = d
		w.file.init(w.logfileName(now.Format(logFileNameFmt)), w.bufsize)
	}
}

func (w *singleWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.file.Flush()
}

func (w *singleWriter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.file.Close()
}

func (w *singleWriter) logfileName(datetime string) string {
	return filepath.Join(w.logdir, datetime+".log")
}
