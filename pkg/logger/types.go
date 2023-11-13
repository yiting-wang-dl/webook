package logger

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func example() {
	var l Logger
	l.Info("users wechat id %d", 123)
}

type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}

type Field struct {
	Key string
	Val any
}

func exampleV1() {
	var l LoggerV1
	// this is a new user union_id=123
	l.Info("This is a new user, ", Field{Key: "union_id", Val: 123})
}

type LoggerV2 interface {
	// args have to be even，like key1,value1,key2,value2
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
