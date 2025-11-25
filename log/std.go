package log

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type StdSink struct{ w io.Writer }

func NewStdSink(w io.Writer) *StdSink {
	return &StdSink{w: w}
}

var (
	lvlStr = map[Level]string{
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERROR",
	}
)

func (l *StdSink) WriteLog(ctx context.Context, level Level, fields [][2]string) {
	sb := strings.Builder{}
	sb.WriteString(time.Now().Format(time.RFC3339Nano))
	ls, ok := lvlStr[level]
	if !ok {
		ls = "UNKNOWN_LEVEL_" + strconv.FormatInt(int64(level), 10)
	}
	sb.WriteString(" [")
	sb.WriteString(ls)
	sb.WriteString("]")
	for _, f := range fields {
		sb.WriteString(" ")
		sb.WriteString(f[0])
		sb.WriteString(": ")
		sb.WriteString(f[1])
	}
	fmt.Fprintln(l.w, sb.String())
}
