package golog

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JSONEncoder struct {
	TypeEncoder
}

func NewJSONEncoder(timeformat string) Encoder {
	if timeformat == "" {
		timeformat = logDatetimeFmt
	}
	return &JSONEncoder{
		TypeEncoder: TypeEncoder{
			Timeformat: timeformat,
			Default: func(buf *bytes.Buffer, v interface{}) {
				bytes, _ := json.Marshal(v)
				buf.Write(bytes)
			},
			Bytes: func(buf *bytes.Buffer, b []byte) {
				buf.Write(b)
			},
		},
	}
}

func (j *JSONEncoder) encodeKeyValue(buf *bytes.Buffer, key string, val interface{}) {
	j.EncodeString(buf, key)
	buf.WriteByte(':')
	j.EncodeVal(buf, val)
}

func (j *JSONEncoder) Encode(buf *bytes.Buffer, log *Log) error {
	buf.WriteByte('{')
	j.encodeKeyValue(buf, "level", log.Level.String())
	buf.WriteByte(',')
	j.encodeKeyValue(buf, "time", log.Time)
	buf.WriteByte(',')
	j.encodeKeyValue(buf, "pos", log.Position)

	var msg string
	if log.Format != "" {
		if log.Prefix != "" {
			log.Format = log.Prefix + log.Format
		}
		msg = fmt.Sprintf(log.Format, log.Args...)
	} else {
		s := fmt.Sprintln(log.Args...)
		if l := len(s); s[l-1] == '\n' {
			s = s[:l-1]
		}
		msg = log.Prefix + s
	}
	buf.WriteByte(',')
	j.encodeKeyValue(buf, "msg", msg)

	for i := range log.Fields {
		buf.WriteString(`,`)
		j.encodeKeyValue(buf, log.Fields[i].Key, log.Fields[i].Value)
	}
	buf.WriteString("}\n")
	return nil
}
