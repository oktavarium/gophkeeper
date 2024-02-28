package storage

import (
	"time"
)

type storageModel struct {
	MasterPass       [32]byte                    `json:"master_pass"`
	Login            string                      `json:"login"`
	Password         string                      `json:"password"`
	Token            token                       `json:"token"`
	ServerAddr       string                      `json:"server_addr"`
	SimpleData       map[string]simpleData       `json:"simple_data"`
	SimpleBinaryData map[string]simpleBinaryData `json:"simple_binary_data"`
	SimpleCardData   map[string]simpleCardData   `json:"simple_card_data"`
}

type token struct {
	Id         string    `json:"id"`
	ValidUntil time.Time `json:"valid_until"`
}

type commonData struct {
	IsDeleted bool      `json:"is_deleted"`
	Modified  time.Time `json:"modified"`
}

type simpleData struct {
	Common commonData   `json:"common"`
	Data   simpleRecord `json:"data"`
}

type simpleRecord struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type simpleBinaryData struct {
	Common commonData         `json:"common"`
	Data   simpleBinaryRecord `json:"data"`
}

type simpleBinaryRecord struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type simpleCardData struct {
	Common commonData       `json:"common"`
	Data   simpleCardRecord `json:"data"`
}

type simpleCardRecord struct {
	Name       string    `json:"name"`
	Number     string    `json:"number"`
	CVV        uint32    `json:"ccv"`
	ValidUntil time.Time `json:"valid_until"`
}
