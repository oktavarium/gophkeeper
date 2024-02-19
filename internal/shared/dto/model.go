package dto

type UserInfo struct {
	Login    string
	Password string
}

type SaveData struct {
	UserInfo UserInfo
	Name     string
	Data     string
}
