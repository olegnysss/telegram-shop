package couchbase

import (
	"github.com/couchbase/gocb/v2"
	"log"
	"time"
)

var usersDocument = "Users"

var collection *gocb.Collection
var UsersMap map[ID]User

func AppendUser(newUser User) (User, error) {
	mops := []gocb.MutateInSpec{
		gocb.ArrayAppendSpec(usersDocument, newUser, nil),
	}

	_, err := collection.MutateIn(usersDocument, mops, &gocb.MutateInOptions{})
	if err != nil {
		return User{}, err
	}

	UsersMap[newUser.UserId] = newUser
	log.Printf("New user %+v is added.", newUser)
	return UsersMap[newUser.UserId], nil
}

func fetchUsers() {
	usersResult, err := collection.Get(usersDocument, nil)
	if err != nil {
		panic(err)
	}

	var users Users
	err = usersResult.Content(&users)
	if err != nil {
		panic(err)
	}

	for _, elem := range users.Users {
		UsersMap[elem.UserId] = elem
	}

	log.Printf("%+v", UsersMap)
}

func ConnectToCouch(config Config) error {
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
	UsersMap = make(map[ID]User)
	fetchUsers()
	return err
}
