package redsync

import (
	"testing"
)

var configs =[]*Config{
	&Config{
		Address: "127.0.0.1:6379",
	},
	&Config{
		Address: "127.0.0.1:6380",
	},
	&Config{
		Address: "127.0.0.1:6381",
	},
	&Config{
		Address: "127.0.0.1:6382",
	},
}

type Config struct {
	Address string
}

func TestRedsync(t *testing.T) {
	pools := newMockPools(configs)

	rs := New(pools)

	mutex := rs.NewMutex("test-redsync")
	err := mutex.Lock()
	if err != nil {

	}
	assertAcquired(t, pools, mutex)
}
