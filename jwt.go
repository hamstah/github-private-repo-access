package main

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTHelper struct {
	RSAKey *rsa.PrivateKey
}

func NewJWTHelper(key string) (*JWTHelper, error) {
	var keyBytes []byte
	file, err := os.Open(key)

	if err == nil {
		defer file.Close()

		keyBytes, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("can't read key file: %w", err)
		}
	} else {
		if strings.HasPrefix(key, "LS0t") {
			keyBytes, err = base64.StdEncoding.DecodeString(key)
			if err != nil {
				return nil, fmt.Errorf("can't decode base64 string: %w", err)
			}
		} else {
			keyBytes = []byte(key)
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
