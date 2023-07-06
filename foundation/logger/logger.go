package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

// 全局日志记录 Log
var zapDefault, _ = zap.NewProduction()
var Log *zap.SugaredLogger = zapDefault.Sugar()
var once sync.Once

// NewLogger 构建了一个可以同时写到 stdout 和 file的Sugared Logger，同时对输出的时间格式进行了 format
func NewLogger(service string) *zap.SugaredLogger {
	once.Do(func() {
		fmt.Println("初始化")
		// 配置console输出
		// 这里我们配置了 console 日志的时间格式
		consoleconfig := zap.NewDevelopmentEncoderConfig()
		consoleconfig.TimeKey = "time"
		consoleconfig.LevelKey = "level"
		consoleconfig.MessageKey = "msg"
		consoleconfig.CallerKey = "caller"
		// 或者配置 time.RFC3339
		consoleconfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		consoleconfig.EncodeLevel = zapcore.CapitalLevelEncoder
		consoleconfig.EncodeCaller = zapcore.ShortCallerEncoder
		consoleconfig.EncodeDuration = zapcore.StringDurationEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(consoleconfig)

		// 配置文件输出
		fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:    "time",
			LevelKey:   "level",
			MessageKey: "msg",
			CallerKey:  "caller",
			// 或者配置 time.RFC3339
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			StacktraceKey:  service,
		})

		logFileHook := &lumberjack.Logger{
			Filename:   fmt.Sprintf("./logs/%s.log", service),
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		}
		errLogFileHook := &lumberjack.Logger{
			Filename:   fmt.Sprintf("./logs/%s-err.log", service),
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		}

		// 用于切割日志
		fileWriterSyncer := zapcore.AddSync(logFileHook)
		// 用于切割日志
		errFileWriterSyncer := zapcore.AddSync(errLogFileHook)
		// 创建不同优先级的日志
		// 用于输出 ERROR 及以上级别的日志
		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		// 用于输出 ERROR 级别以下的日志
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl < zapcore.ErrorLevel
		})

		core := zapcore.NewTee(
			// 这里没有选择 os.Stderr 作为 highPriority 的输出，统一输出到标准输出，避免新通道的开辟，提升性能
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), highPriority),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), lowPriority),
			zapcore.NewCore(fileEncoder, fileWriterSyncer, lowPriority),
			zapcore.NewCore(fileEncoder, errFileWriterSyncer, highPriority),
		)
		// zap.AddStacktrace(zap.NewAtomicLevelAt(zap.ErrorLevel)
		// ERROR 及以上级别输出调用栈
		// AddCaller 配置 Logger 以使用文件名注释每条消息，
		// zap 调用者的行号和函数名。 另请参阅 WithCaller。
		// WithCaller 配置 Logger 以使用文件名注释每条消息，
		// zap 的调用者的行号和函数名，或者不是，取决于
		// 启用的值。 这是 AddCaller 的通用形式。
		//logger := zap.New(core, zap.AddStacktrace(zap.NewAtomicLevelAt(zap.ErrorLevel)), zap.WithCaller(true), zap.AddCallerSkip(1))
		logger := zap.New(core, zap.AddStacktrace(zap.NewAtomicLevelAt(zap.ErrorLevel)), zap.WithCaller(true))

		// 创建一个要添加的字段
		extraFields := []zap.Field{
			zap.String("service", service),
		}

		Log = logger.With(extraFields...).Sugar()
	})
	return Log
}
