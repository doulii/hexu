package log

import (
	"context"
	"fmt"
	"os"
)

type Sink interface {
	WriteLog(ctx context.Context, level Level, fields [][2]string)
}

type Level int

type Config struct {
	Sink   Sink
	MsgKey string
}

type logger struct {
	sink   Sink
	msgKey string
}

const (
	LevelInfo Level = iota + 1
	LevelWarn
	LevelError

	defaultMessageKey = "message"
)

func New(cfg *Config) *logger {
	var (
		sink   Sink
		msgKey string
	)
	if cfg != nil {
		sink = cfg.Sink
		msgKey = cfg.MsgKey
	}
	if sink == nil {
		sink = NewStdSink(os.Stderr)
	}
	if len(msgKey) == 0 {
		msgKey = defaultMessageKey
	}
	return &logger{
		sink:   sink,
		msgKey: msgKey,
	}
}

func (l *logger) writef(ctx context.Context, level Level, format string, args ...interface{}) {
	l.sink.WriteLog(ctx, level, [][2]string{
		{l.msgKey, fmt.Sprintf(format, args...)},
	})
}

func (l *logger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.writef(ctx, LevelInfo, format, args...)
}
func (l *logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.writef(ctx, LevelWarn, format, args...)
}
func (l *logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.writef(ctx, LevelError, format, args...)
}

var (
	l = New(nil)
)

func Init(cfg *Config) {
	l = New(cfg)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	l.Infof(ctx, format, args...)
}
func Warnf(ctx context.Context, format string, args ...interface{}) {
	l.Warnf(ctx, format, args...)
}
func Errorf(ctx context.Context, format string, args ...interface{}) {
	l.Errorf(ctx, format, args...)
}
