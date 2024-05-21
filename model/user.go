package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type User struct {
	UserId      int64
	UserChat_id int64
	UserRole    Role
	UserLogin   string
	UserAlias   string
}

func NewUser(update tgbotapi.Update) User {
	return User{UserId: update.Message.From.ID,
		UserChat_id: update.Message.Chat.ID,
		UserRole:    Role{RoleName: "Клиент", RoleType: 1},
		UserLogin:   update.Message.From.UserName,
		UserAlias:   update.Message.Text}
}

type Users []User

type UsersInterface interface {
	IsExistUserById(user_id int64) bool
}

func (us Users) IsExistUserById(user_id int64) bool {
	for _, u := range us {
		if u.UserId == user_id {
			return true
		}
	}
	return false
}

func (us Users) FindUserById(user_id int64) *User {
	for _, u := range us {
		if u.UserId == user_id {
			return &u
		}
	}
	return nil
}

func (us Users) FindUserByIdAndUpdateAlias(user_id int64, newAlias string) {
	for i, u := range us {
		if u.UserId == user_id {
			u.UserAlias = newAlias
			us[i] = u
			return
		}
	}
}

type Role struct {
	RoleName string
	RoleType int
}

type Pair struct {
	PairUsers []User
}
