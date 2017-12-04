package instance

import (
	"net/http"

	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"

	"net/url"

	"strconv"

	"fmt"

	"github.com/sandwichcloud/deli-cli/api"
	"golang.org/x/net/context/ctxhttp"
)

type InstanceClient struct {
	APIServer  *string
	HttpClient *http.Client
}

func (client *InstanceClient) Create(name, imageID, networkID string, keypairIDs []string, tags map[string]string) (*api.Instance, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name       string            `json:"name"`
		ImageID    string            `json:"image_id"`
		NetworkID  string            `json:"network_id"`
		KeypairIDs []string          `json:"keypair_ids,omitempty"`
		Tags       map[string]string `json:"tags"`
	}

	body := createBody{Name: name, ImageID: imageID, NetworkID: networkID, KeypairIDs: keypairIDs, Tags: tags}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+"/v1/instances", "application/json", bytes.NewBuffer(jsonBody))
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

	instance := &api.Instance{}
	json.Unmarshal(responseData, instance)
	return instance, nil
}

func (client *InstanceClient) Get(id string) (*api.Instance, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, client.HttpClient, *client.APIServer+"/v1/instances/"+id)
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

	instance := &api.Instance{}
	json.Unmarshal(responseData, instance)
	return instance, nil
}

func (client *InstanceClient) Delete(id string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*client.APIServer + "/v1/instances/" + id)
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

	if response.StatusCode != http.StatusAccepted {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}

	return nil
}

func (client *InstanceClient) List(imageID string, limit int, marker string) (*api.InstanceList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	if len(imageID) > 0 {
		parameters.Add("image_id", imageID)
	}

	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*client.APIServer + "/v1/instances")
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

	instances := &api.InstanceList{}
	json.Unmarshal(responseData, instances)
	return instances, nil
}

func (client *InstanceClient) ActionStop(id string, hard bool, timeout int) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type stopBody struct {
		Hard    bool `json:"hard"`
		Timeout int  `json:"timeout"`
	}

	body := stopBody{Hard: hard, Timeout: timeout}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, *client.APIServer+fmt.Sprintf("/v1/instances/%s/action/stop", id), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
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

	if response.StatusCode != http.StatusAccepted {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}
	return nil
}

func (client *InstanceClient) ActionStart(id string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	req, err := http.NewRequest(http.MethodPut, *client.APIServer+fmt.Sprintf("/v1/instances/%s/action/start", id), nil)
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

	if response.StatusCode != http.StatusAccepted {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}
	return nil
}

func (client *InstanceClient) ActionRestart(id string, hard bool, timeout int) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type restartBody struct {
		Hard    bool `json:"hard"`
		Timeout int  `json:"timeout"`
	}

	body := restartBody{Hard: hard, Timeout: timeout}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, *client.APIServer+fmt.Sprintf("/v1/instances/%s/action/restart", id), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
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

	if response.StatusCode != http.StatusAccepted {
		apiError, err := api.ParseErrors(response.StatusCode, responseData)
		if err != nil {
			return err
		}
		return apiError
	}
	return nil
}

func (client *InstanceClient) ActionImage(id string, name string, visibility string) (*api.Image, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name       string `json:"name"`
		Visibility string `json:"visibility"`
	}

	body := createBody{Name: name, Visibility: visibility}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+fmt.Sprintf("/v1/instances/%s/action/image", id), "application/json", bytes.NewBuffer(jsonBody))
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

	image := &api.Image{}
	json.Unmarshal(responseData, image)
	return image, nil
}

func (client *InstanceClient) ActionResetState(id string, active bool) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type resetStateBody struct {
		Active bool `json:"active"`
	}

	body := resetStateBody{Active: active}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, *client.APIServer+fmt.Sprintf("/v1/instances/%s/action/reset_state", id), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
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
