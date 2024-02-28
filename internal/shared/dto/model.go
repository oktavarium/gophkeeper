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

type CommonData struct {
	IsDeleted bool
	Modified  time.Time
}

type SimpleCardData struct {
	Common CommonData
	Data   SimpleCardRecord
}

type SimpleCardRecord struct {
	Name       string
	Number     string
	CVV        uint32
	ValidUntil time.Time
}

type SimpleCardDataEncrypted struct {
	Common CommonData
	Data   []byte
}
