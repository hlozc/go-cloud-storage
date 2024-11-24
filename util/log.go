package util

// 对 logrus 的封装，这个包就初始化一个日志对象(全局一个日志对象, logrus 自带多线程保护)，然后对外提供日志接口

import (
	"fmt"
	"go_cloud_storage/pkg/config"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	l *logrus.Logger
}

// 全局日志对象
var logger *Logger

// 简化一下哈希表的使用
type H = map[string]interface{}

// 初始化 logrus 日志对象
func newLogger() *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 设置输出位置
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	return &Logger{l: log}
}

// 往 logrus.Entry 中添加文件名以及行号信息，再将这个返回
func entryFileInfo(e *logrus.Entry, skip int) *logrus.Entry {
	if e != nil {
		// 获取更底层的堆栈，偏移量为 skip
		e.Data["file"] = fileInfo(skip)
	}
	return e
}

// 根据偏移量获取文件信息
func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 0
	} else {
		// 获取当前堆栈的路径并截取文件信息
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func Info(fields map[string]interface{}, format string, value ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	msg := fmt.Sprintf(format, value...)
	e := logger.l.WithFields(logrus.Fields(fields))
	entryFileInfo(e, 3)
	e.Infoln(msg)
}

func Debug(fields map[string]interface{}, format string, value ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	msg := fmt.Sprintf(format, value...)
	e := logger.l.WithFields(logrus.Fields(fields))
	entryFileInfo(e, 3)
	e.Debugln(msg)
}

func Warn(fields map[string]interface{}, format string, value ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	msg := fmt.Sprintf(format, value...)
	e := logger.l.WithFields(logrus.Fields(fields))
	entryFileInfo(e, 3)
	e.Warnln(msg)
}

func Error(fields map[string]interface{}, format string, value ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	msg := fmt.Sprintf(format, value...)
	e := logger.l.WithFields(logrus.Fields(fields))
	entryFileInfo(e, 3)
	e.Errorln(msg)
}

func Fatal(fields map[string]interface{}, format string, value ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	msg := fmt.Sprintf(format, value...)
	e := logger.l.WithFields(logrus.Fields(fields))
	entryFileInfo(e, 3)
	e.Fatalln(msg)
}

func init() {
	logger = newLogger()
	config.LogModuleInit("Log")
}
