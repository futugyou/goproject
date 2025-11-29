package security

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateClientCredentials() (string, string, error) {
	clientID, err := generateRandomString(24)
	if err != nil {
		return "", "", err
	}

	clientSecret, err := generateRandomString(48)
	if err != nil {
		return "", "", err
	}

	return clientID, clientSecret, nil
}
