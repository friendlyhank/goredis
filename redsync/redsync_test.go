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
		t.Errorf("%v",err)
		return
	}
	//assertAcquired(t, pools, mutex)
}
