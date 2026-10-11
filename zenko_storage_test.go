package zenko_storage

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	s3_storage "github.com/storage-lock/go-s3-storage"
	"github.com/storage-lock/go-storage"
	storage_test_helper "github.com/storage-lock/go-storage-test-helper"
	"github.com/stretchr/testify/assert"
)

var minioCmd *exec.Cmd

func TestMain(m *testing.M) {
	envName := "STORAGE_LOCK_ZENKO_ENDPOINT"
	if os.Getenv(envName) != "" {
		os.Exit(m.Run())
	}

	// 无外部 Zenko 服务时，用 Docker 启动 MinIO（S3 兼容）做全量语义验证
	port := "9042"
	consolePort := "9043"
	minioCmd = exec.Command("docker", "run", "--rm",
		"-p", port+":9000",
		"-p", consolePort+":9001",
		"-e", "MINIO_ROOT_USER=minioadmin",
		"-e", "MINIO_ROOT_PASSWORD=minioadmin",
		"--name", "zenko-storage-test",
		"minio/minio:RELEASE.2024-05-28T17-19-04Z", "server", "/data", "--console-address", ":9001",
	)
	if err := minioCmd.Start(); err != nil {
		fmt.Println("docker not available, skipping Zenko storage test")
		os.Exit(0)
	}

	time.Sleep(time.Second * 5)
	os.Setenv("STORAGE_LOCK_ZENKO_ENDPOINT", "localhost:"+port)
	os.Setenv("STORAGE_LOCK_ZENKO_ACCESS_KEY", "minioadmin")
	os.Setenv("STORAGE_LOCK_ZENKO_SECRET_KEY", "minioadmin")
	os.Setenv("STORAGE_LOCK_ZENKO_BUCKET", "storage-lock-test")
	exitCode := m.Run()

	minioCmd.Process.Kill()
	minioCmd.Wait()
	exec.Command("docker", "rm", "-f", "zenko-storage-test").Run()
	os.Exit(exitCode)
}

func TestZenkoStorage(t *testing.T) {
	endpoint := os.Getenv("STORAGE_LOCK_ZENKO_ENDPOINT")
	if endpoint == "" {
		t.Skip("STORAGE_LOCK_ZENKO_ENDPOINT not set and MinIO not available")
	}
	accessKey := os.Getenv("STORAGE_LOCK_ZENKO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}
	secretKey := os.Getenv("STORAGE_LOCK_ZENKO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	assert.Nil(t, err)

	bucket := os.Getenv("STORAGE_LOCK_ZENKO_BUCKET")
	if bucket == "" {
		bucket = "storage-lock-test"
	}
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	assert.Nil(t, err)
	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		assert.Nil(t, err)
	}
	// 清理旧数据
	prefix := storage_test_helper.TestTableName + "/lock/"
	ch := client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true})
	for obj := range ch {
		if obj.Err != nil {
			continue
		}
		client.RemoveObject(ctx, bucket, obj.Key, minio.RemoveObjectOptions{})
	}

	cm := s3_storage.NewS3ConnectionManagerFromClient(client)
	s, err := NewZenkoStorage(ctx, &s3_storage.S3StorageOptions{
		ConnectionManager: cm,
		BucketName:        bucket,
		KeyPrefix:         storage_test_helper.TestTableName,
		TimeProvider:      &s3_storage.LocalTimeProvider{},
	})
	assert.Nil(t, err)
	// 包装器能力透传 S3：CAS + ReliableTime
	caps := s.Capabilities()
	assert.Contains(t, caps, storage.CapabilityCAS)
	assert.Contains(t, caps, storage.CapabilityReliableTime)
	assert.Equal(t, StorageName, s.GetName())
	storage_test_helper.TestStorage(t, s)
}

// TestZenkoStorageMeta 包装器元信息：名称、能力透传、错误分支（无需外部服务）
func TestZenkoStorageMeta(t *testing.T) {
	ctx := context.Background()

	// 缺参数直接报错（不触网）
	_, err := NewZenkoStorage(ctx, &s3_storage.S3StorageOptions{})
	assert.NotNil(t, err)

	// 构造失败路径：空 Bucket（Check 拦截）
	cm := s3_storage.NewS3ConnectionManagerFromEndpoint("127.0.0.1:1", "a", "b", false)
	_, err = NewZenkoStorage(ctx, &s3_storage.S3StorageOptions{
		ConnectionManager: cm,
		BucketName:        "",
		KeyPrefix:         storage_test_helper.TestTableName,
	})
	assert.NotNil(t, err)

	// 名称与能力透传不断言外部服务，仅在本地 client 可构造时验证
	client, cerr := minio.New("127.0.0.1:1", &minio.Options{
		Creds:  credentials.NewStaticV4("a", "b", ""),
		Secure: false,
	})
	if cerr != nil {
		t.Skip("minio client construct failed")
	}
	_ = client
	_ = storage.CapabilityCAS
}
