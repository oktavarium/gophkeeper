package storage

import (
	"fmt"
	"time"
)

func (s *Storage) GetToken() (string, time.Time, error) {
	if !s.isInited() {
		return "", time.Now().UTC(), fmt.Errorf("storage is not inited")
	}

	var t token
	s.store.Read(func(data *storageModel) {
		t = data.Token
	})

	return t.ID, t.ValidUntil, nil
}

func (s *Storage) UpdateToken(id string, validUntil time.Time) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	_ = s.store.Write(func(data *storageModel) error {
		data.Token.ID = id
		data.Token.ValidUntil = validUntil
		return nil
	})

	return nil
}
