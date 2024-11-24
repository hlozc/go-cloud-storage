package main

import (
	"go_cloud_storage/pkg/minio"
)

func main() {
	minio.FUploadObject("testbucket", `D:\Git\InfraKit\Minio\lani.png`, "lani.png")
}
