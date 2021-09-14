package couchbase

import (
	"github.com/couchbase/gocb/v2"
	"log"
)

type UsersAdapter struct {
	UsersMap      map[ID]User
	entityAdapter *EntityAdapter
	UserDocument  string
}

func InitUsersAdapter() *UsersAdapter {
	return &UsersAdapter{
		map[ID]User{},
		&EntityAdapter{},
		"Users",
	}
}

func (u *UsersAdapter) CheckUser(chatId int64, username string) (User, error) {
	user, ok := u.UsersMap[ID(chatId)]
	if !ok {
		newUser := User{
			UserId: ID(chatId),
			Name:   username,
		}
		return u.AppendUser(newUser)
	} else {
		log.Printf("Пользователь %d уже зарегистрирован", chatId)
		return user, nil
	}
}

func (u *UsersAdapter) AppendUser(newUser User) (User, error) {
	mops := []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("", newUser, nil),
	}

	_, err := collection.MutateIn(u.UserDocument, mops, &gocb.MutateInOptions{})
	if err != nil {
		return User{}, err
	}

	u.UsersMap[newUser.UserId] = newUser
	log.Printf("New user %+v is added.", newUser)
	return u.UsersMap[newUser.UserId], nil
}

func (u *UsersAdapter) fetchUsers() (map[ID]User, error) {
	usersResult, err := u.entityAdapter.fetch(u.UserDocument)
	var users []User
	if usersResult == nil {
		return u.UsersMap, nil
	}

	err = usersResult.Content(&users)
	if err != nil {
		return nil, err
	}

	for _, elem := range users {
		u.UsersMap[elem.UserId] = elem
	}
	return u.UsersMap, nil
}

type User struct {
	UserId          ID
	Name            string
	Balance         float64
	PurchasesAmount float64
}
