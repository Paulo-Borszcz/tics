package glpi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/pauloborszcz/tics/internal/config"
)

type Client struct {
	cfg          *config.Config
	httpClient   *http.Client
	sessionToken string
	mu           sync.Mutex
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{},
	}
}

func (c *Client) InitSession() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	req, err := http.NewRequest("GET", c.cfg.GLPIURL+"/initSession", nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "user_token "+c.cfg.UserToken)
	req.Header.Set("Content-Type", "application/json")
	if c.cfg.AppToken != "" {
		req.Header.Set("App-Token", c.cfg.AppToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("init session request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("init session failed (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		SessionToken string `json:"session_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decoding session response: %w", err)
	}

	c.sessionToken = result.SessionToken
	return nil
}

func (c *Client) KillSession() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sessionToken == "" {
		return
	}

	req, _ := http.NewRequest("GET", c.cfg.GLPIURL+"/killSession", nil)
	c.setHeaders(req)
	c.httpClient.Do(req)
	c.sessionToken = ""
}

func (c *Client) doRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := c.cfg.GLPIURL + endpoint
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.mu.Lock()
	c.setHeaders(req)
	c.mu.Unlock()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Re-auth on 401
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		if err := c.InitSession(); err != nil {
			return nil, fmt.Errorf("re-auth failed: %w", err)
		}
		req, err = http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		c.mu.Lock()
		c.setHeaders(req)
		c.mu.Unlock()
		return c.httpClient.Do(req)
	}

	return resp, nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Session-Token", c.sessionToken)
	if c.cfg.AppToken != "" {
		req.Header.Set("App-Token", c.cfg.AppToken)
	}
}

func (c *Client) get(endpoint string) ([]byte, error) {
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(data))
	}

	return data, nil
}

func (c *Client) post(endpoint string, payload string) ([]byte, error) {
	resp, err := c.doRequest("POST", endpoint, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// BaseURL returns the GLPI base URL (without /apirest.php) for building web links.
func (c *Client) BaseURL() string {
	return strings.TrimSuffix(c.cfg.GLPIURL, "/apirest.php")
}

// SessionUser holds basic user information from a GLPI session.
type SessionUser struct {
	ID   int
	Name string
}

// GetFullSession fetches the current session info and returns the user ID and name.
func (c *Client) GetFullSession() (*SessionUser, error) {
	data, err := c.get("/getFullSession")
	if err != nil {
		return nil, fmt.Errorf("getting full session: %w", err)
	}

	var result struct {
		Session struct {
			GlpiID         int    `json:"glpiID"`
			GlpiFriendName string `json:"glpifriendlyname"`
			GlpiName       string `json:"glpiname"`
		} `json:"session"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decoding session: %w", err)
	}

	name := result.Session.GlpiFriendName
	if name == "" {
		name = result.Session.GlpiName
	}

	return &SessionUser{
		ID:   result.Session.GlpiID,
		Name: name,
	}, nil
}

// GetUserPicture fetches the user's profile picture as raw bytes.
func (c *Client) GetUserPicture(userID int) ([]byte, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/User/%d/Picture", userID), nil)
	if err != nil {
		return nil, fmt.Errorf("fetching user picture: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("user picture not available (status %d)", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading user picture: %w", err)
	}
	return data, nil
}
