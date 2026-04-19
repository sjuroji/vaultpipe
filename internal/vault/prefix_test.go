package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterByPrefix_ReturnsMatchingKeys(t *testing.T) {
	secrets := map[string]string{
		"APP_DB_HOST": "localhost",
		"APP_DB_PASS": "secret",
		"OTHER_KEY":   "value",
	}
	result := FilterByPrefix(secrets, "APP_")
	assert.Len(t, result, 2)
	assert.Equal(t, "localhost", result["APP_DB_HOST"])
	assert.Equal(t, "secret", result["APP_DB_PASS"])
	_, ok := result["OTHER_KEY"]
	assert.False(t, ok)
}

func TestFilterByPrefix_EmptyPrefix_ReturnsAll(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	result := FilterByPrefix(secrets, "")
	assert.Equal(t, secrets, result)
}

func TestFilterByPrefix_NoMatch_ReturnsEmpty(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	result := FilterByPrefix(secrets, "NOPE_")
	assert.Empty(t, result)
}

func TestAddPrefix_PrependsPrefixToAllKeys(t *testing.T) {
	secrets := map[string]string{"HOST": "localhost", "PORT": "5432"}
	result := AddPrefix(secrets, "DB_")
	assert.Len(t, result, 2)
	assert.Equal(t, "localhost", result["DB_HOST"])
	assert.Equal(t, "5432", result["DB_PORT"])
}

func TestAddPrefix_EmptyPrefix_ReturnsUnchanged(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	result := AddPrefix(secrets, "")
	assert.Equal(t, secrets, result)
}
