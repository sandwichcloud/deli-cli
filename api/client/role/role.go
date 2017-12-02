package role

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

type RoleClient struct {
	APIServer  *string
	HttpClient *http.Client
}

func (client *RoleClient) Create(name, roleType, description string) (*api.Role, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
	}

	body := createBody{Name: name, Type: roleType, Description: description}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+"/v1/roles", "application/json", bytes.NewBuffer(jsonBody))
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

	role := &api.Role{}
	json.Unmarshal(responseData, role)
	return role, nil
}

func (client *RoleClient) Get(id string) (*api.Role, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, client.HttpClient, *client.APIServer+"/v1/roles/"+id)
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

	role := &api.Role{}
	json.Unmarshal(responseData, role)
	return role, nil
}

func (client *RoleClient) Delete(id string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*client.APIServer + "/v1/roles/" + id)
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

	if response.StatusCode != http.StatusNoContent {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}

	return nil
}

func (client *RoleClient) List(roleType string, limit int, marker string) (*api.RoleList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	parameters.Add("type", roleType)
	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*client.APIServer + "/v1/roles")
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

	if response.StatusCode != http.StatusOK {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return nil, err
		}
		return nil, apiError
	}

	roles := &api.RoleList{}
	json.Unmarshal(responseData, roles)
	return roles, nil
}
