package couchbase

type ID int64
type Roubles float32

type Config struct {
	ConnString    string
	CouchUsername string
	CouchPassword string
	BucketName    string
}

type Transaction struct {
	TxnId    ID
	PersonId ID
	Sum      uint32
	Comment  string
}

type User struct {
	UserId          ID
	Name            string
	Balance         Roubles
	PurchasesAmount Roubles
}

type Users struct {
	Users []User
}
