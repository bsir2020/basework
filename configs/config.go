package configs

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Desc     string
	Pgsql    *PgSqlConfig
	Redis    *RedisConfig
	Rabbitmq *RabbitmqConfig
	Log      *LogConfig
	Authkey  *AuthkeyConfig
	Fillter  *FillterConfig
}

type PgSqlConfig struct {
	Hosts    []string
	Ports    []int
	User     string
	Password string
	Dbname   string
}

type RedisConfig struct {
	Hosts    []string
	Password string
}

type RabbitmqConfig struct {
	User     string   // mq user
	Password string   // mq password
	Ip       []string // mq ip
	Port     []int    // mq port
	Vhost    []string // vhost
	QuName   []string // 队列名称
	RtKey    []string // key值
	ExName   []string // 交换机名称
	ExType   []string // 交换机类型
}

type LogConfig struct {
	Logfile string
	Sqlog   string
}

type AuthkeyConfig struct {
	Key        string
	Subject    string
	PrivateKey string
	Publickey  string
}

type FillterConfig struct {
	Array []string
}

var (
	confPath  string
	env       string
	logfile   string
	sqlfile   string
	EnvConfig *Config
	WhiteList map[string]string
)

func init() {
	flag.StringVar(&env, "env", "dev", "set running env")
	flag.StringVar(&logfile, "logfile", "", "set log file")
	//flag.StringVar(&sqlfile, "sqllog", "sql", "set sql log file")

	confPath = "./configs/datasources-" + env + ".toml"

	_, err := toml.DecodeFile(confPath, &EnvConfig)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		WhiteList = make(map[string]string)

		if logfile != "" {
			EnvConfig.Log.Logfile = logfile
		}
		//EnvConfig.Log.Sqlog = "./" + sqlfile

		fmt.Println(EnvConfig.Desc)

		//
		for _, path := range EnvConfig.Fillter.Array {
			WhiteList[path] = path
		}
	}
}
