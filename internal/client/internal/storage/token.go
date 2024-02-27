package storage

import (
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func (s *JsonStorage) GetToken() (dto.Token, error) {
	if !s.isInited() {
		return dto.Token{}, fmt.Errorf("storage is not inited")
	}

	var t token
	s.store.Read(func(data *storageModel) {
		t = data.Token
	})

	return dto.Token{
		Id:         t.Id,
		ValidUntil: t.ValidUntil,
	}, nil
}

func (s *JsonStorage) UpdateToken(token dto.Token) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	_ = s.store.Write(func(data *storageModel) error {
		data.Token.Id = token.Id
		data.Token.ValidUntil = token.ValidUntil
		return nil
	})

	return nil
}
