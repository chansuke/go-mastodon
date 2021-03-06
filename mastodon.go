package mastodon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Config is a setting for access mastodon APIs.
type Config struct {
	Server       string
	ClientID     string
	ClientSecret string
	AccessToken  string
}

// Client is a API client for mastodon.
type Client struct {
	http.Client
	config *Config
}

func (c *Client) doAPI(method string, uri string, params url.Values, res interface{}) error {
	url, err := url.Parse(c.config.Server)
	if err != nil {
		return err
	}
	url.Path = path.Join(url.Path, uri)

	var resp *http.Response
	req, err := http.NewRequest(method, url.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	resp, err = c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if res == nil {
		return nil
	}

	if method == http.MethodGet && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad request: %v", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(&res)
}

// NewClient return new mastodon API client.
func NewClient(config *Config) *Client {
	return &Client{
		Client: *http.DefaultClient,
		config: config,
	}
}

// Authenticate get access-token to the API.
func (c *Client) Authenticate(username, password string) error {
	params := url.Values{}
	params.Set("client_id", c.config.ClientID)
	params.Set("client_secret", c.config.ClientSecret)
	params.Set("grant_type", "password")
	params.Set("username", username)
	params.Set("password", password)
	params.Set("scope", "read write follow")

	url, err := url.Parse(c.config.Server)
	if err != nil {
		return err
	}
	url.Path = path.Join(url.Path, "/oauth/token")

	req, err := http.NewRequest(http.MethodPost, url.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad authorization: %v", resp.Status)
	}

	res := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	c.config.AccessToken = res.AccessToken
	return nil
}

// Toot is struct to post status.
type Toot struct {
	Status      string  `json:"status"`
	InReplyToID int64   `json:"in_reply_to_id"`
	MediaIDs    []int64 `json:"media_ids"`
	Sensitive   bool    `json:"sensitive"`
	SpoilerText string  `json:"spoiler_text"`
	Visibility  string  `json:"visibility"`
}

// Mention hold information for mention.
type Mention struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
	ID       int64  `json:"id"`
}

// Tag hold information for tag.
type Tag struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Attachment hold information for attachment.
type Attachment struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`
	URL        string `json:"url"`
	RemoteURL  string `json:"remote_url"`
	PreviewURL string `json:"preview_url"`
	TextURL    string `json:"text_url"`
}
