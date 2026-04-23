package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PolicyListResponse holds the list of policy names.
type PolicyListResponse struct {
	Policies []string `json:"policies"`
}

// PolicyDetail holds a single policy's name and rules.
type PolicyDetail struct {
	Name  string `json:"name"`
	Rules string `json:"rules"`
}

// ListPolicies returns all ACL policies in Vault.
func ListPolicies(client *Client) (*PolicyListResponse, error) {
	resp, err := client.Get("/v1/sys/policies/acl?list=true")
	if err != nil {
		return nil, fmt.Errorf("list policies request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &PolicyListResponse{Policies: wrapper.Data.Keys}, nil
}

// GetPolicy retrieves a single ACL policy by name.
func GetPolicy(client *Client, name string) (*PolicyDetail, error) {
	resp, err := client.Get("/v1/sys/policies/acl/" + name)
	if err != nil {
		return nil, fmt.Errorf("get policy request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("policy %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data PolicyDetail `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &wrapper.Data, nil
}
