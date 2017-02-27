package log

import (
	"bytes"
	"fmt"
)

type TextEncoder struct {
	TypeEncoder
	Separator string
}

func NewTextEncoder(timeformat, separator string) Encoder {
	if timeformat == "" {
		timeformat = logDatetimeFmt
	}
	if separator == "" {
		separator = "="
	}
	return &TextEncoder{
		TypeEncoder: TypeEncoder{
			Timeformat: timeformat,
			Default: func(buf *bytes.Buffer, v interface{}) {
				fmt.Fprint(buf, v)
			},
			String: func(buf *bytes.Buffer, s string) {
				buf.WriteString(s)
			},
		},
		Separator: separator,
	}
}
func (t *TextEncoder) encodeKeyValue(buf *bytes.Buffer, key string, val interface{}) {
	buf.WriteString(key)
	buf.WriteString(t.Separator)
	t.EncodeVal(buf, val)
}

func (t *TextEncoder) Encode(buf *bytes.Buffer, log *Log) error {
	t.encodeKeyValue(buf, "level", log.Level.String())
	buf.WriteByte(' ')
	t.encodeKeyValue(buf, "time", log.Time)
	buf.WriteByte(' ')
	t.encodeKeyValue(buf, "pos", log.Position)

	var msg string
	if log.Format != "" {
		msg = fmt.Sprintf(log.Prefix+log.Format, log.Args...)
	} else {
		msg = log.Prefix + fmt.Sprintln(log.Args...)
		if l := len(msg); l > 0 && msg[l-1] == '\n' {
			msg = msg[:l-1]
		}
	}
	buf.WriteByte(' ')
	t.encodeKeyValue(buf, "msg", msg)

	for i := range log.Fields {
		buf.WriteByte(' ')
		t.encodeKeyValue(buf, log.Fields[i].Key, log.Fields[i].Value)
	}
	buf.WriteByte('\n')
	return nil
}
