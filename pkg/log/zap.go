package log

import (
	cfg "github.com/bsir2020/basework/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var logger *zap.Logger
var logfile string

type ZapLog struct {
	logger *zap.Logger
}

func init() {
	logfile = cfg.EnvConfig.Log.Logfile

	hook := lumberjack.Logger{
		Filename:   logfile, // 日志文件路径
		MaxSize:    128,     // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 300,     // 日志文件最多保存多少个备份
		MaxAge:     120,     // 文件最多保存多少天
		Compress:   true,    // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	logger = zap.New(core, caller, development, filed)
	logger.Info("log 初始化成功")
}

func New() *ZapLog {
	return &ZapLog{
		logger: logger,
	}
}

func (z *ZapLog) Info(methodName, msg string, err error) {
	z.logger.Info(msg, zap.String(methodName, err.Error()))
}

func (z *ZapLog) Error(methodName, msg string, err error) {
	z.logger.Error(msg, zap.String(methodName, err.Error()))
}

func (z *ZapLog) Debug(methodName, msg string, err error) {
	z.logger.Debug(msg, zap.String(methodName, err.Error()))
}

func (z *ZapLog) Fatal(methodName, msg string, err error) {
	z.logger.Fatal(msg, zap.String(methodName, err.Error()))
}
