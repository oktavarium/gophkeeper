package crypto

import (
	"time"

	"github.com/google/uuid"
)

const tokenValidPeriod = 30 * time.Minute

func GenerateToken() (string, time.Time) {
	id := uuid.New().String()
	validUntil := time.Now().UTC().Add(tokenValidPeriod)
	return id, validUntil
}
