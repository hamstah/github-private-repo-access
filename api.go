package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type AppInstallationResponse struct {
	ID      int `json:"id"`
	Account struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"account"`
	RepositorySelection string `json:"repository_selection"`
	AccessTokensURL     string `json:"access_tokens_url"`
	RepositoriesURL     string `json:"repositories_url"`
	HTMLURL             string `json:"html_url"`
	AppID               int    `json:"app_id"`
	AppSlug             string `json:"app_slug"`
	TargetID            int    `json:"target_id"`
	TargetType          string `json:"target_type"`
	Permissions         struct {
		Contents string `json:"contents"`
		Metadata string `json:"metadata"`
	} `json:"permissions"`
	Events                 []interface{} `json:"events"`
	CreatedAt              time.Time     `json:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at"`
	SingleFileName         interface{}   `json:"single_file_name"`
	HasMultipleSingleFiles bool          `json:"has_multiple_single_files"`
	SingleFilePaths        []interface{} `json:"single_file_paths"`
	SuspendedBy            interface{}   `json:"suspended_by"`
	SuspendedAt            interface{}   `json:"suspended_at"`
}

type AccessTokenResponse struct {
	Token       string    `json:"token"`
	ExpiresAt   time.Time `json:"expires_at"`
	Permissions struct {
		Contents string `json:"contents"`
		Metadata string `json:"metadata"`
	} `json:"permissions"`
	RepositorySelection string `json:"repository_selection"`
}

type Client struct {
	Token      string
	HTTPClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) Request(ctx context.Context, method, path string, output interface{}) error {
	url := path
	if !strings.HasPrefix(path, "https://") {
		url = "https://api.github.com/" + path
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		logrus.WithError(err).Error("failed to build request")
		return fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("authorization", "Bearer "+c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("failed to execute request")
		return fmt.Errorf("failed to execute request: %w", err)
	}
	if res.StatusCode >= 300 {
		logrus.WithField("status_code", res.StatusCode).Error("request failed")
		return fmt.Errorf("request failed")
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.WithError(err).Error("failed to read response")
		return fmt.Errorf("failed to read response: %w", err)
	}

	if output != nil {
		err := json.Unmarshal(bodyBytes, output)
		if err != nil {
			logrus.WithError(err).Error("failed to parse body")
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) GetAppInstallations(ctx context.Context) ([]*AppInstallationResponse, error) {
	appInstallations := []*AppInstallationResponse{}

	err := c.Request(ctx, http.MethodGet, "app/installations", &appInstallations)
	if err != nil {
		return nil, err
	}

	return appInstallations, nil
}

func (c *Client) GetAppInstallationByLogin(ctx context.Context, login string) (*AppInstallationResponse, error) {
	appInstallations, err := c.GetAppInstallations(ctx)
	if err != nil {
		return nil, err
	}

	for _, appInstallation := range appInstallations {
		if appInstallation.Account.Login == login {
			return appInstallation, nil
		}
	}
	return nil, nil
}
