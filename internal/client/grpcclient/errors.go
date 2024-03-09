package grpcclient

import (
	"errors"
)

var ErrTokenExpired = errors.New("token expired")
