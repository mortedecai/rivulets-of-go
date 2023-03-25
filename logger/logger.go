package logger

import "go.uber.org/zap"

// Logger represents the basic interface for logging in RoG.
//    It makes available the *w functions from zap.SugaredLogger.
type Logger interface {
	// WithName sets the name of the current logger
	WithName(string) Logger
	// Debugw creates a debug log entry with the specified method identifier
	Debugw(method string, keysAndValues ...interface{})
	// Errorw creates an error log entry with the specified method identifier
	Errorw(method string, keysAndValues ...interface{})
	// Infow creates an info log entry with the specified method identifier
	Infow(method string, keysAndValues ...interface{})
	// Panicw creates a panic log entry with the specified method identifier
	Panicw(method string, keysAndValues ...interface{})
	// Warnw creates a warn log entry with the specified method identifier
	Warnw(method string, keysAndValues ...interface{})
	// Fatalw creates a fatal log entry with the specified method identifier
	Fatalw(method string, keysAndValues ...interface{})
}

type logger struct {
	// currently the logger struct only embeds the zap.SugaredLogger
	*zap.SugaredLogger
}

// New creates a new Logger for use.
func New(name string, dev bool) (Logger, error) {
	var lgr *logger
	var l *zap.Logger
	var err error

	creator := zap.NewProduction
	if dev {
		creator = zap.NewDevelopment
	}

	if l, err = creator(); err != nil {
		return nil, err
	}

	lgr = &logger{
		SugaredLogger: l.Sugar().Named(name),
	}
	return lgr, nil
}

func (l *logger) WithName(name string) Logger {
	return &logger{SugaredLogger: l.SugaredLogger.Named(name)}
}
