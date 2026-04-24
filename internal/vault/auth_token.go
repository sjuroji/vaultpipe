package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TokenLoginResponse holds the result of a token auth login.
type TokenLoginResponse struct {
	ClientToken string
	Accessor    string
	Policies    []string
	Renewable   bool
	TTL         int
}

// TokenLogin authenticates with Vault using a raw token and validates it by
// performing a self-lookup. It returns metadata about the token.
func TokenLogin(c *Client, token string) (*TokenLoginResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/auth/token/lookup-self", nil)
	if err != nil {
		return nil, fmt.Errorf("auth_token: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth_token: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth_token: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			ID        string   `json:"id"`
			Accessor  string   `json:"accessor"`
			Policies  []string `json:"policies"`
			Renewable bool     `json:"renewable"`
			TTL       int      `json:"ttl"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("auth_token: decode response: %w", err)
	}

	return &TokenLoginResponse{
		ClientToken: token,
		Accessor:    body.Data.Accessor,
		Policies:    body.Data.Policies,
		Renewable:   body.Data.Renewable,
		TTL:         body.Data.TTL,
	}, nil
}
