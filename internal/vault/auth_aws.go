package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// AWSLogin authenticates using the AWS IAM or EC2 auth method.
// mount defaults to "aws" if empty.
func AWSLogin(c *Client, role, iamRequestURL, iamRequestBody, iamRequestHeaders, mount string) (string, error) {
	if mount == "" {
		mount = "aws"
	}

	payload := map[string]string{
		"role":                 role,
		"iam_http_request_method": "POST",
		"iam_request_url":     iamRequestURL,
		"iam_request_body":    iamRequestBody,
		"iam_request_headers": iamRequestHeaders,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("aws login: marshal payload: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/%s/login", mount)
	resp, err := c.Post(path, strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("aws login: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("aws login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("aws login: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("aws login: empty token in response")
	}

	return result.Auth.ClientToken, nil
}
