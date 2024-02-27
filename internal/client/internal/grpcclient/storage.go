package grpcclient

import "github.com/oktavarium/gophkeeper/internal/shared/dto"

type storage interface {
	GetToken() (dto.Token, error)
	UpdateToken(dto.Token) error
}
