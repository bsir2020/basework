package datasource

import (
	"github.com/bsir2020/basework/configs"
	elastic "github.com/elastic/go-elasticsearch/v7"
)

var (
	url     string
	logFile string
	user string
	passwd string
)

func init() {
	url = configs.EnvConfig.ES.Url
	logFile = configs.EnvConfig.ES.LogFile
	user = configs.EnvConfig.ES.User
	passwd = configs.EnvConfig.ES.Passwd
}

type ESClient struct {
	client *elastic.Client
}

func GetEsClient() (client *ESClient, err error) {
	esConfig := &elastic.Config{
		Addresses: []string{url},
		Username: user,
		Password: passwd,
	}

	esclient, err := elastic.NewClient(*esConfig)
	if esclient != nil {
		return &ESClient{
			client: esclient,
		}, nil
	}

	return nil, err
}

