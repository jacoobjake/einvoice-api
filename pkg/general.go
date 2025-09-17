package pkg

import (
	"crypto/rand"
	"math/big"
	"regexp"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const stringChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func ComparePassword(hashedPw []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPw, password)
}

func HashPassword(password string) ([]byte, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "error encrypting password")
	}
	return hashedPw, nil
}

func GenerateAndHashPassword(length int) (string, []byte, error) {
	password, err := GenerateRandomString(length)
	if err != nil {
		return "", nil, errors.Wrap(err, "error generating random password string")
	}

	hashedPw, err := HashPassword(password)
	if err != nil {
		return "", nil, errors.Wrap(err, "error hashing password")
	}

	return password, hashedPw, nil
}

func GenerateRandomString(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("invalid string length")
	}

	result := make([]byte, n)
	for i := range n {
		// Pick a random index
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(stringChars))))
		if err != nil {
			return "", errors.Wrap(err, "error genering random index")
		}
		result[i] = stringChars[num.Int64()]
	}

	return string(result), nil
}

func IsPasswordValid(pw string) bool {
	// Length check
	if len(pw) < 8 || len(pw) > 32 {
		return false
	}

	// At least one lowercase
	if matched, _ := regexp.MatchString(`[a-z]`, pw); !matched {
		return false
	}

	// At least one uppercase
	if matched, _ := regexp.MatchString(`[A-Z]`, pw); !matched {
		return false
	}

	// At least one digit
	if matched, _ := regexp.MatchString(`[0-9]`, pw); !matched {
		return false
	}

	// At least one special character (anything not letter or digit)
	if matched, _ := regexp.MatchString(`[^A-Za-z0-9]`, pw); !matched {
		return false
	}

	return true
}
