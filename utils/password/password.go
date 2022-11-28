package password

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plain string) (string, error) {
	b := []byte(plain)

	hashed, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func ComparePassword(plain string, hashed string) error {
	bPlain := []byte(plain)
	bHashed := []byte(hashed)

	return bcrypt.CompareHashAndPassword(bHashed, bPlain)
}
