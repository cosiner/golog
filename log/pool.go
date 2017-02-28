package log

import (
	"bytes"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func allocBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func freeBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

var logPool = sync.Pool{
	New: func() interface{} {
		return &Log{}
	},
}

func allocLog() *Log {
	return logPool.Get().(*Log)
}

func freeLog(log *Log) {
	log.Level = 0
	log.Time.Truncate(0)
	log.Position = ""
	log.Prefix = ""
	log.Format = ""
	log.Args = nil
	log.Fields = log.Fields[:0]
	log.logger = nil
	logPool.Put(log)
}
