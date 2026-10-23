package zenko_storage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	s3_storage "github.com/storage-lock/go-s3-storage"
	"github.com/storage-lock/go-storage"
	"github.com/stretchr/testify/assert"
)

// TestZenkoStorageMeta 包装器元信息测试：名称、能力透传、错误分支与成功构造。
// 全程离线：不连接任何外部 Zenko/S3 服务，用内存假 S3 server 模拟 BucketExists。
func TestZenkoStorageMeta(t *testing.T) {
	ctx := context.Background()

	// 1. 缺参数直接报错（S3StorageOptions.Check 拦截，不触网）
	_, err := NewZenkoStorage(ctx, &s3_storage.S3StorageOptions{})
	assert.NotNil(t, err)

	// 2. 构造失败路径：ConnectionManager 指向不可达端点，BucketExists 失败
	cmBad := s3_storage.NewS3ConnectionManagerFromEndpoint("127.0.0.1:1", "a", "b", false)
	_, err = NewZenkoStorage(ctx, &s3_storage.S3StorageOptions{
		ConnectionManager: cmBad,
		BucketName:        "storage-lock-test",
		KeyPrefix:         storage.DefaultStorageTableName,
		TimeProvider:      &s3_storage.LocalTimeProvider{},
	})
	assert.NotNil(t, err)

	// 3. 成功路径：内存假 S3 server，GET ?location= 返回合法 Location XML，HEAD 返回 200。
	//    minio-go 的 BucketExists 内部先取 bucket region（GET ?location=）再 HEAD 检查存在性。
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">us-east-1</LocationConstraint>`))
	}))
	defer ts.Close()

	endpoint := strings.TrimPrefix(ts.URL, "http://")
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("a", "b", ""),
		Secure: false,
	})
	assert.Nil(t, err)

	cmOK := s3_storage.NewS3ConnectionManagerFromClient(client)
	s, err := NewZenkoStorage(ctx, &s3_storage.S3StorageOptions{
		ConnectionManager: cmOK,
		BucketName:        "storage-lock-test",
		KeyPrefix:         storage.DefaultStorageTableName,
		TimeProvider:      &s3_storage.LocalTimeProvider{},
	})
	assert.Nil(t, err)
	assert.NotNil(t, s)

	// 4. 名称与能力透传（S3：CAS + ReliableTime + Watch，包装器不改写）
	assert.Equal(t, StorageName, s.GetName())
	assert.Contains(t, s.Capabilities(), storage.CapabilityCAS)
	assert.Contains(t, s.Capabilities(), storage.CapabilityReliableTime)
}
