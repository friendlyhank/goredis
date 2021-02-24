package redsync

import "github.com/gomodule/redigo/redis"

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

// A Pool maintains a pool of Redis connections.
type Pool interface {
	Get() redis.Conn
}
