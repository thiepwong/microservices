package common

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string, salt int) (string, error) {

	_pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(_pwd, salt)
	if err != nil {
		return "", err
	}
	return string(hash), nil

}

func PasswordCompare(rawPassword string, passwordHash string, salt int) bool {
	if rawPassword == "" || passwordHash == "" {
		return false
	}

	byteHash := []byte(passwordHash)
	_pwd := []byte(rawPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, _pwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
