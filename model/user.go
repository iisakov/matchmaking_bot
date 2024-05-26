package model

import (
	"errors"
	"fmt"
	"matchmaking_bot/stl"
	"strconv"

	tgbotapi "github.com/iisakov/telegram-bot-api"
)

type User struct {
	UserId        int64
	UserChat_id   int64
	UserRole      Role
	UserLogin     string
	UserAlias     string
	Answers       []string
	Gender        int
	LastMessageId int
}

func NewUser(update tgbotapi.Update) User {
	return User{UserId: update.Message.From.ID,
		UserChat_id: update.Message.Chat.ID,
		UserRole:    Role{RoleName: "Клиент", RoleType: 1},
		UserLogin:   update.Message.From.UserName,
		UserAlias:   update.Message.Text}
}

type Users []User

func (us Users) GetUsers() (result string) {
	result = ""
	for _, u := range us {
		result += strconv.Itoa(u.Gender) + " "
		result += u.UserLogin + " (" + u.UserAlias + "): "
		for _, a := range u.Answers {
			result += a + ", "
		}
		result += fmt.Sprintf("lust_message_id: %d", u.LastMessageId)
		result += "\n"
	}
	return result
}

func (us Users) GetUsersByGender(gender int) (result Users) {
	for _, u := range us {
		if u.Gender == gender {
			result = append(result, u)
		}
	}

	return
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

func (us Users) FindUserByIdAndSetLastMessageId(user_id int64, message_id int) {
	for i, u := range us {
		if u.UserId == user_id {
			u.LastMessageId = message_id
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
	NumMatches int
	PairUsers  Users
}

func (p Pair) GetIdUsers() (result [2]int64) {
	return [2]int64{p.PairUsers[0].UserChat_id, p.PairUsers[1].UserChat_id}
}

func NewPair(us Users, nm int) (Pair, error) {
	if len(us) == 2 {
		return Pair{PairUsers: us, NumMatches: nm}, nil
	}
	return Pair{}, errors.New("слишком много элементов")
}

func IsExistUserInPairs(pairs Pairs, user User) bool {
	for _, pair := range pairs {
		for _, uId := range pair.GetIdUsers() {
			if uId == user.UserChat_id {
				return true
			}
		}
	}
	return false
}

type Pairs []Pair

func (ps Pairs) GetConversationPartner(uId int64) (result int64, ok bool) {
	for _, p := range ps {
		for i, pUId := range p.GetIdUsers() {
			if pUId == uId {
				if i == 0 {
					return p.GetIdUsers()[1], true
				} else {
					return p.GetIdUsers()[0], true
				}

			}
		}
	}
	return -1, false
}

func (ps Pairs) GetPairs() string {
	result := ""
	for _, p := range ps {
		result += "совпадений: " + strconv.Itoa(p.NumMatches) + "\n"

		result += strconv.Itoa(p.PairUsers[0].Gender) + " "
		result += p.PairUsers[0].UserLogin + " (" + p.PairUsers[0].UserAlias + "): "
		for _, a := range p.PairUsers[0].Answers {
			result += a + ", "
		}
		result += "\n"

		result += strconv.Itoa(p.PairUsers[1].Gender) + " "
		result += p.PairUsers[1].UserLogin + " (" + p.PairUsers[1].UserAlias + "): "
		for _, a := range p.PairUsers[1].Answers {
			result += a + ", "
		}

		result += "\n"
		result += "\n"
	}
	return result
}

func CheckCompatibility(sU1, sU2 Users) (result Pairs) {
	var biggerS, smallerS Users
	var maxMatches = 0
	var subResult = make(map[int][]Users)
	if len(sU1) > len(sU2) {
		biggerS = sU1
		smallerS = sU2
	} else {
		biggerS = sU2
		smallerS = sU1
	}

	for _, vS := range smallerS {
		for _, vB := range biggerS {
			numMatches := stl.GetNumberMatches(vS.Answers, vB.Answers)
			if maxMatches < numMatches {
				maxMatches = numMatches
			}

			if _, ok := subResult[numMatches]; !ok {
				subResult[numMatches] = []Users{}
			}
			subResult[numMatches] = append(subResult[numMatches], Users{vS, vB})
		}
	}

	for i := maxMatches; i >= 0; i-- {
		for _, subPair := range subResult[i] {
			if IsExistUserInPairs(result, subPair[0]) || IsExistUserInPairs(result, subPair[1]) {
				continue
			}
			np, _ := NewPair(subPair, i)
			result = append(result, np)
		}
	}
	return
}
