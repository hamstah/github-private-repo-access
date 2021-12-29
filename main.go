package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	appID         = kingpin.Flag("app-id", "The Github App ID").Required().String()
	appPrivateKey = kingpin.Flag("app-private-key", "Filename or private key in PEM format.").Required().String()
	name          = kingpin.Flag("name", "Username or organizations name.").Required().String()
	setGoPrivate  = kingpin.Flag("go-private", "Set the GOPRIVATE environment variable").Bool()
)

func handleError(err error, message string) {
	if err != nil {
		logrus.WithError(err).Fatalln(message)
	}
}

func main() {
	kingpin.Parse()
	ctx := context.Background()

	jwt, err := NewJWTHelper(*appPrivateKey)
	handleError(err, "failed to initialise helper")

	token, err := jwt.NewToken(*appID)
	handleError(err, "failed to generate token")

	client := NewClient(token)

	appInstallation, err := client.GetAppInstallationByLogin(ctx, *name)
	handleError(err, "failed to generate JWT")

	if appInstallation == nil {
		logrus.WithField("name", *name).Fatalln("can't find an app installation for that name")
	}

	accessToken := AccessTokenResponse{}
	err = client.Request(ctx, http.MethodPost, appInstallation.AccessTokensURL, &accessToken)
	handleError(err, "failed to generate access token")

	fmt.Printf("git config --global url.\"https://x-access-token:%s@github.com/%s\".insteadOf \"https://github.com/%s\"\n", accessToken.Token, *name, *name) // nolint
	if *setGoPrivate {
		fmt.Printf("go env -w GOPRIVATE=github.com/%s/*,$(go env GOPRIVATE)\n", *name) // nolint
	}
}
