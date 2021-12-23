package oauth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type cognitoResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

const (
	TOKEN_ENPOINT = "https://algotrade-user-pool.auth.ap-southeast-1.amazoncognito.com/oauth2/token"
)

var (
	BearerToken = ""
)

func RetrieveAccessToken(ch chan string) {
	for val := range ch {
		if val == "tokenRefreshJob" {
			urlParams := map[string]string{
				"grant_type":    "client_credentials",
				"client_id":     "4ufuh3cobkdbga7cav6ag8jr1h",
				"client_secret": "1knrin05de5uqc1kp82o5st2knmpe2tgvfhiqsnr9fhk6muvg7sc",
				"scope":         "https://user-resource.api.algotrade.dev/api.readwrite",
			}

			param := url.Values{}
			for k, v := range urlParams {
				param.Add(k, v)
			}

			var client http.Client
			var result cognitoResult

			resp, err := client.Post(TOKEN_ENPOINT, "application/x-www-form-urlencoded", strings.NewReader(param.Encode()))
			if err != nil {
				log.Warn(err)
				ch <- "tokenRefreshJob"
			}

			if resp.StatusCode == http.StatusOK {
				err = json.NewDecoder(resp.Body).Decode(&result)
				if err != nil {
					log.Warn(err)
					ch <- "tokenRefreshJob"
				}
				BearerToken = result.AccessToken
				log.Info("Refreshed bearer token from aws cognito!")
			}

			resp.Body.Close()
		}
	}
}
