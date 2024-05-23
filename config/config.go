package config

import (
	"matchmaking_bot/model"
	"os"
	"strconv"
)

var TOKEN string = os.Getenv("MATCHMAKER_BOT_TOKEN")
var TEST_BOT_CHAT, _ = strconv.Atoi(os.Getenv("TEST_BOT_CHAT"))
var MODERATOR_BOT_CHAT, _ = strconv.Atoi(os.Getenv("MODERATOR_BOT_CHAT"))

var ADMINS, MODERATORS, CUSTOMERS model.Users

var PAIRS model.Pairs

func MockUsers(s model.Users) model.Users {
	return append(
		s,
		model.User{UserId: 1, UserChat_id: 1, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST1", UserAlias: "ATEST1", Answers: []string{"1", "2", "3", "4"}, Gender: 1},
		model.User{UserId: 2, UserChat_id: 2, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST2", UserAlias: "ATEST2", Answers: []string{"5", "6", "7", "8"}, Gender: 1},
		model.User{UserId: 3, UserChat_id: 3, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST3", UserAlias: "ATEST3", Answers: []string{"9", "10", "11", "12"}, Gender: 1},
		model.User{UserId: 4, UserChat_id: 4, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST4", UserAlias: "ATEST4", Answers: []string{"1", "3", "5", "7"}, Gender: 1},
		model.User{UserId: 5, UserChat_id: 5, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST5", UserAlias: "ATEST5", Answers: []string{"9", "11", "13", "15"}, Gender: 1},
		model.User{UserId: 6, UserChat_id: 6, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST6", UserAlias: "ATEST6", Answers: []string{"2", "4", "8", "10"}, Gender: 1},
		model.User{UserId: 7, UserChat_id: 7, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST7", UserAlias: "ATEST7", Answers: []string{"5", "6", "7", "8"}, Gender: 0},
		model.User{UserId: 8, UserChat_id: 8, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST8", UserAlias: "ATEST8", Answers: []string{"6", "9", "12", "15"}, Gender: 0},
		model.User{UserId: 9, UserChat_id: 9, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST9", UserAlias: "ATEST9", Answers: []string{"2", "5", "8", "11"}, Gender: 0},
		model.User{UserId: 10, UserChat_id: 10, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST10", UserAlias: "ATEST10", Answers: []string{"3", "6", "9", "12"}, Gender: 0},
		model.User{UserId: 11, UserChat_id: 11, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST11", UserAlias: "ATEST11", Answers: []string{"4", "7", "10", "13"}, Gender: 0},
		model.User{UserId: 12, UserChat_id: 12, UserRole: model.Role{RoleName: "Клиент", RoleType: 1}, UserLogin: "TEST12", UserAlias: "ATEST12", Answers: []string{"5", "8", "11", "14"}, Gender: 0})
}
