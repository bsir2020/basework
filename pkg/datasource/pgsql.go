package datasource

import (
	"database/sql"
	"fmt"
	cfg "github.com/bsir2020/basework/configs"
	"github.com/bsir2020/basework/pkg/log"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
	"time"
	"xorm.io/core"
)

var (
	host     string
	port     int
	user     string
	password string
	dbname   string
	logfile  string
)

func init() {
	host = cfg.EnvConfig.Pgsql.Hosts[0]
	port = cfg.EnvConfig.Pgsql.Ports[0]
	user = cfg.EnvConfig.Pgsql.User
	password = cfg.EnvConfig.Pgsql.Password
	dbname = cfg.EnvConfig.Pgsql.Dbname
	logfile = cfg.EnvConfig.Log.Sqlog
}

func GetPGSql() (pgEngine *xorm.Engine) {
	logger := log.New()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	pgEngine, err := xorm.NewEngine("postgres", psqlInfo)
	if err != nil {
		println(err.Error())
		logger.Fatal("GetPGSql", zap.String("create engine failed", err.Error()))

		return
	}

	// 设置日志
	logFile, err := os.Create(logfile)
	if err != nil {
		logger.Fatal("GetPGSql", zap.String("create sql.log failed", err.Error()))

		println(err.Error())
		return
	}

	pgEngine.Logger().SetLevel(core.LOG_DEBUG)
	pgEngine.SetLogger(xorm.NewSimpleLogger(logFile))
	pgEngine.SetMaxIdleConns(10)
	pgEngine.SetMaxOpenConns(1000)
	pgEngine.SetConnMaxLifetime(time.Second * 10)
	pgEngine.ShowExecTime(true)
	pgEngine.ShowSQL(true)

	if err = pgEngine.Ping(); err != nil {
		logger.Fatal("GetPGSql", zap.String("database connect", err.Error()))

		fmt.Printf("database connect failed : %s", err.Error())
	} else {
		logger.Info("GetPGSql", zap.String("database connect ok", err.Error()))
		fmt.Printf("database connect ok")
	}

	return
}

func GetPG() (db *sql.DB) {
	logger := log.New()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatal("GETPG", zap.String("database connect failed", err.Error()))
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal("GETPG", zap.String("database connect failed", err.Error()))
	} else {
		logger.Info("GETPG", zap.String("database connect ok", ""))
		fmt.Printf("database connect ok")
	}

	return
}
