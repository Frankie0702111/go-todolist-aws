package log

import (
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var errorLogger *zap.SugaredLogger

func init() {
	// Set log format
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "time",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			// Normal zapcore.SecondsDurationEncoder, Execution time consumed converted to floating point seconds
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// Implement two interfaces to determine the log level
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	// Get the io.Writer abstraction for info, error log files
	infoWriter := getWriter("./log/info.log")
	errorWriter := getWriter("./log/error.log")

	// Create Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)
	// zap.AddCaller() <-- Show error line and path only after add
	logger := zap.New(core, zap.AddCaller())
	errorLogger = logger.Sugar()
}

func getWriter(filename string) io.Writer {
	// demo.log is latest logs
	// Keep logs for 30 days and split log files at the end of each day
	hook, err := rotatelogs.New(
		// File name (If n < 0, there is no limit on the number of replacements.)
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d.log",
		rotatelogs.WithLinkName(filename),
		// Log retention time
		rotatelogs.WithMaxAge(time.Hour*24*30),
		// Splitting of files by time
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		Error("io.Writer found failed : " + err.Error())
	}

	return hook
}

func Debug(args ...interface{}) {
	errorLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	errorLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	errorLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	errorLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	errorLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	errorLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	errorLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	errorLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	errorLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	errorLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}
