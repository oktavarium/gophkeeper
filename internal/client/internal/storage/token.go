package storage

import (
	"fmt"
	"time"
)

func (s *JsonStorage) GetToken() (string, time.Time, error) {
	if !s.isInited() {
		return "", time.Now(), fmt.Errorf("storage is not inited")
	}

	var t token
	s.store.Read(func(data *storageModel) {
		t = data.Token
	})

	return t.Id, t.ValidUntil, nil
}

func (s *JsonStorage) UpdateToken(id string, validUntil time.Time) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	_ = s.store.Write(func(data *storageModel) error {
		data.Token.Id = id
		data.Token.ValidUntil = validUntil
		return nil
	})

	return nil
}
