package golog

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"bitbucket.org/cosiner/goutils/stringutil"
	"bitbucket.org/cosiner/goutils/timeutil"
)

const (
	logDirPerm     = 0755
	logFileDateFmt = "20060102"
)

type FileLogOptions struct {
	// log file expire days, <0 to disable
	ExpireDays int
	// log file store directory
	LogDir string
	// log file buffer size
	Bufsize int
}

func (f *FileLogOptions) merge(o FileLogOptions) {
	if o.ExpireDays != 0 {
		f.ExpireDays = o.ExpireDays
	}
	if o.LogDir != "" {
		f.LogDir = o.LogDir
	}
	if o.Bufsize > 0 {
		f.Bufsize = o.Bufsize
	}
}

func newDefaultFileLogOptions(options ...FileLogOptions) FileLogOptions {
	opts := FileLogOptions{
		ExpireDays: 14,
		LogDir:     "logs",
		Bufsize:    40960,
	}
	for _, o := range options {
		opts.merge(o)
	}
	return opts
}

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
		bf.Writer.Flush()
		bf.file.Close()
	}
}

type singleFileWriter struct {
	opts FileLogOptions

	mu       sync.Mutex
	day      int
	filedate string
	file     buffedFile
}

func SingleFile(options ...FileLogOptions) (Writer, error) {
	opts := newDefaultFileLogOptions(options...)
	err := os.MkdirAll(opts.LogDir, logDirPerm)
	if err != nil {
		return nil, err
	}

	w := &singleFileWriter{
		opts: opts,
		day:  -1,
	}
	w.checkDaily()

	return w, nil
}

func (w *singleFileWriter) Write(level Level, bytes []byte) (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.checkDaily()
	_, err = w.file.Write(bytes)
	return
}

func (w *singleFileWriter) checkDaily() {
	now := time.Now()
	if d := uniqueLogDay(now); d != w.day {
		w.day = d
		date := now.Format(logFileDateFmt)
		err := w.file.init(w.logfileName(date), w.opts.Bufsize)
		if err == nil {
			w.filedate = date
		}

		cleanLogFiles(&w.opts, w.filedate, w.parseLogDate)
	}
}

func (w *singleFileWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.file.Flush()
}

func (w *singleFileWriter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.file.Close()
}

func (w *singleFileWriter) logfileName(datetime string) string {
	return filepath.Join(w.opts.LogDir, datetime+".log")
}

func (w *singleFileWriter) parseLogDate(filename string) (string, bool) {
	filename = filepath.Base(filename)
	secs := stringutil.SplitNonEmpty(filename, ".")
	if len(secs) != 2 || secs[1] != "log" || secs[0] == "" {
		return "", false
	}
	return secs[0], true
}

func uniqueLogDay(t time.Time) int {
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}

func cleanLogFiles(opts *FileLogOptions, nowdate string, parseLogDate func(string) (string, bool)) {
	if opts.ExpireDays <= 0 {
		return
	}

	items, err := ioutil.ReadDir(opts.LogDir)
	if err != nil {
		return
	}
	nowT, err := timeutil.ParseLocal(logFileDateFmt, nowdate)
	if err != nil {
		return
	}
	expiredate := nowT.Add(-timeutil.Day * time.Duration(opts.ExpireDays)).Format(logFileDateFmt)
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		date, ok := parseLogDate(item.Name())
		if !ok {
			continue
		}
		if date < expiredate {
			os.Remove(filepath.Join(opts.LogDir, item.Name()))
		}
	}
}
