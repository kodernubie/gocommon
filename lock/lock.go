package lock

import (
	"time"

	"github.com/kodernubie/gocommon/cache"
	"github.com/oklog/ulid/v2"
)

func (o *AcquiredLock) Unlock() {

	val := cache.Get("lock_" + o.name)

	if val == o.val {
		cache.Del("lock_" + o.name)
	}
}

func Lock(name string, ttl ...time.Duration) *AcquiredLock {

	targetTtl := 10 * time.Second

	if len(ttl) > 0 {
		targetTtl = ttl[0]
	}

	val := ulid.Make().String()

	for {
		ok := cache.SetNX("lock_"+name, val, targetTtl)

		if ok {
			return &AcquiredLock{
				name: name,
				val:  val,
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func TryLock(name string, ttl ...time.Duration) *AcquiredLock {

	targetTtl := 10 * time.Second

	if len(ttl) > 0 {
		targetTtl = ttl[0]
	}

	val := ulid.Make().String()
	ok := cache.SetNX("lock_"+name, val, targetTtl)

	if ok {
		return &AcquiredLock{
			name: name,
			val:  val,
		}
	}

	return nil
}
