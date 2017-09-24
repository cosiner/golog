package golog

import (
	"testing"
	"time"
)

//
//func TestConsoleLog(t *testing.T) {
//	DefaultLogger = New(LevelDebug, 0, 0, NewTextEncoder("", "")).AddWriter(Console())
//	defer DefaultLogger.Close()
//	Debug("aaa1")
//	Info("aaa2")
//	Warn("aaa4")
//	Error("aaa3")
//}

func TestFileLog(t *testing.T) {
	logger := New(LevelDebug, 0, 0, NewJSONEncoder(""))
	defer logger.Close()
	logger.AddWriter(Console())
	fw, err := Multifile(logger.Level(), "logs", 1024*10)
	if err != nil {
		t.Fatal(err)
	}
	logger.AddWriter(fw)
	logger.Warn("DDDDDDDDDDDDDDDD")
	logger.Info("DDDDDDDDDDDDDDDD")
	logger.Debug("DDDDDDDDDDDDDDDD")
	logger.Fatal("DDDDDDDDDDDDDDDD")
	time.Sleep(100 * time.Millisecond)
}
