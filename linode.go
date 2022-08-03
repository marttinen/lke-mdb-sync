package main

import (
	"net/http"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

func linodeClient(linodeToken string, debug bool) linodego.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: linodeToken})
	oauth2Client := &http.Client{Transport: &oauth2.Transport{Source: tokenSource}}
	client := linodego.NewClient(oauth2Client)
	client.SetDebug(debug)
	return client
}
