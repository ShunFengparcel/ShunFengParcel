package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// Level 日志级别
type Level int8

const (
	LevelDebug Level = iota - 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// Config 日志配置
type Config struct {
	Level      Level  // 日志级别
	Format     string // 日志格式: json 或 text
	OutputPath string // 日志输出路径
	MaxSize    int64  // 单个日志文件最大大小（MB）
	Console    bool   // 是否同时输出到控制台
}

// Logger 结构化日志器
type Logger struct {
	config Config
	file   *os.File
	helper *log.Helper
}

// logEntry 日志条目
type logEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	ServiceID string                 `json:"service_id,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
	Caller    string                 `json:"caller,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// NewLogger 创建新的日志器
func NewLogger(config Config) (*Logger, error) {
	l := &Logger{
		config: config,
	}

	// 创建日志目录
	if config.OutputPath != "" {
		dir := filepath.Dir(config.OutputPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// 打开日志文件
		file, err := os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		l.file = file
	}

	return l, nil
}

// Log 实现 log.Logger 接口
func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	entry := logEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Fields:    make(map[string]interface{}),
	}

	// 解析键值对
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key := fmt.Sprint(keyvals[i])
			value := keyvals[i+1]

			switch key {
			case "level":
				entry.Level = fmt.Sprint(value)
			case "msg":
				entry.Message = fmt.Sprint(value)
			case "service.id":
				entry.ServiceID = fmt.Sprint(value)
			case "trace.id":
				entry.TraceID = fmt.Sprint(value)
			case "span.id":
				entry.SpanID = fmt.Sprint(value)
			case "caller":
				entry.Caller = fmt.Sprint(value)
			default:
				entry.Fields[key] = value
			}
		}
	}

	// 格式化输出
	var output string
	if l.config.Format == "json" {
		data, _ := json.Marshal(entry)
		output = string(data) + "\n"
	} else {
		// 文本格式
		output = fmt.Sprintf("[%s] %s %s", entry.Timestamp, entry.Level, entry.Message)
		if entry.TraceID != "" {
			output += fmt.Sprintf(" trace_id=%s", entry.TraceID)
		}
		if len(entry.Fields) > 0 {
			fieldsJSON, _ := json.Marshal(entry.Fields)
			output += fmt.Sprintf(" fields=%s", string(fieldsJSON))
		}
		output += "\n"
	}

	// 写入文件
	if l.file != nil {
		if _, err := l.file.WriteString(output); err != nil {
			return err
		}
	}

	// 同时输出到控制台
	if l.config.Console {
		fmt.Print(output)
	}

	return nil
}

// Close 关闭日志器
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Helper 返回日志助手
func (l *Logger) Helper() *log.Helper {
	if l.helper == nil {
		l.helper = log.NewHelper(l)
	}
	return l.helper
}

// MultiWriter 多输出写入器
type MultiWriter struct {
	writers []io.Writer
}

// NewMultiWriter 创建多输出写入器
func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	return &MultiWriter{writers: writers}
}

// Write 实现 io.Writer 接口
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
	}
	return len(p), nil
}

// WithContext 从上下文中提取日志字段
func WithContext(ctx context.Context, keyvals ...interface{}) []interface{} {
	// 可以从 context 中提取 trace_id, user_id 等信息
	// 这里是示例，实际使用时需要根据你的 context 结构调整
	return keyvals
}

// RotateLogger 日志轮转器
type RotateLogger struct {
	config   Config
	file     *os.File
	size     int64
	basePath string
}

// NewRotateLogger 创建支持轮转的日志器
func NewRotateLogger(config Config) (*RotateLogger, error) {
	rl := &RotateLogger{
		config:   config,
		basePath: config.OutputPath,
	}

	if err := rl.rotate(); err != nil {
		return nil, err
	}

	return rl, nil
}

// Write 实现 io.Writer 接口
func (rl *RotateLogger) Write(p []byte) (n int, err error) {
	// 检查是否需要轮转
	if rl.config.MaxSize > 0 && rl.size+int64(len(p)) > rl.config.MaxSize*1024*1024 {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.file.Write(p)
	rl.size += int64(n)
	return
}

// rotate 执行日志轮转
func (rl *RotateLogger) rotate() error {
	// 关闭旧文件
	if rl.file != nil {
		rl.file.Close()
	}

	// 创建日志目录
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 如果文件存在，重命名为带时间戳的文件
	if _, err := os.Stat(rl.basePath); err == nil {
		timestamp := time.Now().Format("20060102-150405")
		newPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)
		os.Rename(rl.basePath, newPath)
	}

	// 创建新文件
	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.size = 0
	return nil
}

// Close 关闭日志器
func (rl *RotateLogger) Close() error {
	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}
