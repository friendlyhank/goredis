package redsync

import (
	"testing"
)

var config =&Config{
	Address:"127.0.0.1:6379",
}

type Config struct {
	Address string
}

func TestRedsync(t *testing.T) {
	pools := newMockPools(1, config)
	rs := New(pools)

	mutex := rs.NewMutex("test-redsync")
	err := mutex.Lock()
	if err != nil {

	}

	assertAcquired(t, pools, mutex)
}
