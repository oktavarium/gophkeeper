package remote

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
		UserInfo: dtoUserInfoToGrpcUserInfo(data.UserInfo),
		Name:     data.Name,
		Data:     data.Data,
	}
}
