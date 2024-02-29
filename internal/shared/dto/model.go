package dto

import "time"

type DataType int

const (
	Simple DataType = iota
	Binary
	Card
)

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
	Type      DataType
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

type SimpleDataEncrypted struct {
	Common CommonData
	Data   []byte
}
