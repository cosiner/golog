package log

import (
	"os"
	"sync/atomic"
	"time"
)

type (
	Logger interface {
		AddWriter(Writer) Logger
		Level() Level
		Flush()
		Close()
		Prefix(string) Logger

		WithField(key string, val interface{}) *Log
		WithFields(...interface{}) *Log

		Debug(...interface{})
		Info(...interface{})
		Warn(...interface{})
		Error(...interface{})
		Panic(...interface{})
		Fatal(...interface{}) // exit process
		Debugf(string, ...interface{})
		Infof(string, ...interface{})
		Warnf(string, ...interface{})
		Errorf(string, ...interface{})
		Panicf(string, ...interface{})
		Fatalf(string, ...interface{})
		Depth(Level, int, ...interface{})
		Depthf(Level, int, string, ...interface{})
		Write(*Log)
	}

	Writer interface {
		Write(Level, []byte) error
		Close()
		Flush()
	}

	logger struct {
		encoder Encoder
		level   Level
		writers []Writer
		prefix  string

		flushInterval time.Duration
		flush         chan struct{}

		closeFlag int32
	}
)

func New(level Level, flushSeconds, backlog int, encoder Encoder) Logger {
	if flushSeconds == 0 {
		flushSeconds = 30
	}
	if backlog == 0 {
		backlog = 100
	}

	l := &logger{
		level:   level,
		encoder: encoder,
	}
	if level <= levelMax {
		l.flush = make(chan struct{}, 1)
		l.flushInterval = time.Duration(flushSeconds) * time.Second
	}
	l.start()
	return l
}

func (l *logger) isClosed() bool {
	return atomic.LoadInt32(&l.closeFlag) == 1
}

func (l *logger) markClosed() bool {
	return atomic.CompareAndSwapInt32(&l.closeFlag, 0, 1)
}

func (l *logger) AddWriter(w Writer) Logger {
	l.writers = append(l.writers, w)
	return l
}

func (l *logger) Prefix(p string) Logger {
	if p == l.prefix {
		return l
	}

	nl := *l
	nl.prefix = p
	return &nl
}

func (l *logger) Level() Level {
	return l.level
}

func (l *logger) start() {
	go func(l *logger) {
		ticker := time.Tick(l.flushInterval)
		for {
			select {
			case <-ticker:
			case _, ok := <-l.flush:
				if !ok {
					return
				}
			}

			for _, writer := range l.writers {
				writer.Flush()
			}
		}
	}(l)
}

func (l *logger) Write(log *Log) {
	level := log.Level
	if level < l.level || l.isClosed() {
		freeLog(log)
		return
	}

	buf := allocBuffer()
	l.encoder.Encode(buf, log)
	bytes := buf.Bytes()
	for _, writer := range l.writers {
		writer.Write(level, bytes)
	}
	freeBuffer(buf)
	freeLog(log)

	if level == LevelPanic {
		panic(string(bytes))
	}
	if level == LevelFatal {
		l.Close()
		os.Exit(-1)
	}
}

func (l *logger) Flush() {
	if !l.isClosed() && l.flush != nil {
		select {
		case l.flush <- struct{}{}:
		default:
		}
	}
}

func (l *logger) Close() {
	if !l.markClosed() {
		return
	}
	for _, w := range l.writers {
		w.Close()
	}
	close(l.flush)
}

func (l *logger) newLog(prefix string) *Log {
	log := allocLog()
	log.Time = time.Now()
	log.Prefix = prefix
	log.logger = l
	return log
}

func (l *logger) WithField(key string, val interface{}) *Log {
	return l.newLog(l.prefix).appendField(key, val)
}

func (l *logger) WithFields(args ...interface{}) *Log {
	return l.newLog(l.prefix).appendFields(args...)
}

func (l *logger) Depth(level Level, depth int, args ...interface{}) {
	l.Depthf(level, depth+1, "", args...)
}

func (l *logger) Depthf(level Level, depth int, format string, args ...interface{}) {
	if level >= l.level {
		l.newLog(l.prefix).Depthf(level, depth+1, format, args...)
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Depthf(LevelDebug, 1, format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Depthf(LevelInfo, 1, format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Depthf(LevelWarn, 1, format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Depthf(LevelError, 1, format, args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.Depthf(LevelPanic, 1, format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Depthf(LevelFatal, 1, format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.Depth(LevelDebug, 1, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.Depth(LevelInfo, 1, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.Depth(LevelWarn, 1, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.Depth(LevelError, 1, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.Depth(LevelPanic, 1, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.Depth(LevelFatal, 1, args...)
}

var DefaultLogger = New(LevelDebug, 0, 0, NewTextEncoder("", "")).AddWriter(Console())

func Depth(level Level, depth int, args ...interface{}) {
	DefaultLogger.Depth(level, depth+1, args...)
}

func Depthf(level Level, depth int, format string, args ...interface{}) {
	DefaultLogger.Depthf(level, depth+1, format, args...)
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Depthf(LevelDebug, 1, format, args...)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Depthf(LevelInfo, 1, format, args...)
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger.Depthf(LevelWarn, 1, format, args...)
}

func Errorf(format string, args ...interface{}) {
	DefaultLogger.Depthf(LevelError, 1, format, args...)
}

func Panicf(format string, args ...interface{}) {
	DefaultLogger.Depthf(LevelPanic, 1, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Depthf(LevelFatal, 1, format, args...)
}

func Debug(args ...interface{}) {
	DefaultLogger.Depth(LevelDebug, 1, args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Depth(LevelInfo, 1, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Depth(LevelWarn, 1, args...)
}

func Error(args ...interface{}) {
	DefaultLogger.Depth(LevelError, 1, args...)
}

func Panic(args ...interface{}) {
	DefaultLogger.Depth(LevelPanic, 1, args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Depth(LevelFatal, 1, args...)
}

func WithField(key string, val interface{}) *Log {
	return DefaultLogger.WithField(key, val)
}

func WithFields(args ...interface{}) *Log {
	return DefaultLogger.WithFields(args...)
}
