package rds

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"

	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

// RedisSource -
type RedisSource struct {
	dbpool *redis.Pool
	psc    *redis.PubSubConn
}

var redisSourceMap = make(map[string]*RedisSource)
var defaultServer string

func newPool(server, password string, maxIdle int) *redis.Pool {
	if maxIdle == 0 {
		maxIdle = 128
	}
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if len(password) == 0 {
				return c, nil
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// AddRedisServer -
func AddRedisServer(server, password string, maxIdle int) bool {
	redisSource := new(RedisSource)
	logs.Debug("|foundation|rds|redis|AddRedisServer|server:%v,password:%v,maxIdle:%v", server, password, maxIdle)
	redisSource.dbpool = newPool(server, password, maxIdle)
	logs.Debug("|foundation|rds|redis|AddRedisServer|server:%v,password:%v,maxIdle:%v|newPool|%+v", server, password, maxIdle, redisSource.dbpool)
	redisSourceMap[server] = redisSource
	return true
}

// InitRedisServer -
func InitRedisServer(server, password string, maxIdle int) {
	defaultServer = server
	AddRedisServer(server, password, maxIdle)
}

// Init -
func Init() {
	logs.Debug("|foundation|init|rds|Init")
	//
	redissource := beego.AppConfig.DefaultString("redis","192.168.85.109:6379,150,123456")
	if rdss := strings.Split(redissource, ","); len(rdss) == 3 {
		// address,connect,password
		maxIdle, _ := strconv.Atoi(rdss[1])
		InitRedisServer(rdss[0], rdss[2], maxIdle)
	}
}

// GetRedisByServerName -
func GetRedisByServerName(server string) *RedisSource {
	if v, ok := redisSourceMap[server]; ok {
		return v
	}
	logs.Warn("Not Found: %s", server)
	return nil
}

// GetRedisDefault -
func GetRedisDefault() *RedisSource {
	if len(defaultServer) == 0 {
		logs.Debug("|foundation|rds|redis|GetRedisDefault|len(defaultServer)==0|redisSourceMap|%+v", redisSourceMap)
		for _, s := range redisSourceMap {
			return s
		}
	}
	return redisSourceMap[defaultServer]
}

// Ping -
func (rs *RedisSource) Ping() error {
	_, err := rs.Do("PING")
	return err
}

// GetConn -
func (rs *RedisSource) GetConn() redis.Conn {
	c := rs.dbpool.Get()
	// 统计redis 连接数
	return c
}

// CloseConn -
func (rs *RedisSource) CloseConn(conn redis.Conn) (err error) {
	err = conn.Close()
	// 统计redis 连接数

	return
}

// TTL - 过期时间
func (rs *RedisSource) TTL(key string) (int, error) {
	return redis.Int(rs.Do("TTL", key))
}

// Do -
func (rs *RedisSource) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := rs.GetConn()
	defer rs.CloseConn(c)

	return c.Do(commandName, args...)
}

// Incr -
func (rs *RedisSource) Incr(key string) uint64 {
	n, _ := redis.Uint64(rs.Do("INCR", key))
	return n
}

// IncrEx -
func (rs *RedisSource) IncrEx(key string) (uint64, error) {
	return redis.Uint64(rs.Do("INCR", key))
}

// IncrBy -
func (rs *RedisSource) IncrBy(key string, val int) uint64 {
	n, _ := redis.Uint64(rs.Do("INCRBY", key, val))
	return n
}

// Del -
func (rs *RedisSource) Del(key string) error {
	_, err := rs.Do("DEL", key)
	return err
}

// Set -
func (rs *RedisSource) Set(key string, val interface{}, expire int) error {
	if expire == 0 {
		_, err := rs.Do("SET", key, val)
		return err
	}
	_, err := rs.Do("SET", key, val, "EX", expire)
	return err
}

// SetNx -
func (rs *RedisSource) SetNx(key string, val interface{}, expire int) (bool, error) {
	var (
		reply string
		err   error
	)

	if expire == 0 {
		reply, err = redis.String(rs.Do("SET", key, val, "NX"))
	} else {
		reply, err = redis.String(rs.Do("SET", key, val, "EX", expire, "NX"))
	}

	if "OK" == reply {
		return true, err
	}
	return false, err
}

// SetPx - ms
func (rs *RedisSource) SetPx(key string, val interface{}, pexpire int) error {
	if pexpire == 0 {
		_, err := rs.Do("SET", key, val)
		return err
	}
	_, err := rs.Do("SET", key, val, "PX", pexpire)
	return err
}

// SetExpire -
func (rs *RedisSource) SetExpire(key string, expire int) error {
	_, err := rs.Do("EXPIRE", key, expire)
	return err
}

// SetPexpire - ms
func (rs *RedisSource) SetPexpire(key string, expire int) error {
	_, err := rs.Do("PEXPIRE", key, expire)
	return err
}

// GetUint64 -
func (rs *RedisSource) GetUint64(key string) (uint64, error) {
	return redis.Uint64(rs.Do("GET", key))
}

// GetUint32 -
func (rs *RedisSource) GetUint32(key string) (uint32, error) {
	v, err := rs.GetUint64(key)
	return uint32(v), err
}

// GetInt64 -
func (rs *RedisSource) GetInt64(key string) (int64, error) {
	return redis.Int64(rs.Do("GET", key))
}

// GetInt64Ex -
func (rs *RedisSource) GetInt64Ex(key string, def int64) (int64, error) {
	v, err := redis.Int64(rs.Do("GET", key))
	if err == redis.ErrNil {
		return def, nil
	}
	return v, err
}

// GetString -
func (rs *RedisSource) GetString(key string) (string, error) {
	return redis.String(rs.Do("GET", key))
}

// LRANGE -
func (rs *RedisSource) LRANGE(key string, start, stop int) ([]string, error) {
	return redis.Strings(rs.Do("LRANGE", key, start, stop))
}

// SMembers -
func (rs *RedisSource) SMembers(key string) ([]string, error) {
	return redis.Strings(rs.Do("SMEMBERS", key))
}

// SMembersUint64 -
func (rs *RedisSource) SMembersUint64(key string) ([]uint64, error) {
	reply, err := rs.Do("SMEMBERS", key)

	var ints []uint64
	if reply == nil {
		return ints, redis.ErrNil
	}
	values, err := redis.Values(reply, err)
	if err != nil {
		return ints, err
	}
	if err := redis.ScanSlice(values, &ints); err != nil {
		return ints, err
	}
	return ints, nil
}

// SAdd -
func (rs *RedisSource) SAdd(key string, val interface{}) error {
	_, err := redis.Strings(rs.Do("SADD", key, val))
	return err
}

// SPop -
func (rs *RedisSource) SPop(key string) (string, error) {
	return redis.String(rs.Do("SPOP", key))
}

// SPopInt64 -
func (rs *RedisSource) SPopInt64(key string) (int64, error) {
	return redis.Int64(rs.Do("SPOP", key))
}

// SCard -
func (rs *RedisSource) SCard(key string) (int, error) {
	return redis.Int(rs.Do("SCARD", key))
}

// SRem -
func (rs *RedisSource) SRem(key string, val interface{}) error {
	_, err := redis.Strings(rs.Do("SREM", key, val))
	return err
}

// GetBytes -
func (rs *RedisSource) GetBytes(key string) ([]byte, error) {
	raw, err := rs.Do("GET", key)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}
	return redis.Bytes(raw, err)
}

