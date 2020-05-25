package datasource

import (
	"fmt"
	"github.com/bsir2020/basework/api"
	cfg "github.com/bsir2020/basework/configs"
	"github.com/bsir2020/basework/pkg/log"
	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"
	"time"
)

const (
	redisMaxIdle        = 30   //最大空闲连接数
	redisIdleTimeoutSec = 60 //最大空闲连接时间
	maxActive = 1000
)

var (
	redisURL      string
	redisPassword string
	db            int
	timeout =  2
)

func init() {
	redisURL = cfg.EnvConfig.Redis.Hosts[0]
	redisPassword = cfg.EnvConfig.Redis.Password
	db = cfg.EnvConfig.Redis.DB
}

func newRedisPool() (redisPool *redis.Pool) {
	logger := log.New()

	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		MaxActive:   maxActive,
		IdleTimeout: redisIdleTimeoutSec * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL, redis.DialDatabase(db), redis.DialPassword(redisPassword),
				redis.DialConnectTimeout(time.Duration(timeout) * time.Second),
				redis.DialReadTimeout(time.Duration(timeout) *time.Second),
				redis.DialWriteTimeout(time.Duration(timeout) *time.Second))
			if err != nil {
				logger.Error("RedisPool", zap.String(api.RedisConnErr.Message, err.Error()))
				return nil, fmt.Errorf("redis connection error: %s", err)
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) (err error) {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err = c.Do("PING")
			if err != nil {
				logger.Error("RedisPool", zap.String(api.RedisConnErr.Message, err.Error()))
			}
			return err
		},
	}
}

func GetRedisConn() (redis.Conn, *api.Errno) {
	return newRedisPool().Get(), nil
}

//key:"lock_uid"
//uid: user_id
func AddLock(val string) bool {
	msg, err := redis.String(
		Exec("set", "lock:LOCK_"+val, val, "nx", "ex", 4),
	)

	if err == redis.ErrNil {
		return false
	}

	if msg == "OK" {
		return true
	}

	return false
}

func DelLock(val string) {
	_, err := Exec("del", "lock:LOCK_"+val)
	if err != nil{
		fmt.Println(api.RedisConnErr, err.Error())
	}
}

//func GetLock(conn redis.Conn, val string) string {
//	defer conn.Close()
//
//	msg, _ := redis.String(conn.Do("get", "lock:LOCK_"+val))
//	return msg
//}

func Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con, _ := GetRedisConn()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()
	parmas := make([]interface{}, 0)
	parmas = append(parmas, key)

	if len(args) > 0 {
		for _, v := range args {
			parmas = append(parmas, v)
		}
	}
	return con.Do(cmd, parmas...)
}
