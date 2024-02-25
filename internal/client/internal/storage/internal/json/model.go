package json

import (
	"time"
)

type storageModel struct {
	MasterPass       [32]byte                      `json:"master_pass"`
	Login            []byte                        `json:"login"`
	Password         []byte                        `json:"password"`
	ServerAddr       []byte                        `json:"server_addr"`
	SimpleData       map[string]simpleRecord       `json:"simple_data"`
	SimpleBinaryData map[string]simpleBinaryRecord `json:"simple_binary_data"`
	SimpleCardData   map[string]simpleCardRecord   `json:"simple_card_data"`
}

type simpleRecord struct {
	Data string `json:"data"`
}

type simpleBinaryRecord struct {
	Data []byte `json:"data"`
}

type simpleCardRecord struct {
	Number     string    `json:"number"`
	Ccv        uint32    `json:"ccv"`
	ValidUntil time.Time `json:"valid_until"`
}
