package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.SugaredLogger
var ZapLogger *zap.Logger
var OverrideDebug = false

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "log_level"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.TimeKey = "timestamp_app"

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig:    encoderConfig,
	}

	fmt.Println("env", os.Getenv("DC_APP_ENV"), os.Getenv("APP_ENV"), os.Getenv("app_env"))

	if os.Getenv("DC_APP_ENV") == "dev" {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.Encoding = "console"
		cfg.DisableCaller = false
	}

	checkDebugMode()

	var err error
	ZapLogger, err = cfg.Build()
	if err != nil {
		panic(err)
	}

	logger = ZapLogger.Sugar()
	logger.Sync()
}

func checkDebugMode() bool {
	return os.Getenv("DC_DEBUG") == "1"
}

type LogContext map[string]interface{}

func WithContext(ctx *gin.Context) {
	newLogger := ZapLogger.Sugar()

	rqId := ctx.GetHeader("X-Request-ID")
	consumerId := ctx.GetHeader("Consumer_Id")
	appVersion := ctx.GetHeader("App_Version")
	deviceOs := ctx.GetHeader("Device_Os")

	if rqId == "" {
		rqId = uuid2.New().String()
		ctx.Request.Header.Set("X-Request-ID", rqId)
	}

	logger = newLogger.With(
		zap.String("X-Request-ID", rqId),
		zap.String("Consumer_ID", consumerId),
		zap.String("Device_Os", deviceOs),
		zap.String("App_Version", appVersion),
	)
}

func OverrideDebugMode(status bool) {
	OverrideDebug = status
}

func Debug(msg string, ctx *LogContext) {
	if checkDebugMode() || OverrideDebug {
		logger.Debugw(msg, zap.Any("event", ctx))
		if OverrideDebug {
			OverrideDebugMode(!OverrideDebug)
		}
	}
}

func Info(msg string, ctx *LogContext) {
	logger.Infow(msg, zap.Any("event", ctx))
}

func Warn(msg string, ctx *LogContext) {
	logger.Warnw(msg, zap.Any("event", ctx))
}

func Error(msg string, err error, ctx *LogContext) {
	logger.Errorw(msg, zap.Error(err), zap.Any("event", ctx))
}

func Fatal(msg string, err error, ctx *LogContext) {
	logger.Fatalw(msg, zap.Error(err), zap.Any("event", ctx))
}
