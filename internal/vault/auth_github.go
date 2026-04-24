package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GitHubLoginResponse holds the client token returned after GitHub auth.
type GitHubLoginResponse struct {
	Token     string
	LeaseDur  int
	Renewable bool
}

// GitHubLogin authenticates with Vault using a GitHub personal access token.
// mount defaults to "github" if empty.
func GitHubLogin(c *Client, githubToken string, mount string) (*GitHubLoginResponse, error) {
	if mount == "" {
		mount = "github"
	}

	body := fmt.Sprintf(`{"token":%q}`, githubToken)
	path := fmt.Sprintf("/v1/auth/%s/login", strings.Trim(mount, "/"))

	req, err := http.NewRequest(http.MethodPost, c.Address+path, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("github login: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github login: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github login: unexpected status %d", resp.StatusCode)
	}

	var raw struct {
		Auth struct {
			ClientToken   string `json:"client_token"`
			LeaseDuration int    `json:"lease_duration"`
			Renewable     bool   `json:"renewable"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("github login: decode: %w", err)
	}
	if raw.Auth.ClientToken == "" {
		return nil, fmt.Errorf("github login: empty token in response")
	}
	return &GitHubLoginResponse{
		Token:     raw.Auth.ClientToken,
		LeaseDur:  raw.Auth.LeaseDuration,
		Renewable: raw.Auth.Renewable,
	}, nil
}
