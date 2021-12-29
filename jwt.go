package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTHelper struct {
	RSAKey *rsa.PrivateKey
}

func NewJWTHelper(key string) (*JWTHelper, error) {
	var keyBytes []byte
	file, err := os.Open(key)
	if os.IsNotExist(err) {
		keyBytes = []byte(key)
	} else {
		defer file.Close()

		keyBytes, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
	}
	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}
	return &JWTHelper{
		RSAKey: rsaKey,
	}, nil
}

func (j *JWTHelper) NewToken(appID string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
		Issuer:    appID,
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.RSAKey)
}
