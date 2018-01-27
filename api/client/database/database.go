package database

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sandwichcloud/deli-cli/api"
	"golang.org/x/net/context/ctxhttp"
)

type DatabaseAuthClient struct {
	APIServer  *string
	HttpClient *http.Client
}

func (client *DatabaseAuthClient) Create(username, password string) (*api.DatabaseUser, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	body := createBody{Username: username, Password: password}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+"/v1/auth/database/users", "application/json", bytes.NewBuffer(jsonBody))
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
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	user := &api.DatabaseUser{}
	json.Unmarshal(responseData, user)
	return user, nil
}

func (client *DatabaseAuthClient) Get(id string) (*api.DatabaseUser, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, client.HttpClient, *client.APIServer+"/v1/auth/database/users/"+id)
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
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	user := &api.DatabaseUser{}
	json.Unmarshal(responseData, user)
	return user, nil
}

func (client *DatabaseAuthClient) Delete(id string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*client.APIServer + "/v1/auth/database/users/" + id)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", Url.String(), nil)
	if err != nil {
		return err
	}
	response, err := ctxhttp.Do(ctx, client.HttpClient, req)
	if err != nil {
		if err == context.DeadlineExceeded {
			return api.ErrTimedOut
		}
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}

	return nil
}

func (client *DatabaseAuthClient) List(limit int, marker string) (*api.DatabaseUserList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*client.APIServer + "/v1/auth/database/users")
	if err != nil {
		return nil, err
	}
	Url.RawQuery = parameters.Encode()

	response, err := ctxhttp.Get(ctx, client.HttpClient, Url.String())
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
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	users := &api.DatabaseUserList{}
	json.Unmarshal(responseData, users)
	return users, nil
}

func (client *DatabaseAuthClient) ChangePassword(id, password string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	var Url *url.URL
	var err error
	if id == "" {
		Url, err = url.Parse(*client.APIServer + "/v1/auth/database/users")
	} else {
		Url, err = url.Parse(*client.APIServer + "/v1/auth/database/users/" + id)
	}
	if err != nil {
		return err
	}

	type passwordBody struct {
		Password string `json:"password"`
	}

	body := passwordBody{Password: password}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("PATCH", Url.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := ctxhttp.Do(ctx, client.HttpClient, req)
	if err != nil {
		if err == context.DeadlineExceeded {
			return api.ErrTimedOut
		}
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}

	return nil
}

func (client *DatabaseAuthClient) UpdateRoles(id string, roles []string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type roleBody struct {
		Roles []string `json:"roles"`
	}

	body := roleBody{Roles: roles}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+"/v1/auth/database/users/"+id+"/roles/update", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		if err == context.DeadlineExceeded {
			return api.ErrTimedOut
		}
		return err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}

	return nil
}
