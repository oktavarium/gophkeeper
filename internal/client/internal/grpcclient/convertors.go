package grpcclient

import (
	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func dtoUserInfoToGrpcUserInfo(data dto.UserInfo) *pbapi.UserInfo {
	return &pbapi.UserInfo{
		Login:    data.Login,
		Password: data.Password,
	}
}

func dtoSavaDataToGrpcSaveData(data dto.SaveData) *pbapi.SaveRequest {
	return &pbapi.SaveRequest{
		Name: data.Name,
		Data: []byte(data.Data),
	}
}
