package log

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestJSONEncoder(t *testing.T) {
	json := NewJSONEncoder("")
	text := NewTextEncoder("", "")
	var buf bytes.Buffer
	log := Log{
		Level:    LevelInfo,
		Time:     time.Now(),
		Position: callerPos(0),
		Prefix:   "[PREFIX] ",

		//Format: "%s %d",
		Args: []interface{}{"aa", 1},
	}

	log.appendFields(
		"A", 1,
		"B", true,
		"C", 1.1,
		"D", complex(1, 1),
		"E", time.Now(),
		"G", `"aaa"`,
		Field{"F", []int{1, 2, 3}},
	)

	json.Encode(&buf, &log)
	fmt.Print(buf.String())
	buf.Reset()

	text.Encode(&buf, &log)
	fmt.Print(buf.String())
}
