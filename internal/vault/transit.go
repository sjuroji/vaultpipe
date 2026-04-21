package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TransitEncryptRequest holds the plaintext (base64-encoded) to encrypt.
type TransitEncryptRequest struct {
	Plaintext string `json:"plaintext"`
}

// TransitEncryptResponse holds the ciphertext returned by Vault.
type TransitEncryptResponse struct {
	Ciphertext string `json:"ciphertext"`
}

// TransitDecryptResponse holds the plaintext returned by Vault.
type TransitDecryptResponse struct {
	Plaintext string `json:"plaintext"`
}

// EncryptTransit encrypts base64-encoded plaintext using the named transit key.
func EncryptTransit(c *Client, keyName, base64Plaintext string) (*TransitEncryptResponse, error) {
	body, _ := json.Marshal(TransitEncryptRequest{Plaintext: base64Plaintext})
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/v1/transit/encrypt/%s", c.Address, keyName),
		bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("transit encrypt: unexpected status %d", resp.StatusCode)
	}
	var wrapper struct {
		Data TransitEncryptResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

// DecryptTransit decrypts a Vault ciphertext using the named transit key.
func DecryptTransit(c *Client, keyName, ciphertext string) (*TransitDecryptResponse, error) {
	body, _ := json.Marshal(map[string]string{"ciphertext": ciphertext})
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/v1/transit/decrypt/%s", c.Address, keyName),
		bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("transit decrypt: unexpected status %d", resp.StatusCode)
	}
	var wrapper struct {
		Data TransitDecryptResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
