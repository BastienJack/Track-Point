package zap

import (
	"commerce/pkg/viper"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	config = viper.Init("log")

	infoPath  = config.Viper.GetString("info")
	errorPath = config.Viper.GetString("error")

	// log writer config params
	maxSize    = 100
	maxBackups = 5
	maxAge     = 30
	compress   = false
)

func InitLogger() *zap.SugaredLogger {
	// log level definition
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zap.ErrorLevel && level >= zap.DebugLevel
	})

	// encoder
	encoder := getEncoder()

	// info level log
	infoSyncer := getInfoWriter()
	infoCore := zapcore.NewCore(encoder, infoSyncer, lowPriority)

	// error level log
	errorSyncer := getErrorWriter()
	errorCore := zapcore.NewCore(encoder, errorSyncer, highPriority)

	// combine log with different level
	core := zapcore.NewTee(infoCore, errorCore)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger := logger.Sugar()

	return sugarLogger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getInfoWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   infoPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func getErrorWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   errorPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	return zapcore.AddSync(lumberJackLogger)
}
