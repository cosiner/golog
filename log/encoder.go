package log

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type Encoder interface {
	Encode(*bytes.Buffer, *Log) error
}

type TypeEncoder struct {
	Timeformat string
	Default    func(buf *bytes.Buffer, v interface{})
	String     func(buf *bytes.Buffer, s string)
}

func (t *TypeEncoder) EncodeVal(buf *bytes.Buffer, val interface{}) {
	switch v := val.(type) {
	case int:
		t.EncodeInt(buf, int64(v))
	case int8:
		t.EncodeInt(buf, int64(v))
	case int16:
		t.EncodeInt(buf, int64(v))
	case int32:
		t.EncodeInt(buf, int64(v))
	case int64:
		t.EncodeInt(buf, int64(v))
	case uint8:
		t.EncodeUint(buf, uint64(v))
	case uint16:
		t.EncodeUint(buf, uint64(v))
	case uint32:
		t.EncodeUint(buf, uint64(v))
	case uint64:
		t.EncodeUint(buf, uint64(v))
	case uint:
		t.EncodeUint(buf, uint64(v))
	case float32:
		t.EncodeFloat(buf, float64(v))
	case float64:
		t.EncodeFloat(buf, float64(v))
	case complex64:
		t.EncodeComplex(buf, complex128(v))
	case complex128:
		t.EncodeComplex(buf, complex128(v))
	case bool:
		t.EncodeBool(buf, v)
	case string:
		t.EncodeString(buf, v)
	case time.Time:
		t.EncodeTime(buf, v)
	case time.Duration:
		t.EncodeDuration(buf, v)
	case []int8:
		for _, v := range v {
			t.EncodeInt(buf, int64(v))
		}
	case []int16:
		for _, v := range v {
			t.EncodeInt(buf, int64(v))
		}
	case []int32:
		for _, v := range v {
			t.EncodeInt(buf, int64(v))
		}
	case []int64:
		for _, v := range v {
			t.EncodeInt(buf, int64(v))
		}
	case []int:
		for _, v := range v {
			t.EncodeInt(buf, int64(v))
		}
	case []uint8:
		for _, v := range v {
			t.EncodeUint(buf, uint64(v))
		}
	case []uint16:
		for _, v := range v {
			t.EncodeUint(buf, uint64(v))
		}
	case []uint32:
		for _, v := range v {
			t.EncodeUint(buf, uint64(v))
		}
	case []uint64:
		for _, v := range v {
			t.EncodeUint(buf, uint64(v))
		}
	case []uint:
		for _, v := range v {
			t.EncodeUint(buf, uint64(v))
		}
	case []float32:
		for _, v := range v {
			t.EncodeFloat(buf, float64(v))
		}
	case []float64:
		for _, v := range v {
			t.EncodeFloat(buf, float64(v))
		}
	case []complex64:
		for _, v := range v {
			t.EncodeComplex(buf, complex128(v))
		}
	case []complex128:
		for _, v := range v {
			t.EncodeComplex(buf, complex128(v))
		}
	case []bool:
		for _, v := range v {
			t.EncodeBool(buf, v)
		}
	case []string:
		for _, v := range v {
			t.EncodeString(buf, v)
		}
	case []time.Time:
		for _, v := range v {
			t.EncodeTime(buf, v)
		}
	case []time.Duration:
		for _, v := range v {
			t.EncodeDuration(buf, v)
		}
	default:
		t.Default(buf, v)
	}
}

func (*TypeEncoder) EncodeInt(buf *bytes.Buffer, n int64) {
	buf.WriteString(strconv.FormatInt(n, 10))
}

func (*TypeEncoder) EncodeUint(buf *bytes.Buffer, n uint64) {
	buf.WriteString(strconv.FormatUint(n, 10))
}

func (*TypeEncoder) EncodeFloat(buf *bytes.Buffer, n float64) {
	buf.WriteString(strconv.FormatFloat(n, 'f', 4, 64))
}

func (*TypeEncoder) EncodeComplex(buf *bytes.Buffer, s complex128) {
	buf.WriteString(fmt.Sprint(s))
}

func (t *TypeEncoder) EncodeTime(buf *bytes.Buffer, time time.Time) {
	t.EncodeString(buf, time.Format(t.Timeformat))
}

func (t *TypeEncoder) EncodeDuration(buf *bytes.Buffer, dur time.Duration) {
	t.EncodeString(buf, dur.String())
}

func (*TypeEncoder) EncodeBool(buf *bytes.Buffer, b bool) {
	if b {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
}

func (t *TypeEncoder) EncodeString(buf *bytes.Buffer, s string) {
	if t.String != nil {
		t.String(buf, s)
		return
	}

	buf.WriteByte('"')
	for _, b := range s {
		if b == '"' {
			buf.WriteByte('\\')
		}
		buf.WriteRune(b)
	}
	buf.WriteByte('"')
}
