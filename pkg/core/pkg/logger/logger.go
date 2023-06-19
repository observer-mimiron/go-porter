package logger

import (
	"fmt"
	"go-porter/pkg/core/pkg/conf/env"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// DefaultLevel the default log level
	DefaultLevel = zapcore.InfoLevel

	// DefaultTimeLayout the default time layout;
	DefaultTimeLayout = time.RFC3339

	levelInfo  = "info"
	levelDebug = "debug"
	levelWarn  = "warn"
	levelError = "error"
	levelFatal = "fatal"

	DefaultMaxSize    = 128
	DefaultMaxAge     = 30
	DefaultMaxBackUps = 300
	DefaultLocalTime  = true
	DefaultCompress   = false
)

type option struct {
	level      zapcore.Level
	fields     map[string]string
	file       io.Writer
	timeLayout string
}

func setupWithLevel(level string) zapcore.Level {
	switch level {
	case levelDebug:
		return zapcore.DebugLevel
	case levelInfo:
		return zapcore.InfoLevel
	case levelWarn:
		return zapcore.WarnLevel
	case levelError:
		return zapcore.ErrorLevel
	case levelFatal:
		return zapcore.FatalLevel
	default:
		panic("invalid log level")
	}
}

func setupWithFiles(logConf LogConf) *lumberjack.Logger {
	fileName := logConf.Filename
	if fileName == "" {
		fileName = fmt.Sprintf("%s/%s.log", env.Active(), "poter")
	}
	dir := filepath.Dir(fileName)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	return &lumberjack.Logger{
		Filename:   fileName,           // 文件路径
		MaxSize:    logConf.MaxSize,    // 单个文件最大尺寸，默认单位 M
		MaxAge:     logConf.MaxAge,     // 最多保留 300 个备份
		MaxBackups: logConf.MaxBackups, // 最大时间，默认单位 day
		LocalTime:  logConf.LocalTime,  // 使用本地时间
		Compress:   logConf.Compress,   // 是否压缩 disabled by default
	}
}

// NewJSONLogger return a json-encoder zap logger,
func NewJSONLogger(logConf LogConf) (*zap.Logger, error) {
	//初始化配置
	opt := &option{level: DefaultLevel, fields: make(map[string]string), timeLayout: DefaultTimeLayout}

	timeLayout := DefaultTimeLayout
	if logConf.TimeFormat != "" {
		timeLayout = opt.timeLayout
	}

	// similar to zap.NewProduction EncoderConfig()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger", // used by logger.Named(key); optional; useless
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(timeLayout))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// set log level 设置日志级别
	opt.level = setupWithLevel(logConf.Level)
	// lowPriority usd by info\debug\warn
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl < zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl >= zapcore.ErrorLevel
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	core := zapcore.NewTee()

	switch logConf.Model {
	case "console":
		core = zapcore.NewTee(
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
			zapcore.NewCore(jsonEncoder,
				zapcore.NewMultiWriteSyncer(stderr),
				highPriority,
			),
		)
	case "file":
		opt.file = setupWithFiles(logConf)
		if opt.file == nil {
			panic("log file not support")
		}
		core = zapcore.NewTee(core,
			zapcore.NewCore(jsonEncoder,
				zapcore.AddSync(opt.file),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= opt.level
				}),
			),
		)
	default:
		panic("log model not support")
	}

	logger := zap.New(core,
		zap.AddCaller(),
		zap.ErrorOutput(stderr),
	)

	// set fields 添加一些字段到日志中
	opt.fields["domain"] = fmt.Sprintf("%s[%s]", "poter", env.Active().Value())
	for key, value := range opt.fields {
		logger = logger.WithOptions(zap.Fields(zapcore.Field{Key: key, Type: zapcore.StringType, String: value}))
	}
	return logger, nil
}

var _ Meta = (*meta)(nil)

// Meta key-value
type Meta interface {
	Key() string
	Value() interface{}
	meta()
}

type meta struct {
	key   string
	value interface{}
}

func (m *meta) Key() string {
	return m.key
}

func (m *meta) Value() interface{} {
	return m.value
}

func (m *meta) meta() {}

// NewMeta create meat
func NewMeta(key string, value interface{}) Meta {
	return &meta{key: key, value: value}
}

// WrapMeta wrap meta to zap fields
func WrapMeta(err error, metas ...Meta) (fields []zap.Field) {
	capacity := len(metas) + 1 // namespace meta
	if err != nil {
		capacity++
	}

	fields = make([]zap.Field, 0, capacity)
	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	fields = append(fields, zap.Namespace("meta"))
	for _, meta := range metas {
		fields = append(fields, zap.Any(meta.Key(), meta.Value()))
	}

	return
}
