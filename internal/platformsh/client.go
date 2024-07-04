package platformsh

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	restyClient *resty.Client
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Project struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
}

type Environment struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Title          string `json:"title"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	DefaultDomain  string `json:"default_domain"`
	EnableSMTP     bool   `json:"enable_smtp"`
	RestrictRobots bool   `json:"restrict_robots"`
}

type EnvironmentsResponse struct {
	Environments []Environment `json:"environments"`
}

func NewClient(apiToken string) (*Client, error) {
	client := resty.New()
	tokenResp, err := client.R().
		SetBasicAuth("platform-api-user", "").
		SetFormData(map[string]string{
			"grant_type": "api_token",
			"api_token":  apiToken,
		}).
		SetResult(&TokenResponse{}).
		Post("https://auth.api.platform.sh/oauth2/token")

	if err != nil {
		return nil, err
	}

	tokenResponse := tokenResp.Result().(*TokenResponse)
	client.SetAuthToken(tokenResponse.AccessToken)

	return &Client{restyClient: client}, nil
}

func (c *Client) GetSession() *resty.Client {
	return c.restyClient
}

func (c *Client) GetProjects() ([]Project, error) {
	var projectsResponse ProjectsResponse
	_, err := c.restyClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.restyClient.Token)).
		SetResult(&projectsResponse).
		Get("https://api.platform.sh/projects")

	if err != nil {
		return nil, err
	}

	return projectsResponse.Projects, nil
}

func (c *Client) GetEnvironments(projectID string) ([]Environment, error) {
	var environmentsResponse EnvironmentsResponse
	_, err := c.restyClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.restyClient.Token)).
		SetResult(&environmentsResponse).
		Get(fmt.Sprintf("https://api.platform.sh/projects/%s/environments", projectID))

	if err != nil {
		return nil, err
	}

	return environmentsResponse.Environments, nil
}

func (c *Client) GetEnvironment(projectID, environmentID string) (*Environment, error) {
	var environment Environment
	_, err := c.restyClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.restyClient.Token)).
		SetResult(&environment).
		Get(fmt.Sprintf("https://api.platform.sh/projects/%s/environments/%s", projectID, environmentID))

	if err != nil {
		return nil, err
	}

	return &environment, nil
}

func (c *Client) CreateEnvironment(projectID string, environment *Environment) (*Environment, error) {
	var createdEnvironment Environment
	_, err := c.restyClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.restyClient.Token)).
		SetBody(environment).
		SetResult(&createdEnvironment).
		Post(fmt.Sprintf("https://api.platform.sh/projects/%s/environments", projectID))

	if err != nil {
		return nil, err
	}

	return &createdEnvironment, nil
}

func (c *Client) UpdateEnvironment(projectID, environmentID string, environment *Environment) (*Environment, error) {
	var updatedEnvironment Environment
	_, err := c.restyClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.restyClient.Token)).
		SetBody(environment).
		SetResult(&updatedEnvironment).
		Patch(fmt.Sprintf("https://api.platform.sh/projects/%s/environments/%s", projectID, environmentID))

	if err != nil {
		return nil, err
	}

	return &updatedEnvironment, nil
}

func (c *Client) DeleteEnvironment(projectID, environmentID string) error {
	_, err := c.restyClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.restyClient.Token)).
		Delete(fmt.Sprintf("https://api.platform.sh/projects/%s/environments/%s", projectID, environmentID))

	return err
}
