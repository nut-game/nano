package interfaces

// Logger interface for nano loggers
type Logger interface {
	Fatal(format ...any)
	Fatalf(format string, args ...any)
	Fatalln(args ...any)

	Debug(args ...any)
	Debugf(format string, args ...any)
	Debugln(args ...any)

	Error(args ...any)
	Errorf(format string, args ...any)
	Errorln(args ...any)

	Info(args ...any)
	Infof(format string, args ...any)
	Infoln(args ...any)

	Warn(args ...any)
	Warnf(format string, args ...any)
	Warnln(args ...any)
}
