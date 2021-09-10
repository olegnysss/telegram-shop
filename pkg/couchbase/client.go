package couchbase

import (
	"github.com/couchbase/gocb/v2"
	"time"
)

var collection *gocb.Collection

func ConnectToCouch(config Config) (map[ID]User, error) {
	cluster, err := gocb.Connect(
		config.ConnString,
		gocb.ClusterOptions{
			Username: config.CouchUsername,
			Password: config.CouchPassword,
		})
	if err != nil {
		panic(err)
	}
	// get a bucket reference
	bucket := cluster.Bucket(config.BucketName)

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		panic(err)
	}

	collection = bucket.DefaultCollection()
	usersMap = make(map[ID]User)
	transactionsMap = make(map[ID][]Transaction)
	return fetchUsers()
}
