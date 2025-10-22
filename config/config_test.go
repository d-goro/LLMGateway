package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	testKeysJSON := `{
		"virtual_keys": {
			"vk_test_openai": {
				"provider": "openai",
				"api_key": "sk-test-key-123"
			},
			"vk_test_anthropic": {
				"provider": "anthropic",
				"api_key": "sk-ant-test-key-456"
			}
		}
	}`

	tmpFile, err := os.CreateTemp("", "keys-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(testKeysJSON))
	require.NoError(t, err)
	tmpFile.Close()

	os.Setenv("KEYS_FILE_PATH", tmpFile.Name())
	defer os.Unsetenv("KEYS_FILE_PATH")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Len(t, cfg.KeysConfig.VirtualKeys, 2)

	openaiKey, exists := cfg.KeysConfig.VirtualKeys["vk_test_openai"]
	assert.True(t, exists)
	assert.Equal(t, "openai", string(openaiKey.Provider))
	assert.Equal(t, "sk-test-key-123", openaiKey.APIKey)

	anthropicKey, exists := cfg.KeysConfig.VirtualKeys["vk_test_anthropic"]
	assert.True(t, exists)
	assert.Equal(t, "anthropic", string(anthropicKey.Provider))
}

func TestLoadDefaultPath(t *testing.T) {
	testKeysJSON := `{
		"virtual_keys": {
			"vk_test": {
				"provider": "openai",
				"api_key": "sk-test-key"
			}
		}
	}`

	err := os.WriteFile("keys.json", []byte(testKeysJSON), 0644)
	require.NoError(t, err)
	defer os.Remove("keys.json")

	os.Unsetenv("KEYS_FILE_PATH")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Len(t, cfg.KeysConfig.VirtualKeys, 1)
}

func TestValidateVirtualKey(t *testing.T) {
	testKeysJSON := `{
		"virtual_keys": {
			"vk_valid": {
				"provider": "openai",
				"api_key": "sk-test-key"
			}
		}
	}`

	tmpFile, err := os.CreateTemp("", "keys-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(testKeysJSON))
	tmpFile.Close()

	os.Setenv("KEYS_FILE_PATH", tmpFile.Name())
	defer os.Unsetenv("KEYS_FILE_PATH")

	cfg, err := Load()
	require.NoError(t, err)

	keyConfig, valid := cfg.ValidateVirtualKey("vk_valid")
	assert.True(t, valid)
	assert.Equal(t, "openai", string(keyConfig.Provider))

	_, valid = cfg.ValidateVirtualKey("vk_invalid")
	assert.False(t, valid)
}

func TestLoadInvalidFile(t *testing.T) {
	os.Setenv("KEYS_FILE_PATH", "nonexistent.json")
	defer os.Unsetenv("KEYS_FILE_PATH")

	_, err := Load()
	require.Error(t, err)
}

func TestLoadEmptyKeys(t *testing.T) {
	testKeysJSON := `{"virtual_keys": {}}`

	tmpFile, err := os.CreateTemp("", "keys-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(testKeysJSON))
	tmpFile.Close()

	os.Setenv("KEYS_FILE_PATH", tmpFile.Name())
	defer os.Unsetenv("KEYS_FILE_PATH")

	_, err = Load()
	require.Error(t, err)
}

func TestLoadEnvironmentVariables(t *testing.T) {
	testKeysJSON := `{
		"virtual_keys": {
			"vk_test": {
				"provider": "openai",
				"api_key": "sk-test-key"
			}
		}
	}`

	tmpFile, err := os.CreateTemp("", "keys-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(testKeysJSON))
	tmpFile.Close()

	os.Setenv("KEYS_FILE_PATH", tmpFile.Name())
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("LOG_TO_FILE", "true")
	os.Setenv("LOG_FILE_PATH", "/tmp/test.log")
	os.Setenv("QUOTA_ENABLED", "false")
	os.Setenv("QUOTA_LIMIT", "200")
	os.Setenv("REQUEST_TIMEOUT", "60")

	defer func() {
		os.Unsetenv("KEYS_FILE_PATH")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("LOG_TO_FILE")
		os.Unsetenv("LOG_FILE_PATH")
		os.Unsetenv("QUOTA_ENABLED")
		os.Unsetenv("QUOTA_LIMIT")
		os.Unsetenv("REQUEST_TIMEOUT")
	}()

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "9090", cfg.ServerPort)
	assert.True(t, cfg.LogToFile)
	assert.Equal(t, "/tmp/test.log", cfg.LogFilePath)
	assert.False(t, cfg.QuotaEnabled)
	assert.Equal(t, int64(200), cfg.QuotaLimit)
	assert.Equal(t, 60, cfg.RequestTimeout)
}
