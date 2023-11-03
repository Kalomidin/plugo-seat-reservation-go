package common

import (
	"context"
	"math/rand"
	"seat-reservation/jwt"

	"github.com/google/uuid"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetUserId(ctx context.Context) (uuid.UUID, error) {
	tokenData, err := jwt.GetJWTToken(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(tokenData.Subject)
}

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		// This will ensure each character in the charset has an equal probability
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
