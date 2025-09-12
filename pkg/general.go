package pkg

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const stringChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func ComparePassword(hashedPw []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPw, password)
}

func HashPassword(password string) ([]byte, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPw, nil
}

func GenerateAndHashPassword(length int) (string, []byte, error) {
	password, err := GenerateRandomString(length)
	if err != nil {
		return "", nil, err
	}

	hashedPw, err := HashPassword(password)
	if err != nil {
		return "", nil, err
	}

	return password, hashedPw, nil
}

func GenerateRandomString(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("invalid length")
	}

	result := make([]byte, n)
	for i := range n {
		// Pick a random index
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(stringChars))))
		if err != nil {
			return "", err
		}
		result[i] = stringChars[num.Int64()]
	}

	return string(result), nil
}
