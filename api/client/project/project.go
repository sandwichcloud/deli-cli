package project

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"context"

	"github.com/sandwichcloud/deli-cli/api"
	"golang.org/x/net/context/ctxhttp"
)

type ProjectClient struct {
	APIServer  *string
	HttpClient *http.Client
}

func (projectClient *ProjectClient) Create(name string) (*api.Project, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type createBody struct {
		Name string `json:"name"`
	}

	body := createBody{Name: name}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, projectClient.HttpClient, *projectClient.APIServer+"/v1/projects", "application/json", bytes.NewBuffer(jsonBody))
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

	project := &api.Project{}
	json.Unmarshal(responseData, project)
	return project, nil
}

func (projectClient *ProjectClient) Get(id string) (*api.Project, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, projectClient.HttpClient, *projectClient.APIServer+"/v1/projects/"+id)
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

	project := &api.Project{}
	json.Unmarshal(responseData, project)
	return project, nil
}

func (projectClient *ProjectClient) Delete(id string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*projectClient.APIServer + "/v1/projects/" + id)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", Url.String(), nil)
	if err != nil {
		return err
	}
	response, err := ctxhttp.Do(ctx, projectClient.HttpClient, req)
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

func (projectClient *ProjectClient) List(all bool, limit int, marker string) (*api.ProjectList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	parameters.Add("all", strconv.FormatBool(all))
	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*projectClient.APIServer + "/v1/projects")
	if err != nil {
		return nil, err
	}
	Url.RawQuery = parameters.Encode()

	response, err := ctxhttp.Get(ctx, projectClient.HttpClient, Url.String())
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

	projects := &api.ProjectList{}
	json.Unmarshal(responseData, projects)
	return projects, nil
}

func (projectClient *ProjectClient) GetQuota() (*api.ProjectQuota, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, projectClient.HttpClient, *projectClient.APIServer+"/v1/project-quota")
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

	projectQuota := &api.ProjectQuota{}
	json.Unmarshal(responseData, projectQuota)
	return projectQuota, nil
}

func (projectClient *ProjectClient) SetQuota(vcpu, ram, disk int) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type quotaBody struct {
		VCPU int `json:"vcpu"`
		Ram  int `json:"ram"`
		Disk int `json:"disk"`
	}

	body := quotaBody{VCPU: vcpu, Ram: ram, Disk: disk}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, projectClient.HttpClient, *projectClient.APIServer+"/v1/project-quota", "application/json", bytes.NewBuffer(jsonBody))
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

func (projectClient *ProjectClient) AddMember(username, driver string) (*api.ProjectMember, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type addBody struct {
		Username string `json:"username"`
		Driver   string `json:"driver"`
	}

	body := addBody{
		Username: username,
		Driver:   driver,
	}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, projectClient.HttpClient, *projectClient.APIServer+"/v1/project-members", "application/json", bytes.NewBuffer(jsonBody))
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

	projectMember := &api.ProjectMember{}
	json.Unmarshal(responseData, projectMember)
	return projectMember, nil
}

func (projectClient *ProjectClient) GetMember(id string) (*api.ProjectMember, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	response, err := ctxhttp.Get(ctx, projectClient.HttpClient, *projectClient.APIServer+"/v1/project-members/"+id)
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

	projectMember := &api.ProjectMember{}
	json.Unmarshal(responseData, projectMember)
	return projectMember, nil
}

func (projectClient *ProjectClient) ListMembers(limit int, marker string) (*api.ProjectMemberList, error) {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	parameters := url.Values{}

	parameters.Add("limit", strconv.FormatInt(int64(limit), 10))

	if len(marker) > 0 {
		parameters.Add("marker", marker)
	}

	Url, err := url.Parse(*projectClient.APIServer + "/v1/project-members")
	if err != nil {
		return nil, err
	}
	Url.RawQuery = parameters.Encode()

	response, err := ctxhttp.Get(ctx, projectClient.HttpClient, Url.String())
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

	projectMemberList := &api.ProjectMemberList{}
	json.Unmarshal(responseData, projectMemberList)
	return projectMemberList, nil
}

func (projectClient *ProjectClient) UpdateMember(id string, roles []string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()

	type updateBody struct {
		Roles []string `json:"roles"`
	}

	body := updateBody{
		Roles: roles,
	}
	jsonBody, _ := json.Marshal(body)

	response, err := ctxhttp.Post(ctx, projectClient.HttpClient, *projectClient.APIServer+"/v1/project-members/"+id, "application/json", bytes.NewBuffer(jsonBody))
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

func (projectClient *ProjectClient) RemoveMember(id string) error {
	ctx, cancel := api.CreateTimeoutContext()
	defer cancel()
	Url, err := url.Parse(*projectClient.APIServer + "/v1/project-members/" + id)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", Url.String(), nil)
	if err != nil {
		return err
	}
	response, err := ctxhttp.Do(ctx, projectClient.HttpClient, req)
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
