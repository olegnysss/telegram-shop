package couchbase

import (
	"github.com/couchbase/gocb/v2"
	"github.com/olegnysss/telebot_qiwi/pkg/config"
	"log"
	"os"
	"time"
)

var Logger *log.Logger
var cluster *gocb.Cluster

type ID int64

type CouchClient struct {
	connString    string
	couchUsername string
	couchPassword string
	bucketName    string

	UsersAdapter        *UsersAdapter
	TransactionsAdapter *TransactionsAdapter
}

func InitCouchClient(config config.Couch) *CouchClient {
	return &CouchClient{
		connString:          config.CouchConnString,
		couchUsername:       config.CouchUsername,
		couchPassword:       config.CouchPassword,
		bucketName:          config.CouchBucketName,
		UsersAdapter:        InitUsersAdapter(),
		TransactionsAdapter: InitTransactionsAdapter(),
	}
}

func (c *CouchClient) ConnectToCouch() (*gocb.Scope, error) {
	err := initLogs()
	if err != nil {
		log.Panic(err)
	}

	cluster, err = gocb.Connect(
		c.connString,
		gocb.ClusterOptions{
			Username: c.couchUsername,
			Password: c.couchPassword,
		})
	if err != nil {
		log.Panic(err)
	}
	// get a bucket reference
	bucket := cluster.Bucket(c.bucketName)

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Panic(err)
	}

	scope := bucket.DefaultScope()
	InitUsersCollection(scope)
	InitTnxCollection(scope)
	return scope, nil
}

func initLogs() error {
	file, err := os.OpenFile("couch.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	Logger = log.New(file, "COUCH: ", log.Ldate|log.Ltime|log.Lshortfile)
	return err
}
