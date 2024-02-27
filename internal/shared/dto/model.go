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

type CommonData struct {
	Name      string
	IsDeleted bool
	Modified  time.Time
}

type SimpleCardRecord struct {
	Common     CommonData
	Number     string
	CVV        uint32
	ValidUntil time.Time
}
