package zenko_storage

import (
	"context"

	s3_storage "github.com/storage-lock/go-s3-storage"
	)

// ZenkoStorage Zenko对象存储实现 (S3兼容)
type ZenkoStorage struct {
	*s3_storage.S3Storage
}

func NewZenkoStorage(ctx context.Context, options *s3_storage.S3StorageOptions) (*ZenkoStorage, error) {
	s3Storage, err := s3_storage.NewS3Storage(ctx, options)
	if err != nil {
		return nil, err
	}
	return &ZenkoStorage{S3Storage: s3Storage}, nil
}

const StorageName = "zenko-storage"


func (x *ZenkoStorage) GetName() string {
	return StorageName
}
