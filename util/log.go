package util

// 对 logrus 的封装，这个包就初始化一个日志对象(全局一个日志对象, logrus 自带多线程保护)，然后对外提供日志接口

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	l *logrus.Logger
}

// 全局日志对象
var logger *Logger

type H = map[string]interface{}

// 由于封装了 logrus，所以行号以及文件名需要进行修改，改成调用该包 Info 接口的位置；而不是调用 logrus.Infoln 接口的位置
// 自定义一个 hook 结构体，并且实现 hook interface 所需要的接口
type callerHook struct{}

// 在函数调用栈中，确定调用 当前文件中的 Info 是哪个堆栈
func getLogCaller() (uintptr, string, int, bool) {
	buf := make([]byte, 1024) // 创建缓冲区存储堆栈信息
	n := runtime.Stack(buf, false)
	fmt.Printf("Current stack trace:\n%s", buf[:n])

	for skip := 2; ; skip++ {
		_, file, _, ok := runtime.Caller(skip)
		if !ok { // 不存在这个堆栈，那么直接退出
			break
		}

		// 判断当前堆栈是否属于当前这个包文件
		if strings.Contains(file, "util/log.go") {
			return runtime.Caller(skip - 1)
		}
	}

	return 0, "", 0, false
}

func (h *callerHook) Fire(e *logrus.Entry) error {
	// 获取第八层堆栈
	pc, file, line, ok := getLogCaller()
	if !ok {
		return nil
	}

	// filepath.Base 就是获取这个路径下的文件名；filepath.Ext 就是只获取后缀名(包括 .)
	filename := filepath.Base(file)                   // 先获取文件名
	fn := filepath.Base(runtime.FuncForPC(pc).Name()) // 获取函数名

	// 提前增加键值对
	e.Data["file"] = fmt.Sprintf("%s:%d", filename, line)
	e.Data["func"] = fn
	return nil
}

func (h *callerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// 初始化 logrus 日志对象
func newLogger() *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 设置输出位置
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	log.AddHook(&callerHook{})

	return &Logger{l: log}
}

func Info(fields map[string]interface{}, format string, value ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	msg := fmt.Sprintf(format, value...)
	logger.l.WithFields(logrus.Fields(fields)).Infoln(msg)
}

func Debug(fields map[string]interface{}, format string, value ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	msg := fmt.Sprintf(format, value...)
	logger.l.WithFields(logrus.Fields(fields)).Debugln(msg)
}

func Warn(fields map[string]interface{}, format string, value ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	msg := fmt.Sprintf(format, value...)
	logger.l.WithFields(logrus.Fields(fields)).Warnln(msg)
}

func Error(fields map[string]interface{}, format string, value ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	msg := fmt.Sprintf(format, value...)
	logger.l.WithFields(logrus.Fields(fields)).Errorln(msg)
}

func Fatal(fields map[string]interface{}, format string, value ...interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	msg := fmt.Sprintf(format, value...)
	logger.l.WithFields(logrus.Fields(fields)).Fatalln(msg)
}

func init() {
	logger = newLogger()
}
