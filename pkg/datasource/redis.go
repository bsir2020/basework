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
	redisMaxIdle        = 3   //最大空闲连接数
	redisIdleTimeoutSec = 100 //最大空闲连接时间
)

var (
	redisURL      string
	redisPassword string
	db            int
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
		IdleTimeout: redisIdleTimeoutSec * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL, redis.DialDatabase(db), redis.DialPassword(redisPassword))
			if err != nil {
				logger.Error("RedisPool", zap.String(api.RedisConnErr.Message, err.Error()))
				return nil, fmt.Errorf("redis connection error: %s", err)
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				logger.Error("RedisPool", zap.String(api.RedisConnErr.Message, err.Error()))

				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

func GetRedisConn() (redis.Conn, *api.Errno) {
	var pool = newRedisPool()
	if conn, err := pool.Dial(); err != nil {
		fmt.Println(api.RedisConnErr, err.Error())
		return nil, api.RedisConnErr
	} else {
		return conn, nil
	}
}

//key:"lock_uid"
//uid: user_id
func AddLock(conn redis.Conn, key, val string, ex int) bool {
	msg, err := redis.String(
		conn.Do("set", key, val, "nx", "ex", 5),
	)

	if err == redis.ErrNil {
		return false
	}

	if msg == "OK" {
		return true
	}

	return false
}

func DelLock(conn redis.Conn, key, requestId string) bool {
	if GetLock(conn, key) == requestId {
		msg, _ := redis.Int64(conn.Do("del", key))
		if msg == 1 || msg == 0 {
			return true
		}
		return false
	}
	return false
}

func GetLock(conn redis.Conn, key string) string {
	msg, _ := redis.String(conn.Do("get", key))
	return msg
}
