package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo struct {
	ID       primitive.ObjectID `bson:"_id"`
	Login    string             `bson:"login,omitempty"`
	Password []byte             `bson:"password,omitempty"`
}

type Token struct {
	TokenID    string             `bson:"token_id,omitempty"`
	IP         string             `bson:"ip, omitempty"`
	ValidUntil time.Time          `bson:"valid_until,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
}

type CommonData struct {
	DataID   string    `bson:"data_id,omitempty"`
	DataType int       `bson:"data_type,omitempty"`
	Deleted  bool      `bson:"is_deleted,omitempty"`
	Modified time.Time `bson:"modified,omitempty"`
	Data     []byte    `bson:"data,omitempty"`
}
