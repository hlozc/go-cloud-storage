package minio

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	. "go_cloud_storage/pkg/config"
	"go_cloud_storage/util"
)

var minioClient *minio.Client
var location string = "" // 暂时用不到

var mimeTypeMap = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".bmp":  "image/bmp",
	".tiff": "image/tiff",
	".mp4":  "video/mp4",
	".avi":  "video/x-msvideo",
	".mkv":  "video/x-matroska",
	".mov":  "video/quicktime",
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".pdf":  "application/pdf",
	".txt":  "text/plain",
	".html": "text/html",
	".css":  "text/css",
	".js":   "application/javascript",
	".json": "application/json",
	".zip":  "application/zip",
	".rar":  "application/x-rar-compressed",
	".7z":   "application/x-7z-compressed",
}

func NewMinioClient() {
	useSSL := false

	var err error
	minioClient, err = minio.New((Cfg.MinIOHost), &minio.Options{
		Creds:  credentials.NewStaticV4(Cfg.MinIOAccessKey, Cfg.MinIOSecretKey, ""),
		Secure: useSSL,
	})

	if err != nil || minioClient == nil {
		util.Error(util.H{"Error": err}, "MinIO client connect fail")
		os.Exit(0)
	}
}

// 判断存储桶是否存在
func BucketExists(bucketName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err == nil {
		return exists, nil
	} else {
		util.Error(util.H{"Error": err}, "[BucketName: %v] not exists", bucketName)
		return exists, err
	}
}

func MakeBucket(bucketName string) error {
	if exists, err := BucketExists(bucketName); exists {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		util.Error(util.H{"Error": err}, "[BucketName: %v] make fail", bucketName)
	}

	util.Info(nil, bucketName+" bucket make success")
	return err
}

// 根据内容来判断文件 mime 类型
func detectContentType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		util.Error(util.H{"Error": err}, "[%v] open fail", filePath)
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		util.Error(util.H{"Error": err}, "[%v] read fail", filePath)
		return "", err
	}

	return http.DetectContentType(buf), nil
}

func getContentType(filePath string) (string, error) {
	// 获取文件拓展名
	ext := filepath.Ext(filePath)
	if ext == "" {
		return detectContentType(filePath)
	}

	if mimeType, exists := mimeTypeMap[ext]; exists {
		return mimeType, nil
	}
	return "application/octet-stream", nil // default
}

// 上传对象（内存级）
func UploadOjbect(bucketName string, data []byte, n int64, objName string) bool {
	reader := bytes.NewReader(data)

	contentType, _ := getContentType(objName)
	util.Debug(nil, "[BucketName: %v | OjbectName: %v] content type is %v", bucketName, objName, contentType)
	info, err := minioClient.PutObject(context.Background(), bucketName, objName, reader, n, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		util.Error(util.H{"Error": err}, "[BucketName: %v | OjbectName: %v] upload fail", bucketName, objName)
		return false
	}
	util.Info(nil, "Successfully uploaded %s of size %d\n", objName, info.Size)
	return true
}

func FUploadObject(bucketName string, filePath string, objName string) bool {
	contentType, _ := getContentType(filePath)
	util.Debug(nil, "[BucketName: %v | OjbectName: %v] content type is %v", bucketName, objName, contentType)

	info, err := minioClient.FPutObject(context.Background(), bucketName, objName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		util.Error(util.H{"Error": err}, "[FilePath: %v -> BucketName: %v | OjbectName: %v] upload fail", filePath, bucketName, objName)
		return false
	}

	util.Info(nil, "Successfully uploaded %s of size %d\n", filePath, info.Size)
	return true
}

// 获取该用户指定对象的 obj 资源的 url (私有的也可以直接通过这个 URL 下载)
func PresignedObjectURL(bucketName, objName string, expiry time.Duration) (string, error) {
	url, err := minioClient.PresignedGetObject(context.Background(), bucketName, objName, expiry, nil)
	if err != nil {
		util.Error(util.H{"Error": err}, "[BucketName: %v | ObjectName: %v] get URL fail", bucketName, objName)
		return "", err
	}

	return url.String(), nil
}

func ObjectInfo(bucketName, objName string) *minio.ObjectInfo {
	info, err := minioClient.StatObject(context.Background(), bucketName, objName, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil
		}
		util.Error(util.H{"Error": err}, "[BucketName: %v | ObjectName: %v] stat fail", bucketName, objName)
		return nil
	}

	return &info
}

func DeleteObject(bucketName, objName string) error {
	err := minioClient.RemoveObject(context.Background(), bucketName, objName, minio.RemoveObjectOptions{})
	if err != nil {
		util.Error(util.H{"Error": err}, "[BucketName: %v | ObjectName: %v] delete fail", bucketName, objName)
		return err
	}

	return nil
}

// 通过管道的方式来删除，和 MinIO 的交互是并发的，并且逐个传递对象，内存开销不会暴增
func DeleteObjects(bucketName string, objNames []string) {
	objCh := make(chan minio.ObjectInfo)

	// 启动一个协程，将要删除的 obj 逐个放到管道里面
	go func() {
		defer close(objCh)

		for _, objName := range objNames {
			objCh <- minio.ObjectInfo{Key: objName}
		}
	}()

	for err := range minioClient.RemoveObjects(context.Background(), bucketName, objCh, minio.RemoveObjectsOptions{}) {
		if err.Err != nil {
			util.Error(util.H{"Error": err}, "[BucketName: %v | ObjectName: %v] delete fail", bucketName, err.ObjectName)
		}
	}
}

func init() {
	NewMinioClient()
	LogModuleInit("MinIO")
}
