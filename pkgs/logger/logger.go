package logger

import (
	"context"
	"os"
	"stock-management/pkgs/setting"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(setting setting.LoggerSetting) *LoggerZap {
	logLevel := setting.Level

	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoder := GetEncodeLog()
	hook := lumberjack.Logger{
		Filename:   setting.FileName,
		MaxSize:    setting.MaxSize,
		MaxBackups: setting.MaxBackups,
		MaxAge:     setting.MaxAge,
		Compress:   setting.Compress,
	}

	fileWriter := zapcore.AddSync(&hook)
	consoleWriter := zapcore.AddSync(os.Stdout)

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(consoleWriter, fileWriter), level)

	return &LoggerZap{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

func GetEncodeLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.TimeKey = "time"
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}

func (l *LoggerZap) Sync() {
	_ = l.Logger.Sync()
}

func (l *LoggerZap) AddTraceID(ctx context.Context) *zap.Logger {
	traceID := "unknown"
	if val := ctx.Value("TraceID"); val != nil {
		traceID = val.(string)
	}
	return l.Logger.With(zap.String("trace_id", traceID))
}