// Get -
func (rs *RedisSource) Get(key string) (interface{}, error) {
	return rs.Do("GET", key)
}

// SetJSON -
func (rs *RedisSource) SetJSON(key string, val interface{}, expire int) error {
	jsdata, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return rs.Set(key, jsdata, expire)
}

// GetJSON -
func (rs *RedisSource) GetJSON(key string, val interface{}) (interface{}, error) {
	v, err := redis.String(rs.Do("GET", key))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(v), val)
	if err != nil {
		return nil, err
	}

	return val, nil
}

// LPush -
func (rs *RedisSource) LPush(key string, val interface{}) bool {
	jsdata, err := json.Marshal(val)
	if err != nil {
		return false
	}
	_, err = rs.Do("LPUSH", key, jsdata)
	return err == nil
}

// RPush -
func (rs *RedisSource) RPush(key string, val interface{}) bool {
	jsdata, err := json.Marshal(val)
	if err != nil {
		return false
	}
	_, err = rs.Do("RPUSH", key, jsdata)
	return err == nil
}

// LPushStr -
func (rs *RedisSource) LPushStr(key string, val interface{}) (int, error) {
	return redis.Int(rs.Do("LPUSH", key, val))
}

// LLen -
func (rs *RedisSource) LLen(key string) uint64 {
	l, _ := redis.Uint64(rs.Do("LLEN", key))
	return l
}

// LTrim -
func (rs *RedisSource) LTrim(key string, start, end int) (bool, error) {
	return redis.Bool(rs.Do("LTRIM", key, start, end))
}

// LPop -
func (rs *RedisSource) LPop(key string, val interface{}) bool {
	raw, err := rs.Do("LPOP", key)
	if err != nil || raw == nil {
		return false
	}
	valbytes, err := redis.Bytes(raw, err)
	if err != nil {
		return false
	}

	err = json.Unmarshal(valbytes, val)
	if err != nil {
		return false
	}
	return true
}

