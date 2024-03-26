package model

type RedisStorage struct {
	Mem map[string]*RedisBucket
}

type RedisBucket struct {
	Value    []byte
	ExpireAt int64
}
