package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultBaseURL = "https://app.ourskylight.com/api"
const oauthURL = "https://app.ourskylight.com/oauth/token"
const oauthClientID = "skylight-mobile"

// Client is an HTTP client for the Skylight Calendar API.
type Client struct {
	baseURL    string
	token      string // OAuth Bearer access_token
	httpClient *http.Client
}

// NewClient creates a new API client with the given base URL and OAuth Bearer token.
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
	for k, v := range pathParams {
		path = strings.ReplaceAll(path, "{"+k+"}", url.PathEscape(v))
	}

	fullURL := c.baseURL + path

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

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
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

// OAuthTokenResponse is the response from POST /oauth/token.
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
}

var csrfRe = regexp.MustCompile(`<meta name="csrf-token" content="([^"]+)"`)

// Login authenticates with email and password using the OAuth Authorization Code flow:
//  1. GET /auth/session/new  — extract CSRF token from HTML
//  2. POST /auth/session     — submit credentials; server sets session cookie
//  3. GET /oauth/authorize   — exchange session for an auth code (captures redirect)
//  4. POST /oauth/token      — exchange auth code for Bearer + refresh tokens
func Login(email, password string) (*OAuthTokenResponse, error) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		// Do not follow the final redirect to ourskylight.com — we need to capture
		// the authorization code from the Location header.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if strings.HasPrefix(req.URL.String(), "https://ourskylight.com") {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// Step 1: fetch login page for CSRF token
	resp, err := client.Get("https://app.ourskylight.com/auth/session/new")
	if err != nil {
		return nil, fmt.Errorf("fetching login page: %w", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	matches := csrfRe.FindSubmatch(body)
	if len(matches) < 2 {
		return nil, errors.New("could not extract CSRF token from login page")
	}
	csrfToken := string(matches[1])

	// Step 2: submit credentials
	form := url.Values{}
	form.Set("authenticity_token", csrfToken)
	form.Set("email", email)
	form.Set("password", password)

	req, _ := http.NewRequest("POST", "https://app.ourskylight.com/auth/session", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "SkylightMobile (web)")
	req.Header.Set("Referer", "https://app.ourskylight.com/auth/session/new")

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("submitting credentials: %w", err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("login failed (%d): check email and password", resp.StatusCode)
	}

	// Step 3: exchange session for an authorization code
	authorizeURL := "https://app.ourskylight.com/oauth/authorize?" + url.Values{
		"client_id":     {oauthClientID},
		"redirect_uri":  {"https://ourskylight.com/welcome"},
		"response_type": {"code"},
		"scope":         {"everything"},
	}.Encode()

	req, _ = http.NewRequest("GET", authorizeURL, nil)
	req.Header.Set("User-Agent", "SkylightMobile (web)")

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting authorization: %w", err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	location := resp.Header.Get("Location")
	redirectURL, err := url.Parse(location)
	if err != nil || redirectURL.Query().Get("code") == "" {
		return nil, fmt.Errorf("no authorization code in redirect: %q", location)
	}
	code := redirectURL.Query().Get("code")

	// Step 4: exchange code for tokens
	return exchangeCode(code)
}

func exchangeCode(code string) (*OAuthTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", oauthClientID)
	form.Set("code", code)
	form.Set("redirect_uri", "https://ourskylight.com/welcome")
	form.Set("scope", "everything")

	req, _ := http.NewRequest("POST", oauthURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "SkylightMobile (web)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var tok OAuthTokenResponse
	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}
	if tok.AccessToken == "" {
		return nil, fmt.Errorf("token response missing access_token: %s", string(body))
	}
	return &tok, nil
}

// RefreshToken exchanges a refresh_token for a new access_token via POST /oauth/token.
func RefreshToken(refreshToken string) (*OAuthTokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", oauthClientID)
	form.Set("refresh_token", refreshToken)
	form.Set("scope", "everything")

	req, _ := http.NewRequest("POST", oauthURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "SkylightMobile (web)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("token refresh failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var tok OAuthTokenResponse
	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, fmt.Errorf("parsing refresh response: %w", err)
	}
	return &tok, nil
}

// Config represents the skylight CLI configuration file.
type Config struct {
	Token        string `yaml:"token"`                   // OAuth Bearer access_token
	RefreshToken string `yaml:"refresh_token,omitempty"` // OAuth refresh_token
	FrameID      string `yaml:"frame_id,omitempty"`
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

// ResolveToken returns the Bearer access token from SKYLIGHT_TOKEN env var or config file.
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
