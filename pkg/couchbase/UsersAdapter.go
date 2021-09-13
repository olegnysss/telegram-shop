package couchbase

import (
	"github.com/couchbase/gocb/v2"
	"log"
)

var usersDocument = "Users"
var usersMap map[ID]User

func CheckUser(chatId int64, username string) (User, error) {
	user, ok := usersMap[ID(chatId)]
	if !ok {
		newUser := User{
			UserId: ID(chatId),
			Name:   username,
		}
		return AppendUser(newUser)
	} else {
		log.Printf("Пользователь %d уже зарегистрирован", chatId)
		return user, nil
	}
}

func AppendUser(newUser User) (User, error) {
	mops := []gocb.MutateInSpec{
		gocb.ArrayAppendSpec(usersDocument, newUser, nil),
	}

	_, err := collection.MutateIn(usersDocument, mops, &gocb.MutateInOptions{})
	if err != nil {
		return User{}, err
	}

	usersMap[newUser.UserId] = newUser
	log.Printf("New user %+v is added.", newUser)
	return usersMap[newUser.UserId], nil
}

func fetchUsers() (map[ID]User, error) {
	usersResult, err := collection.Get(usersDocument, nil)
	if err != nil {
		if kvErr, ok := err.(*gocb.KeyValueError); ok {
			_, err := handleKVError(kvErr, usersDocument, User{})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var users []User
	err = usersResult.Content(&users)
	if err != nil {
		return nil, err
	}

	for _, elem := range users {
		usersMap[elem.UserId] = elem
	}

	return usersMap, nil
}

func handleKVError(kvErr *gocb.KeyValueError, document string, types interface{}) (mutOut *gocb.MutationResult, errOut error) {
	if kvErr.ErrorDescription == "Not Found" {
		Logger.Printf("%s is not existed. Creating new document.", document)
		result, err := collection.Insert(document, types, nil)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, kvErr
}
