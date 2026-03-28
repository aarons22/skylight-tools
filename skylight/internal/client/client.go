package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultBaseURL = "https://api.ourskylight.com/api"

// Client is an HTTP client for the Skylight Calendar API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new API client with the given base URL and auth token.
// The token should be the base64-encoded "user_id:user_token" string.
func NewClient(baseURL, token string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		httpClient: &http.Client{},
	}
}

// Do executes an HTTP request and returns the raw response body.
// pathParams are substituted into the path using {key} placeholders.
// queryParams are appended as URL query parameters.
// body is JSON-encoded for POST/PUT/PATCH requests when non-nil.
func (c *Client) Do(method, path string, pathParams, queryParams map[string]string, body interface{}) ([]byte, error) {
	// Substitute path parameters
	for k, v := range pathParams {
		path = strings.ReplaceAll(path, "{"+k+"}", url.PathEscape(v))
	}

	fullURL := c.baseURL + path

	// Append query parameters
	if len(queryParams) > 0 {
		params := url.Values{}
		for k, v := range queryParams {
			if v != "" {
				params.Set(k, v)
			}
		}
		if encoded := params.Encode(); encoded != "" {
			fullURL += "?" + encoded
		}
	}

	// Encode body
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encoding request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.token))
	req.Header.Set("User-Agent", "SkylightMobile (web)")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	return respBody, nil
}

// LoginResponse is the response from POST /sessions.
type LoginResponse struct {
	UserID    int    `json:"user_id"`
	UserToken string `json:"user_token"`
}

// Login authenticates with email and password, returning a base64 token string.
func Login(baseURL, email, password string) (string, error) {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", baseURL+"/sessions", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("creating login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SkylightMobile (web)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("login failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var lr LoginResponse
	if err := json.Unmarshal(body, &lr); err != nil {
		return "", fmt.Errorf("parsing login response: %w", err)
	}
	if lr.UserID == 0 || lr.UserToken == "" {
		return "", fmt.Errorf("login response missing user_id or user_token")
	}

	raw := fmt.Sprintf("%d:%s", lr.UserID, lr.UserToken)
	return base64.StdEncoding.EncodeToString([]byte(raw)), nil
}

// Config represents the skylight CLI configuration file.
type Config struct {
	Token   string `yaml:"token"`
	FrameID string `yaml:"frame_id,omitempty"`
}

// ConfigPath returns the default config file path.
func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "skylight", "config.yaml")
}

// LoadConfig reads the config file at path (uses ConfigPath() if empty).
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = ConfigPath()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// SaveConfig writes the config to path (uses ConfigPath() if empty).
func SaveConfig(path string, cfg *Config) error {
	if path == "" {
		path = ConfigPath()
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// ResolveToken returns the token from the SKYLIGHT_TOKEN env var or config file.
func ResolveToken(configPath string) (string, error) {
	if t := os.Getenv("SKYLIGHT_TOKEN"); t != "" {
		return t, nil
	}
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return "", err
	}
	if cfg.Token == "" {
		return "", fmt.Errorf("no auth token found; run 'skylight account login' or set SKYLIGHT_TOKEN")
	}
	return cfg.Token, nil
}

// ResolveFrameID returns the frame ID from flag value, SKYLIGHT_FRAME_ID env var, or config file.
func ResolveFrameID(flagValue, configPath string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}
	if f := os.Getenv("SKYLIGHT_FRAME_ID"); f != "" {
		return f, nil
	}
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return "", err
	}
	if cfg.FrameID == "" {
		return "", fmt.Errorf("no frame ID found; use --frame-id flag or set SKYLIGHT_FRAME_ID")
	}
	return cfg.FrameID, nil
}
