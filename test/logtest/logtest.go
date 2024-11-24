package logtest

import (
	"go_cloud_storage/util"
)

func Test() {
	util.Info(nil, "hello world")
	util.Error(util.H{"Error": "open file fail", "server": "MinIO"}, "open %v fail", "./README.md")
}
