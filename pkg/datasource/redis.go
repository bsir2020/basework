package datasource

import (
	"fmt"
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
)

func init() {
	redisURL = cfg.EnvConfig.Redis.Hosts[0]
	redisPassword = cfg.EnvConfig.Redis.Password
}

func newRedisPool() *redis.Pool {
	logger := log.GetLogger()

	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeoutSec * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL)
			if err != nil {
				loger(logger, "redis connection error", err)
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			//验证redis密码
			if _, authErr := c.Do("AUTH", redisPassword); authErr != nil {
				loger(logger, "redis auth password error", err)

				return nil, fmt.Errorf("redis auth password error: %s", authErr)
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				loger(logger, "ping redis error", err)

				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

func GetRedisConn() redis.Conn {
	return newRedisPool().Get()
}

func loger(logger *zap.Logger, msg string, err error) {
	logger.Fatal(msg, zap.String("err", err.Error()))
}
