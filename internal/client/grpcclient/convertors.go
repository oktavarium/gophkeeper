package grpcclient

import (
	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func dtoUserInfoToGrpcUserInfo(data models.UserInfo) *pbapi.UserInfo {
	return &pbapi.UserInfo{
		Login:    data.Login,
		Password: data.Password,
	}
}
