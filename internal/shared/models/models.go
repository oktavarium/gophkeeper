package models

import "time"

type DataType int

const (
	Simple DataType = iota
	Binary
	Card
)

func DataTypeToString(t DataType) string {
	switch t {
	case Simple:
		return "Simple"
	case Binary:
		return "File"
	case Card:
		return "Card"
	default:
		return "Unspecified"
	}
}

type Token struct {
	ID         string
	ValidUntil time.Time
}

type UserInfo struct {
	Login    string
	Password string
}

type CommonData struct {
	Deleted  bool
	Modified time.Time
	Type     DataType
}

type SimpleCardData struct {
	Common CommonData
	Data   SimpleCardRecord
}

type SimpleCardRecord struct {
	Name       string    `json:"name"`
	Number     string    `json:"number"`
	CVV        uint32    `json:"cvv"`
	ValidUntil time.Time `json:"valid"`
}

type SimpleData struct {
	Common CommonData
	Data   SimpleRecord
}

type SimpleRecord struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type SimpleBinaryData struct {
	Common CommonData
	Data   SimpleBinaryRecord
}

type SimpleBinaryRecord struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type SimpleDataEncrypted struct {
	Common CommonData
	Data   []byte
}
