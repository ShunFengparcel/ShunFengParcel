package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Logger 全局zap日志实例
	Logger *zap.Logger
	// SugaredLogger 全局zap sugar日志实例（更便捷的API）
	SugaredLogger *zap.SugaredLogger
)

// InitZap 初始化Zap日志器
func InitZap() {
	// 创建日志目录
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}

	// 配置日志编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     customTimeEncoder,              // 自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行时间格式
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
	}

	// 配置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel) // 默认Info级别，可改为DebugLevel

	// 配置日志输出：控制台 + 文件
	core := zapcore.NewTee(
		// 控制台输出（彩色）
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			atomicLevel,
		),
		// Info及以上级别文件输出
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(getLogWriter(filepath.Join(logDir, "app.log"))),
			atomicLevel,
		),
		// Error级别单独文件输出
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(getLogWriter(filepath.Join(logDir, "error.log"))),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			}),
		),
	)

	// 创建Logger实例
	Logger = zap.New(core,
		zap.AddCaller(),                   // 添加调用者信息
		zap.AddCallerSkip(1),              // 跳过1层调用栈
		zap.AddStacktrace(zap.ErrorLevel), // Error级别添加堆栈跟踪
	)

	// 创建SugaredLogger
	SugaredLogger = Logger.Sugar()

	// 将Zap适配为Kratos的log.Logger接口
	kratosLogger := NewKratosLogger(Logger)
	log.SetLogger(kratosLogger)

	fmt.Println("✅ Zap日志系统初始化成功")
}

// getLogWriter 获取日志写入器（支持日志轮转）
func getLogWriter(filename string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    100,      // 单个文件最大尺寸，单位MB
		MaxBackups: 30,       // 最多保留30个备份
		MaxAge:     7,        // 最多保留7天
		Compress:   true,     // 是否压缩旧日志
		LocalTime:  true,     // 使用本地时间
	}
	return zapcore.AddSync(lumberJackLogger)
}

// customTimeEncoder 自定义时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// NewKratosLogger 创建Kratos日志适配器
func NewKratosLogger(zapLogger *zap.Logger) log.Logger {
	return NewZapLogger(zapLogger)
}

// --- 便捷的日志记录函数 ---

// Info 记录Info级别日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Infof 记录Info级别日志（格式化）
func Infof(template string, args ...interface{}) {
	SugaredLogger.Infof(template, args...)
}

// Debug 记录Debug级别日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Debugf 记录Debug级别日志（格式化）
func Debugf(template string, args ...interface{}) {
	SugaredLogger.Debugf(template, args...)
}

// Warn 记录Warn级别日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Warnf 记录Warn级别日志（格式化）
func Warnf(template string, args ...interface{}) {
	SugaredLogger.Warnf(template, args...)
}

// Error 记录Error级别日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Errorf 记录Error级别日志（格式化）
func Errorf(template string, args ...interface{}) {
	SugaredLogger.Errorf(template, args...)
}

// --- Kratos日志适配器 ---

// NewZapLogger 创建Zap日志适配器（兼容Kratos）
func NewZapLogger(zapLogger *zap.Logger) log.Logger {
	return kratoszap.NewLogger(zapLogger)
}

// --- 结构化日志辅助函数 ---

// LogPaymentError 记录支付相关错误
func LogPaymentError(orderNo, errMsg string, err error) {
	Error("支付处理错误",
		zap.String("order_no", orderNo),
		zap.String("error_msg", errMsg),
		zap.Error(err),
	)
}

// LogOrderError 记录订单相关错误
func LogOrderError(orderNo, operation string, err error) {
	Error("订单操作错误",
		zap.String("order_no", orderNo),
		zap.String("operation", operation),
		zap.Error(err),
	)
}

// LogDatabaseError 记录数据库错误
func LogDatabaseError(operation, table string, err error) {
	Error("数据库操作错误",
		zap.String("operation", operation),
		zap.String("table", table),
		zap.Error(err),
	)
}

// LogAPIError 记录API调用错误
func LogAPIError(api, method string, statusCode int, err error) {
	Error("API调用错误",
		zap.String("api", api),
		zap.String("method", method),
		zap.Int("status_code", statusCode),
		zap.Error(err),
	)
}

// LogAlipayError 记录支付宝相关错误
func LogAlipayError(operation, tradeNo string, err error, params map[string]interface{}) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("trade_no", tradeNo),
		zap.Error(err),
	}

	// 添加额外参数
	for k, v := range params {
		fields = append(fields, zap.Any(k, v))
	}

	Error("支付宝操作错误", fields...)
}

// LogBusinessInfo 记录业务信息日志
func LogBusinessInfo(module, action string, data map[string]interface{}) {
	fields := []zap.Field{
		zap.String("module", module),
		zap.String("action", action),
	}

	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}

	Info("业务操作", fields...)
}

// LogRequest 记录HTTP请求日志
func LogRequest(method, path, clientIP string, latency time.Duration, statusCode int) {
	Info("HTTP请求",
		zap.String("method", method),
		zap.String("path", path),
		zap.String("client_ip", clientIP),
		zap.Duration("latency", latency),
		zap.Int("status_code", statusCode),
	)
}
