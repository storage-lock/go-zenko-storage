module github.com/storage-lock/go-zenko-storage

go 1.25.0

require (
	github.com/minio/minio-go/v7 v7.2.1
	github.com/storage-lock/go-s3-storage v0.0.2
	github.com/storage-lock/go-storage v0.0.2
	github.com/storage-lock/go-storage-test-helper v0.0.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang-infrastructure/go-goroutine-id v0.0.0-20230917075314-f1af231c0b06 // indirect
	github.com/golang-infrastructure/go-if-expression v0.0.3 // indirect
	github.com/golang-infrastructure/go-iterator v0.0.0-20230524171120-56988a9b127c // indirect
	github.com/golang-infrastructure/go-pointer v0.0.5 // indirect
	github.com/golang-infrastructure/go-reflect-utils v0.0.0-20221130143747-965ef2eb09c3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/klauspost/cpuid/v2 v2.2.11 // indirect
	github.com/klauspost/crc32 v1.3.0 // indirect
	github.com/minio/crc64nvme v1.1.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/storage-lock/go-events v0.0.3 // indirect
	github.com/storage-lock/go-storage-events v0.0.5 // indirect
	github.com/storage-lock/go-storage-lock v0.0.3 // indirect
	github.com/storage-lock/go-utils v0.0.2 // indirect
	github.com/tinylib/msgp v1.6.1 // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.51.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	gopkg.in/ini.v1 v1.67.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/storage-lock/go-storage => ../go-storage

replace github.com/storage-lock/go-s3-storage => ../go-s3-storage

replace github.com/storage-lock/go-storage-lock => ../go-storage-lock
