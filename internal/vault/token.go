package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TokenInfo holds metadata about a Vault token.
type TokenInfo struct {
	ID           string
	DisplayName  string
	Policies     []string
	TTL          time.Duration
	Renewable    bool
	CreationTime time.Time
}

type tokenLookupResponse struct {
	Data struct {
		ID           string   `json:"id"`
		DisplayName  string   `json:"display_name"`
		Policies     []string `json:"policies"`
		TTL          int      `json:"ttl"`
		Renewable    bool     `json:"renewable"`
		CreationTime int64    `json:"creation_time"`
	} `json:"data"`
}

// LookupSelfToken retrieves metadata for the token currently configured on the client.
func LookupSelfToken(c *Client) (*TokenInfo, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/auth/token/lookup-self", nil)
	if err != nil {
		return nil, fmt.Errorf("token lookup: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token lookup: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token lookup: unexpected status %d", resp.StatusCode)
	}

	var out tokenLookupResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("token lookup: decode: %w", err)
	}

	return &TokenInfo{
		ID:           out.Data.ID,
		DisplayName:  out.Data.DisplayName,
		Policies:     out.Data.Policies,
		TTL:          time.Duration(out.Data.TTL) * time.Second,
		Renewable:    out.Data.Renewable,
		CreationTime: time.Unix(out.Data.CreationTime, 0),
	}, nil
}

// IsExpired returns true when the token TTL has reached zero.
func (t *TokenInfo) IsExpired() bool {
	return t.TTL <= 0
}
