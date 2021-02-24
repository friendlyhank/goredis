package redsync

import (
	"testing"
)

func TestRedsync(t *testing.T) {
	pools := newMockPools(configs)

	rs := New(pools)

	mutex := rs.NewMutex(1,"test-redsync")
	err := mutex.Lock()
	if err != nil {

	}
	//assertAcquired(t, pools, mutex)
}
