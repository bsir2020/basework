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
	logger := log.New()

	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeoutSec * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL)
			if err != nil {
				logger.Fatal("newRedisPool", zap.String("redis connection error", err.Error()))
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			//验证redis密码
			if _, authErr := c.Do("AUTH", redisPassword); authErr != nil {
				logger.Fatal("newRedisPool", zap.String("redis connection error", err.Error()))

				return nil, fmt.Errorf("redis auth password error: %s", authErr)
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				logger.Fatal("newRedisPool", zap.String("redis connection error", err.Error()))

				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

func GetRedisConn() redis.Conn {
	return newRedisPool().Get()
}
