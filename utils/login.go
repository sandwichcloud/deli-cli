package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
	"github.com/sandwichcloud/deli-cli/api/client"
	"github.com/sandwichcloud/deli-cli/api/client/auth"
	"golang.org/x/oauth2"
)

func Login(authClient client.AuthClientInterface, username, password, authMethod string, interactive bool) (*oauth2.Token, error) {
	apiDiscover, err := authClient.DiscoverAuth()
	if err != nil {
		return nil, err
	}
	var token *oauth2.Token
	switch authMethod {
	case "github":
		if apiDiscover.Github == nil {
			return nil, errors.New("Github auth method is not enabled on this API Server.")
		}

		if username == "" {
			return nil, errors.New("Username is required for Github Authentication")
		}

		for interactive && password == "" {
			passwordBytes, err := gopass.GetPasswdPrompt("Please enter your GitHub password: ", true, os.Stdin, os.Stdout)
			if err != nil {
				return nil, err
			}
			password = string(passwordBytes)
		}

		token, err = authClient.GithubLogin(*apiDiscover.Github, username, password, "")
		if interactive && err == auth.ErrOTPRequired {
			var otpBytes []byte
			otpBytes, err = gopass.GetPasswdPrompt("Please enter your GitHub OTP Code: ", true, os.Stdin, os.Stdout)
			if err != nil {
				return nil, err
			}
			otpCode := string(otpBytes)
			token, err = authClient.GithubLogin(*apiDiscover.Github, username, password, otpCode)
		}
	default:
		return nil, errors.New(fmt.Sprintf("Unknown API Auth Driver %s", authMethod))
	}
	if err != nil {
		return nil, err
	}
	return token, nil
}
