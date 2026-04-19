package vault

import (
	"context"
	"fmt"
	"time"
)

// TokenInfo holds information about the current Vault token.
type TokenInfo struct {
	Accessor   string
	TTL        time.Duration
	Renewable  bool
	ExpireTime time.Time
}

// LookupSelfResponse is the minimal structure we care about from auth/token/lookup-self.
type lookupSelfResponse struct {
	Data struct {
		Accessor  string `json:"accessor"`
		TTL       int    `json:"ttl"`
		Renewable bool   `json:"renewable"`
	} `json:"data"`
}

// LookupToken calls auth/token/lookup-self and returns token metadata.
func (c *Client) LookupToken(ctx context.Context) (*TokenInfo, error) {
	var resp lookupSelfResponse
	if err := c.getJSON(ctx, "v1/auth/token/lookup-self", &resp); err != nil {
		return nil, fmt.Errorf("lookup token: %w", err)
	}
	ttl := time.Duration(resp.Data.TTL) * time.Second
	return &TokenInfo{
		Accessor:  resp.Data.Accessor,
		TTL:       ttl,
		Renewable: resp.Data.Renewable,
		ExpireTime: time.Now().Add(ttl),
	}, nil
}

// RenewToken attempts to renew the current token by the given increment.
// If increment is zero, Vault uses the token's default TTL.
func (c *Client) RenewToken(ctx context.Context, increment time.Duration) error {
	body := map[string]interface{}{}
	if increment > 0 {
		body["increment"] = int(increment.Seconds())
	}
	if err := c.postJSON(ctx, "v1/auth/token/renew-self", body); err != nil {
		return fmt.Errorf("renew token: %w", err)
	}
	return nil
}
