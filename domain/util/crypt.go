package util

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

// HashCost defines the cost of the bcrypt hash.
const HashCost = 12

// Hash generates a hashed string using the bcrypt key derivation function.
func Hash(informationToHash string) (string, error) {
	generated, err := bcrypt.GenerateFromPassword([]byte(informationToHash), HashCost)
	if err != nil {
		return "", derr.JoinInternalError(err, "failed to hash password")
	}

	return string(generated), nil
}

// GetUserIDFromToken extracts the user ID from a PASETO token string
func GetUserIDFromToken(token string, pasetoSecurityKey string) (int, bool, error) {
	symmetricKey := []byte(pasetoSecurityKey)
	now := time.Now()

	var payload paseto.JSONToken
	var footer string
	_, err := paseto.Parse(token, &payload, &footer, symmetricKey, nil)
	if err != nil {
		return 0, false, derr.JoinInternalError(err, "failed to parse token")
	}

	if now.After(payload.Expiration) {
		return 0, true, nil
	}

	var userID int
	err = json.Unmarshal([]byte(payload.Subject), &userID)
	if err != nil {
		return 0, false, derr.JoinInternalError(err, "failed to parse unmarshal token payload")
	}

	return userID, false, nil
}

// GetNewAuthToken generates a PASETO token for the provided user
func GetNewAuthToken(userID int64, pasetoSecurityKey string) (string, error) {
	v2 := paseto.NewV2()
	now := time.Now()
	expiration := now.Add(7 * 24 * time.Hour).UTC()
	symmetricKey := []byte(pasetoSecurityKey)

	tokenUUID, err := uuid.NewRandom()
	if err != nil {
		return "", derr.JoinInternalError(err, "failed to generate new UUID")
	}

	subjectJS, err := json.Marshal(userID)
	if err != nil {
		return "", derr.JoinInternalError(err, "failed to marshal")
	}

	// Filling a new token with relevant data
	jsonToken := paseto.JSONToken{
		Issuer:     "heart-shaped-box/api",
		Jti:        tokenUUID.String(),
		Subject:    string(subjectJS),
		Expiration: expiration,
		IssuedAt:   now,
		NotBefore:  now,
	}

	footer := struct {
		ExpiresAt time.Time `json:"expires_at"`
	}{
		ExpiresAt: expiration,
	}

	encrypted, err := v2.Encrypt(symmetricKey, jsonToken, footer)
	if err != nil {
		return "", derr.JoinInternalError(err, "failed to encrypt json token")
	}

	return encrypted, nil
}

func GenerateRandomPassword() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	buffer := make([]byte, 16)
	for i := range buffer {
		buffer[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(buffer)
}

// CheckValidPassword verifies if the provided input password (after hashing it) matches the provided hashed password
func CheckValidPassword(input, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(input))
	return err == nil
}
