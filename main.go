package main

import (
	"go_cloud_storage/pkg/minio"
	_ "go_cloud_storage/util"
)

func main() {
	minio.FUploadObject("testbucket", `D:\Git\InfraKit\Minio\lani.png`, "lani.png")
}
