package redsync

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/hashicorp/go-multierror"
)

// A DelayFunc is used to decide the amount of time to wait between retries.
type DelayFunc func(tries int) time.Duration

// A Mutex is a distributed mutual exclusion lock.
type Mutex struct {
	Id int64
	name   string
	expiry time.Duration

	tries     int
	delayFunc DelayFunc

	factor float64

	quorum int

	genValueFunc func() (string, error)
	value        string
	until        time.Time

	pools []Pool
}

// Lock locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *Mutex) Lock() error {
	//生成唯一随机串，默认用base64
	value, err := m.genValueFunc()
	if err != nil {
		return err
	}

	//失败尝试
	for i := 0; i < m.tries; i++ {
		//设置延迟系数
		if i != 0 {
			time.Sleep(m.delayFunc(i))
		}

		start := time.Now()

		n, err := m.actOnPoolsAsync(func(nodeId int,pool Pool) (bool, error) {
			var tryId = i
			return m.acquire(nodeId,tryId,pool, value)
		})
		if n == 0 && err != nil {
			return err
		}

		now := time.Now()
		until := now.Add(m.expiry - now.Sub(start) - time.Duration(int64(float64(m.expiry)*m.factor)))
		if n >= m.quorum && now.Before(until) {
			m.value = value
			m.until = until
			return nil
		}
		m.actOnPoolsAsync(func(nodeId int,pool Pool) (bool, error) {
			return m.release(pool, value)
		})
	}

	return ErrFailed
}

// Unlock unlocks m and returns the status of unlock.
func (m *Mutex) Unlock() (bool, error) {
	n, err := m.actOnPoolsAsync(func(nodeId int,pool Pool) (bool, error) {
		return m.release(pool, m.value)
	})
	if n < m.quorum {
		return false, err
	}
	return true, nil
}

// Extend resets the mutex's expiry and returns the status of expiry extension.
func (m *Mutex) Extend() (bool, error) {
	n, err := m.actOnPoolsAsync(func(nodeId int,pool Pool) (bool, error) {
		return m.touch(pool, m.value, int(m.expiry/time.Millisecond))
	})
	if n < m.quorum {
		return false, err
	}
	return true, nil
}

func (m *Mutex) Valid() (bool, error) {
	n, err := m.actOnPoolsAsync(func(nodeId int,pool Pool) (bool, error) {
		return m.valid(pool)
	})
	return n >= m.quorum, err
}

func (m *Mutex) valid(pool Pool) (bool, error) {
	conn := pool.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("GET", m.name))
	if err != nil {
		return false, err
	}
	return m.value == reply, nil
}

func genValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (m *Mutex) acquire(nodeID int,tryID int,pool Pool, value string) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("SET", m.name, value, "NX", "PX", int(m.expiry/time.Millisecond)))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return false, err
	}
	fmt.Printf("=====%v锁=====|host:%v|tryID:%v|获取锁成功\t",m.Id,configs[nodeID],tryID)
	return reply == "OK", nil
}

var deleteScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

func (m *Mutex) release(pool Pool, value string) (bool, error) {
	conn := pool.Get()
	defer conn.Close()
	status, err := redis.Int64(deleteScript.Do(conn, m.name, value))

	return err == nil && status != 0, err
}

var touchScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("pexpire", KEYS[1], ARGV[2])
	else
		return 0
	end
`)

func (m *Mutex) touch(pool Pool, value string, expiry int) (bool, error) {
	conn := pool.Get()
	defer conn.Close()
	status, err := redis.Int64(touchScript.Do(conn, m.name, value, expiry))

	return err == nil && status != 0, err
}

func (m *Mutex) actOnPoolsAsync(actFn func(int,Pool) (bool, error)) (int, error) {
	type result struct {
		Status bool
		Err    error
	}

	ch := make(chan result)
	//随机多个连接池异步获取锁资源
	for i, pool := range m.pools {
		go func(i int,pool Pool) {
			r := result{}
			r.Status, r.Err = actFn(i,pool)
			ch <- r
		}(i,pool)
	}
	n := 0
	var err error
	for range m.pools {
		r := <-ch
		if r.Status {
			n++
		} else if r.Err != nil {
			err = multierror.Append(err, r.Err)
		}
	}
	return n, err
}
