package memory

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type Token struct {
	name       string
	validUntil int64
}

type loginData struct {
	password string
	data     map[string]string
	token    Token
}

type Storage struct {
	loginsData map[string]loginData
}

func NewStorage() *Storage {
	return &Storage{
		loginsData: make(map[string]loginData),
	}
}

func (s *Storage) Register(ctx context.Context, in dto.UserInfo) error {
	if len(in.Login) == 0 || len(in.Password) == 0 {
		return fmt.Errorf("empty credentials")
	}

	_, ok := s.loginsData[in.Login]
	if ok {
		return fmt.Errorf("already registerd")
	}

	newLogin := loginData{
		password: in.Password,
		data:     make(map[string]string),
		token: Token{
			name:       "token name",
			validUntil: time.Now().Add(time.Hour).Unix(),
		},
	}

	s.loginsData[in.Login] = newLogin

	return nil
}

func (s *Storage) Login(ctx context.Context, in dto.UserInfo) error {
	if len(in.Login) == 0 || len(in.Password) == 0 {
		return fmt.Errorf("empty credentials")
	}

	loginData, ok := s.loginsData[in.Login]
	if !ok {
		return fmt.Errorf("not registerd")
	}

	loginData.token.validUntil = time.Now().Add(time.Hour).Unix()
	s.loginsData[in.Login] = loginData

	return nil
}

func (s *Storage) Save(ctx context.Context, in dto.SaveData) error {
	log.Println(in)
	return nil
}
