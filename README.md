# go-zenko-storage

这个包是 `s3_storage.S3Storage` 的薄包装：结构体嵌入了 S3Storage，
`Capabilities()`、CAS 条件写等全部由基类提供。它自己只覆盖 `GetName()`
和构造函数。

能否建锁取决于底层网关是否真的实现了 S3 的 ETag/If-Match 条件写。
声明为 S3 兼容但忽略 `If-Match` 的网关会让 CAS 退化成普通覆盖写——
这种情况本库无法从客户端侧检测，部署前必须自行验证。

## 分布式锁的前提

锁的正确性依赖两件事：同一个 `lockId` 上有**线性一致的 CAS**，以及
**可靠、单调且跨节点可协调的时间源**。能力标志是接口自述，不是证明；
具体约束与故障模型见
[go-storage-lock 的 README](https://github.com/storage-lock/go-storage-lock#能力边界)。

## 用法

```go
store, err := zenko_storage.NewZenkoStorage(...)
lock, err := storage_lock.NewStorageLock(store, "my-lock-key")
owner := lock.NewOwnerId()
err = lock.Lock(ctx, owner)
```

同一套 `构造 Storage → NewStorageLock → Lock(ctx, ownerId)` 在
[go-storage-lock-example](https://github.com/storage-lock/go-storage-lock-example) 里有
各后端的完整示例。

## License

MIT
