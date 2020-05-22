package datasource

import (
	"context"
	"fmt"
	"github.com/bsir2020/basework/configs"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

var (
	url     string
	logFile string
)

func init() {
	url = configs.EnvConfig.ES.Url
	logFile = configs.EnvConfig.ES.LogFile
}

type ESClient struct {
	client *elastic.Client
}

func GetEsClient() (client *ESClient, err error) {
	lf, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err.Error())
	}

	cfg := []elastic.ClientOptionFunc{
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetInfoLog(log.New(lf, "ES-INFO: ", 0)),
		elastic.SetTraceLog(log.New(lf, "ES-TRACE: ", 0)),
		elastic.SetErrorLog(log.New(lf, "ES-ERROR: ", 0)),
	}

	esclient, err := elastic.NewClient(cfg...)
	if esclient != nil {
		return &ESClient{
			client: esclient,
		}, nil
	}

	return nil, err
}

func (e *ESClient) Add(index, typ string, id string, data interface{}) (*elastic.IndexResponse, error) {
	return e.client.Index().Index(index).Type(typ).Id(id).BodyJson(data).Do(context.Background())
}

func (e *ESClient) Get(index, typ string, id string, data interface{}) (*elastic.IndexResponse, error) {
	return e.client.Index().Index(index).Type(typ).Id(id).Do(context.Background())
}
