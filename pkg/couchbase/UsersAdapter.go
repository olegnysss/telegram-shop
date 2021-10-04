package couchbase

import (
	"github.com/couchbase/gocb/v2"
	"strconv"
)

var usersCollection *gocb.Collection

type UsersAdapter struct {
	entityAdapter  *EntityAdapter
	UserCollection string
}

func InitUsersAdapter() *UsersAdapter {
	return &UsersAdapter{
		&EntityAdapter{},
		"users",
	}
}

func InitUsersCollection(scope *gocb.Scope) *gocb.Collection {
	usersCollection = scope.Collection("users")
	return usersCollection
}

func (u *UsersAdapter) CheckUser(chatId int64, username string) (User, error) {
	userResult, err := u.retrieve(chatId)
	if err != nil {
		if kvErr, ok := err.(*gocb.KeyValueError); ok && kvErr.ErrorDescription == "Not Found" {
			_, err := u.registerUser(chatId, username)
			if err != nil {
				return User{}, err
			}
			newUserResult, err := u.entityAdapter.retrieve(usersCollection, int(chatId))
			return ContentUser(newUserResult)
		} else {
			return User{}, err
		}
	}
	return ContentUser(userResult)
}

func (u *UsersAdapter) retrieve(chatId int64) (*gocb.GetResult, error) {
	return u.entityAdapter.retrieve(usersCollection, int(chatId))
}

func ContentUser(userResult *gocb.GetResult) (User, error) {
	var user User
	err := userResult.Content(&user)
	if err != nil {
		panic(err)
	}
	return user, nil
}

func (u *UsersAdapter) registerUser(chatId int64, username string) (*gocb.MutationResult, error) {
	newUser := User{
		UserId: ID(chatId),
		Name:   username,
	}
	result, err := u.entityAdapter.insert(usersCollection, newUser)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type User struct {
	UserId          ID
	Name            string
	Balance         float64
	PurchasesAmount float64
}

func (u User) GetId() ID {
	return u.UserId
}

func (u User) GetStringId() string {
	return strconv.Itoa(int(u.UserId))
}
