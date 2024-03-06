package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var Logger *zap.Logger

var once = sync.Once{}

func init() {
	once.Do(func() {
		// TODO - add & check production/development config variable
		Logger, _ = zap.NewProduction()

		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		defaultEncoder := zapcore.NewJSONEncoder(config)
		// TODO - add to config
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./logs/log.json",
			LocalTime:  false,
			MaxSize:    1, // maximum size in megabytes of the log file before it gets rotated.
			MaxBackups: 2, // maximum number of old log files to retain.
			MaxAge:     3, // maximum number of days to retain old log files
		})

		stdOutWriter := zapcore.AddSync(os.Stdout)
		defaultLogLevel := zapcore.InfoLevel
		core := zapcore.NewTee(
			zapcore.NewCore(defaultEncoder, writer, defaultLogLevel),
			zapcore.NewCore(defaultEncoder, stdOutWriter, zap.InfoLevel),
		)
		Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	})
}