// RPop -
func (rs *RedisSource) RPop(key string, val interface{}) bool {
	raw, err := rs.Do("RPOP", key)
	if err != nil || raw == nil {
		return false
	}
	valbytes, err := redis.Bytes(raw, err)
	if err != nil {
		return false
	}

	err = json.Unmarshal(valbytes, val)
	if err != nil {
		return false
	}
	return true
}

// LPopInt64 -
func (rs *RedisSource) LPopInt64(key string) (int64, error) {
	return redis.Int64(rs.Do("LPOP", key))
}

// ZINCRBY -
func (rs *RedisSource) ZINCRBY(key string, v interface{}, element string) error {
	_, err := rs.Do("ZINCRBY", key, v, element)
	return err
}

// ZRANK -
func (rs *RedisSource) ZRANK(key, element string) (int, error) {
	return redis.Int(rs.Do("ZRANK", key, element))
}

// Keys -
func (rs *RedisSource) Keys(key string) ([]string, error) {
	return redis.Strings(rs.Do("KEYS", key))
}

// Hkeys -
func (rs *RedisSource) Hkeys(key string) ([]string, error) {
	return redis.Strings(rs.Do("HKEYS", key))
}

// Uint64 -
func (rs *RedisSource) Uint64(command, key string) uint64 {
	v, _ := redis.Uint64(rs.Do(command, key))
	return v
}

// HGetStr -
func (rs *RedisSource) HGetStr(key, field string) (string, error) {
	return redis.String(rs.Do("HGET", key, field))
}

// HINCRBY -
func (rs *RedisSource) HINCRBY(key, field string, incr int) error {
	_, err := rs.Do("HINCRBY", key, field, incr)
	return err
}

// HGetInt64 -
func (rs *RedisSource) HGetInt64(key, field string) (int64, error) {
	return redis.Int64(rs.Do("HGET", key, field))
}

// HGetAllInt64Map -
func (rs *RedisSource) HGetAllInt64Map(key string) (map[string]int64, error) {
	return redis.Int64Map(rs.Do("HGETALL", key))
}

// HSet -
func (rs *RedisSource) HSet(key, field string, val interface{}) error {
	_, err := rs.Do("HSET", key, field, val)
	return err

}

// HDel -
func (rs *RedisSource) HDel(key, field string) error {
	_, err := rs.Do("HDEL", key, field)
	return err
}

// SIsMember  -
func (rs *RedisSource) SIsMember(key, field string) (bool, error) {
	resInt, err := redis.Int(rs.Do("SISMEMBER", key, field))
	if 1 == resInt {
		return true, err
	}
	return false, err
}

/******************************************** 发布与订阅 ******************************************************/


// NewPubSubCoon -
func (rs *RedisSource) NewPubSubCoon() *redis.PubSubConn {
	pubSubConn := &redis.PubSubConn{Conn: rs.dbpool.Get()}
	// 统计redis 连接数
	return pubSubConn
}

// Subscribe -
func (rs *RedisSource) Subscribe(channel string) error {
	if rs.psc == nil {
		rs.psc = rs.NewPubSubCoon()
	}
	return rs.psc.Subscribe(channel)
}

// PSubscribe -
func (rs *RedisSource) PSubscribe(pattern string) error {
	if rs.psc == nil {
		rs.psc = &redis.PubSubConn{Conn: rs.dbpool.Get()}
		// 统计redis 连接数
	}
	return rs.psc.PSubscribe(pattern)
}

// Receive -
func (rs *RedisSource) Receive() interface{} {
	if rs.psc == nil {
		return fmt.Errorf("please subscribe first")
	}
	return rs.psc.Receive()
}

// Publish -
func (rs *RedisSource) Publish(channel, value interface{}) error {
	c := rs.GetConn()
	defer rs.CloseConn(c)
	_, err := c.Do("PUBLISH", channel, value)
	return err
}

/******************************************** 位图 ******************************************************/
//SetBit-
func (rs *RedisSource)SetBit(key string,offset int64,status bool)error{
	var value int
	if status{
		value =1
	}
	_,err := rs.Do("setbit",key,offset,value)
	return err
}

//GetBit -
func (rs *RedisSource)GetBit(key string,offset int64) (int,error){
	return redis.Int(rs.Do("getbit",key,offset))
}

/*
 *BitCount -
 *@param key string
 *@param cods []interface{} [start,end]
 *@return
 */
func (rs *RedisSource)BitCount(key string,cods ...interface{})(int,error){
	if len(cods) != 2{
		return redis.Int(rs.Do("bitcount",key))
	}
	var args []interface{}
	args = append(args,key)
	args = append(args,cods...)
	return redis.Int(rs.Do("bitcount",args...))
}