package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

// LEVEL 日志级别，后面应写进配置文件
const LEVEL = logrus.DebugLevel

var logger *MyLogger

func Init() {
	logger = NewMyLogger()
}

func InfoC(ctx context.Context, format string, args ...any) {
	logger.InfoC(ctx, format, args...)
}

func ErrorC(ctx context.Context, format string, args ...any) {
	logger.ErrorC(ctx, format, args...)
}

func WarnC(ctx context.Context, format string, args ...any) {
	logger.WarnC(ctx, format, args...)
}

type Logger interface {
	Info(args ...any)
	InfoC(ctx context.Context, format string, args ...any)
	InfoF(format string, args ...any)
	Error(args ...any)
	ErrorC(ctx context.Context, format string, args ...any)
	ErrorF(format string, args ...any)
	Warn(args ...any)
	WarnF(format string, args ...any)
	WarnC(ctx context.Context, format string, args ...any)
}

type MyLogger struct {
	entry *logrus.Entry
}

func (l *MyLogger) InfoC(ctx context.Context, format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *MyLogger) ErrorC(ctx context.Context, format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *MyLogger) WarnC(ctx context.Context, format string, args ...any) {
	l.addCaller()
	//TODO implement me
	panic("implement me")
}

func (l *MyLogger) addCaller() *logrus.Entry {
	_, file, line, _ := runtime.Caller(2)
	return l.entry.WithField("file", fmt.Sprintf("%v:%v\n", file, line))
}

func (l *MyLogger) Info(args ...any) {
	l.addCaller().Info(args)
}

func (l *MyLogger) InfoF(format string, args ...any) {
	l.addCaller().Infof(format, args)
}

func (l *MyLogger) Error(args ...any) {
	l.addCaller().Error(args)
}

func (l *MyLogger) ErrorF(format string, args ...any) {
	l.addCaller().Errorf(format, args)
}

func (l *MyLogger) Warn(args ...any) {
	l.addCaller().Warn(args)
}

func (l *MyLogger) WarnF(format string, args ...any) {
	l.addCaller().Warn(format, args)
}

func NewMyLogger() *MyLogger {
	logger := logrus.New()
	// JSON 格式
	logger.SetFormatter(&logrus.JSONFormatter{})
	// 设置日志级别
	logger.SetLevel(LEVEL)

	// 创建第一个日志文件
	outputFile, err := os.OpenFile("/Users/zengyehong/GolandProjects/reply_mvp/service_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal(err)
	}

	//multiWriter := io.MultiWriter(os.Stdout, rabbitMQWriter)
	// 输出到标准输出流和消息队列
	multiWriter := io.MultiWriter(os.Stdout, outputFile)
	logger.SetOutput(multiWriter)
	entry := logger.WithField("category", "my log")
	return &MyLogger{entry: entry}
}
