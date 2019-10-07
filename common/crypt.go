package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Hash method
func Hash(password string, salt int) (string, error) {

	_pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(_pwd, salt)
	if err != nil {
		return "", err
	}
	return string(hash), nil

}

// PasswordCompare method
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

// RsaConfigSetup method
func RsaConfigSetup(rsaPrivateKeyLocation, rsaPrivateKeyPassword string, rsaPublicKeyLocation string) (*rsa.PrivateKey, error) {
	if rsaPrivateKeyLocation == "" {

		return GenRSA(1024)
	}

	priv, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {

		privkey, err := GenRSA(1024)
		if err != nil {
			return nil, err
		}

		privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
		privkeyPem := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: privkeyBytes,
			},
		)

		ioutil.WriteFile(rsaPrivateKeyLocation, privkeyPem, 0644)
		return privkey, nil
	}

	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte
	if privPem.Type != "RSA PRIVATE KEY" {

	}

	if rsaPrivateKeyPassword != "" {
		privPemBytes, err = x509.DecryptPEMBlock(privPem, []byte(rsaPrivateKeyPassword))
	} else {
		privPemBytes = privPem.Bytes
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`

			return GenRSA(1024)
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {

		return GenRSA(1024)
	}

	pub, err := ioutil.ReadFile(rsaPublicKeyLocation)
	if err != nil {

		return GenRSA(1024)
	}
	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {

		return GenRSA(1024)
	}
	if pubPem.Type != "RSA PUBLIC KEY" {

		return GenRSA(1024)
	}

	if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {

		return GenRSA(1024)
	}

	var pubKey *rsa.PublicKey
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {

		return GenRSA(1024)
	}

	privateKey.PublicKey = *pubKey

	return privateKey, nil
}

// GenRSA method
func GenRSA(bits int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)

	return key, err
}
