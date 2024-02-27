package storage

import (
	"time"
)

type storageModel struct {
	MasterPass       [32]byte                      `json:"master_pass"`
	Login            string                        `json:"login"`
	Password         string                        `json:"password"`
	Token            token                         `json:"token"`
	ServerAddr       string                        `json:"server_addr"`
	SimpleData       map[string]simpleRecord       `json:"simple_data"`
	SimpleBinaryData map[string]simpleBinaryRecord `json:"simple_binary_data"`
	SimpleCardData   map[string]simpleCardRecord   `json:"simple_card_data"`
}

type token struct {
	Id         string    `json:"id"`
	ValidUntil time.Time `json:"valid_until"`
}

type commonData struct {
	Name      string    `json:"name"`
	IsDeleted bool      `json:"is_deleted"`
	Modified  time.Time `json:"modified"`
}

type simpleRecord struct {
	Common commonData `json:"common"`
	Data   string     `json:"data"`
}

type simpleBinaryRecord struct {
	Common commonData `json:"common"`
	Data   []byte     `json:"data"`
}

type simpleCardRecord struct {
	Common     commonData `json:"common"`
	Number     string     `json:"number"`
	CVV        uint32     `json:"ccv"`
	ValidUntil time.Time  `json:"valid_until"`
}
