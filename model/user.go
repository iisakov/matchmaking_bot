package model

import (
	"matchmaking_bot/stl"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	UserId      int64
	UserChat_id int64
	UserRole    Role
	UserLogin   string
	UserAlias   string
	Answers     []string
	Gender      int
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

func (us Users) FindUserByIdAndAddAnswer(user_id int64, newAnswer string) {
	for i, u := range us {
		if u.UserId == user_id {
			if !u.FindAnswer(newAnswer) {
				u.Answers = append(u.Answers, newAnswer)
				us[i] = u
			}
			return
		}
	}
}

func (us Users) FindUserByIdAndUpdateAnswer(user_id int64, oldAnswers []string, newAnswer string) {
	for i, u := range us {
		if u.UserId == user_id {
			for _, oa := range oldAnswers {
				index, ok := stl.IndexElemInSliseString(u.Answers, oa)
				if !ok {
					continue
				}
				u.Answers = stl.DeleteElementByIndex(u.Answers, index)
			}
			u.Answers = append(u.Answers, newAnswer)
			us[i] = u
			return
		}
	}
}

func (us Users) FindUserByIdSetGender(user_id int64, gender string) {
	var genderInt int
	if gender == "Юноша" {
		genderInt = 1
	} else {
		genderInt = 0
	}

	for i, u := range us {
		if u.UserId == user_id {
			u.Gender = genderInt
			us[i] = u
			return
		}
	}
}

func (u User) FindAnswer(answer string) bool {
	for _, a := range u.Answers {
		if a == answer {
			return true
		}
	}
	return false
}

type Role struct {
	RoleName string
	RoleType int
}

type Pair struct {
	PairUsers []User
}
