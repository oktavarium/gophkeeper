package grpcserver

import (
	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func grpcUserInfoToDtoUserInfo(data *pbapi.UserInfo) dto.UserInfo {
	return dto.UserInfo{
		Login:    data.GetLogin(),
		Password: data.GetPassword(),
	}
}

func grpcSaveDataToDtoSavaData(data *pbapi.SaveRequest) dto.SaveData {
	return dto.SaveData{
		Name: data.GetName(),
		Data: string(data.GetData()),
	}
}
