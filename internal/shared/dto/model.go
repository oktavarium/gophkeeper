package dto

import "time"

type Token struct {
	Id         string
	ValidUntil time.Time
}

type UserInfo struct {
	Login    string
	Password string
}

type SaveData struct {
	Token Token
	Name  string
	Data  string
}
