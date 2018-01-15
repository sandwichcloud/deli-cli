package image

import (
	"net/http"

	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"

	"net/url"
	"strconv"

	"github.com/sandwichcloud/deli-cli/api"
	"golang.org/x/net/context/ctxhttp"
)

type ImageClient struct {
	APIServer  *string
	HttpClient *http.Client
}

func (client *ImageClient) Create(name, regionID, fileName, visibility string) (*api.Image, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name       string `json:"name"`
		RegionID   string `json:"region_id"`
		FileName   string `json:"file_name"`
		Visibility string `json:"visibility"`
	}

	body := createBody{Name: name, RegionID: regionID, FileName: fileName, Visibility: visibility}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+"/v1/images", "application/json", bytes.NewBuffer(jsonBody))
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

func (client *ImageClient) Get(id string) (*api.Image, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, client.HttpClient, *client.APIServer+"/v1/images/"+id)
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

func (client *ImageClient) Delete(id string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*client.APIServer + "/v1/images/" + id)
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

func (client *ImageClient) List(visibility string, limit int, marker string) (*api.ImageList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	parameters.Add("visibility", visibility)
	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*client.APIServer + "/v1/images")
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

	images := &api.ImageList{}
	json.Unmarshal(responseData, images)
	return images, nil
}

func (client *ImageClient) ActionLock(id string) error {
	return nil
}

func (client *ImageClient) ActionUnlock(id string) error {
	return nil
}

func (client *ImageClient) MemberAdd(id, projectID string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type memberAddBody struct {
		ProjectID string `json:"project_id"`
	}

	body := memberAddBody{ProjectID: projectID}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, client.HttpClient, *client.APIServer+"/v1/images/"+id+"/members", "application/json", bytes.NewBuffer(jsonBody))
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

func (client *ImageClient) MemberList(id string) (*api.ImageMemberList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	Url, err := url.Parse(*client.APIServer + "/v1/images/" + id + "/members")
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

	imageMembers := &api.ImageMemberList{}
	json.Unmarshal(responseData, imageMembers)
	return imageMembers, nil
}

func (client *ImageClient) MemberRemove(id, projectID string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*client.APIServer + "/v1/images/" + id + "/members/" + projectID)
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
