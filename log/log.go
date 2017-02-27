package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type Level uint8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal

	levelMin = LevelDebug
	levelMax = LevelFatal

	logDatetimeFmt = "20060102150405"
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelPanic:
		return "PANIC"
	case LevelFatal:
		return "FATAL"
	default:
		return "INFO"
	}
}

// ParseLevel parse level from string regardless of string case, if match nothing, LevelInfo was returned.
func ParseLevel(level string) Level {
	level = strings.ToUpper(strings.TrimSpace(level))
	for l := levelMin; l <= levelMax; l++ {
		if level == l.String() {
			return l
		}
	}

	return LevelInfo
}

type (
	Field struct {
		Key   string
		Value interface{}
	}

	Log struct {
		logger Logger

		Level    Level
		Time     time.Time
		Position string

		Prefix string
		Format string
		Args   []interface{}
		Fields []Field
	}
)

func (log *Log) WithField(key string, val interface{}) *Log {
	log.Fields = append(log.Fields, Field{
		Key:   key,
		Value: val,
	})
	return log
}

func (log *Log) WithFields(args ...interface{}) *Log {
	l := len(args)
	for i := 0; i < l; {
		arg := args[i]
		switch arg := arg.(type) {
		case string:
			if i < l-1 {
				log.WithField(arg, args[i+1])
			}
			i += 2
		case Field:
			log.Fields = append(log.Fields, arg)
			i++
		default:
			i += 2
		}
	}
	return log
}

func (l *Log) Depth(level Level, depth int, args ...interface{}) {
	l.Depthf(level, depth+1, "", args...)
}

func (l *Log) Depthf(level Level, depth int, format string, args ...interface{}) {
	l.Level = level
	l.Position = callerPos(depth + 1)
	l.Format = format
	l.Args = args

	l.logger.Write(l)
}

func (l *Log) Debugf(format string, args ...interface{}) {
	l.Depthf(LevelDebug, 1, format, args...)
}

func (l *Log) Infof(format string, args ...interface{}) {
	l.Depthf(LevelInfo, 1, format, args...)
}

func (l *Log) Warnf(format string, args ...interface{}) {
	l.Depthf(LevelWarn, 1, format, args...)
}

func (l *Log) Errorf(format string, args ...interface{}) {
	l.Depthf(LevelError, 1, format, args...)
}

func (l *Log) Panicf(format string, args ...interface{}) {
	l.Depthf(LevelPanic, 1, format, args...)
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	l.Depthf(LevelFatal, 1, format, args...)
}

func (l *Log) Debug(args ...interface{}) {
	l.Depth(LevelDebug, 1, args...)
}

func (l *Log) Info(args ...interface{}) {
	l.Depth(LevelInfo, 1, args...)
}

func (l *Log) Warn(args ...interface{}) {
	l.Depth(LevelWarn, 1, args...)
}

func (l *Log) Error(args ...interface{}) {
	l.Depth(LevelError, 1, args...)
}

func (l *Log) Panic(args ...interface{}) {
	l.Depth(LevelPanic, 1, args...)
}

func (l *Log) Fatal(args ...interface{}) {
	l.Depth(LevelFatal, 1, args...)
}

func isPathSeparator(r rune) bool {
	return r == '/' || r == os.PathSeparator
}

func callerPos(depth int) string {
	_, file, line, _ := runtime.Caller(depth + 1)
	i := strings.LastIndexFunc(file, isPathSeparator)
	if i >= 0 {
		j := strings.LastIndexFunc(file[:i], isPathSeparator)
		if j >= 0 {
			i = j
		}
		file = file[i+1:]
	}
	return fmt.Sprintf("%s:%d", file, line)
}
