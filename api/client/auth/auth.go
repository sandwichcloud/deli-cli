package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"context"

	"github.com/sandwichcloud/deli-cli/api"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
)

type AuthClient struct {
	APIServer  *string
	HttpClient *http.Client
}

func (authClient *AuthClient) DiscoverAuth() (api.AuthDiscover, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	authDiscover := api.AuthDiscover{}

	resp, err := ctxhttp.Get(ctx, http.DefaultClient, *authClient.APIServer+"/v1/auth/discover")
	if err != nil {
		if err == context.DeadlineExceeded {
			return authDiscover, api.ErrTimedOut
		}
		return authDiscover, err
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return authDiscover, err
	}

	if resp.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(resp.StatusCode, responseData)
		if err != nil {
			return authDiscover, err
		}
		return authDiscover, apiError
	}

	json.Unmarshal(responseData, &authDiscover)
	return authDiscover, nil
}

func (authClient *AuthClient) ScopeToken(project *api.Project) (*oauth2.Token, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type scopeBody struct {
		ProjectID uuid.UUID `json:"project_id"`
	}

	body := scopeBody{ProjectID: project.ID}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, authClient.HttpClient, *authClient.APIServer+"/v1/auth/tokens/scope", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, api.ErrTimedOut
		}
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	token := &oauth2.Token{}
	json.Unmarshal(responseData, token)
	return token, nil
}

func (authClient *AuthClient) TokenInfo() (*api.TokenInfo, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, authClient.HttpClient, *authClient.APIServer+"/v1/auth/tokens")
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, api.ErrTimedOut
		}
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	tokenInfo := &api.TokenInfo{}
	json.Unmarshal(responseData, tokenInfo)
	return tokenInfo, nil
}

func (authClient *AuthClient) BuiltInLogin(options api.BuiltInAuthDriver, username, password string) (*oauth2.Token, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	data := requestBody{Username: username, Password: password}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("Error parsing Database Auth data into json")
	}

	resp, err := ctxhttp.Post(ctx, http.DefaultClient, *authClient.APIServer+"/v1/auth/builtin/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, api.ErrTimedOut
		}
		return nil, err
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(resp.StatusCode, responseData)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	token := &oauth2.Token{}
	json.Unmarshal(responseData, token)
	return token, nil
}

func (authClient *AuthClient) GithubLogin(options api.GithubAuthDriver, username, password, otpCode string) (*oauth2.Token, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		OTPCode  string `json:"otp_code"`
	}

	data := requestBody{Username: username, Password: password}

	if otpCode != "" {
		data.OTPCode = otpCode
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("Error parsing Github Auth data into json")
	}

	resp, err := ctxhttp.Post(ctx, http.DefaultClient, *authClient.APIServer+"/v1/auth/github/authorization", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, api.ErrTimedOut
		}
		return nil, err
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if _, ok := resp.Header["X-Github-Otp"]; ok {
		if otpCode != "" {
			return nil, ErrOTPInvalid
		}
		return nil, ErrOTPRequired
	}

	if resp.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(resp.StatusCode, responseData)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	token := &oauth2.Token{}
	json.Unmarshal(responseData, token)
	return token, nil
}
