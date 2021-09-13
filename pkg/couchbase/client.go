package couchbase

import (
	"github.com/couchbase/gocb/v2"
	"log"
	"os"
	"time"
)

var collection *gocb.Collection
var Logger *log.Logger

func ConnectToCouch(config Config) (map[ID]User, error) {
	err := initLogs()
	if err != nil {
		log.Panic(err)
	}

	cluster, err := gocb.Connect(
		config.ConnString,
		gocb.ClusterOptions{
			Username: config.CouchUsername,
			Password: config.CouchPassword,
		})
	if err != nil {
		log.Panic(err)
	}
	// get a bucket reference
	bucket := cluster.Bucket(config.BucketName)

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Panic(err)
	}

	collection = bucket.DefaultCollection()
	usersMap = make(map[ID]User)
	transactionsMap = make(map[ID][]Transaction)
	return fetchUsers()
}

func initLogs() error {
	file, err := os.OpenFile("couch.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	Logger = log.New(file, "TG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return err
}
