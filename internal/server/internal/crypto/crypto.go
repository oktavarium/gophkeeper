package crypto

import (
	"time"

	"github.com/google/uuid"
)

func GenerateToken() (string, time.Time) {
	id := uuid.New().String()
	validUntil := time.Now().UTC().Add(30 * time.Minute)
	return id, validUntil
}
