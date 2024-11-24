package config

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

type GlobalConfig struct {
	ListenPort int

	MinIOHost      string
	MinIOSecretKey string
	MinIOAccessKey string

	RedisHost string
}

var Cfg *GlobalConfig

func fatalExit(msg string, err error) {
	log.WithFields(log.Fields{"Error": err}).Fatalln(msg)
	os.Exit(0)
}

func init() {
	// 打印启动日志
	hello()

	// ini 接口结构体，就是解析 config.ini 文件之后的句柄/对象
	var i *ini.File

	i, err := ini.Load("conf/config.ini")
	if err != nil {
		fatalExit("Parse conf.ini file fail", err)
	}

	// 开始解析配置文件里面的各个标签
	server, err := i.GetSection("server")
	if err != nil {
		fatalExit("config.ini not exists \"server\" field", err)
	}

	minio, err := i.GetSection("MinIO")
	if err != nil {
		fatalExit("config.ini not exists \"MinIO\" field", err)
	}

	redis, err := i.GetSection("Redis")
	if err != nil {
		fatalExit("config.ini not exists \"Redis\" field", err)
	}

	// 将每一个标签中的 Key 读取出来，MustInt, MustString 表示转成 int/string，如果不存在 Key，则转成某个默认值
	Cfg = &GlobalConfig{
		ListenPort: server.Key("LISTEN_PORT").MustInt(8080),

		MinIOHost:      minio.Key("HOST").MustString(""),
		MinIOAccessKey: minio.Key("ACCESS_KEY").MustString(""),
		MinIOSecretKey: minio.Key("SECRET_KEY").MustString(""),

		RedisHost: redis.Key("HOST").MustString(""),
	}
}

func hello() {
	// 定义彩色输出
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	// 图案部分
	logo := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
		yellow("               *       .--."),
		yellow("                   / /  `"),
		yellow("      +           | |"),
		yellow("             '     \\ \\__,"),
		yellow("         *          '--'  *      🚀"),
	)

	// 项目名和描述
	title := bold(cyan("🌟 Go-Storage-Cloud 🌟"))
	description := green("Secure. Fast. Reliable.")

	// 输出
	fmt.Printf("%s\n   %s\n   %s\n\n", logo, title, description)
	fmt.Println(green(`
🌟 Version: 1.0.0
🌐 Access your cloud storage at: http://localhost:8080
📦 Ready to store, retrieve, and share your files securely!
`))
}

func LogModuleInit(moduleName string) {
	// 定义彩色输出
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// 输出模块初始化日志
	fmt.Printf("[%s] %s module initialized successfully!\n",
		yellow(time.Now().Format("2006-01-02 15:04:05")),
		green(moduleName))
}
