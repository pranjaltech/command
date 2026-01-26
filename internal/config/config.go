package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds persistent settings.
// Provider stores credentials for an AI model provider.
type Provider struct {
	APIKey string `mapstructure:"api_key"`
	APIURL string `mapstructure:"api_url"`
}

// Config holds persistent settings. Multiple providers can be stored so users
// may switch without re-entering credentials.
type Config struct {
	Provider         string              `mapstructure:"provider"`
	Providers        map[string]Provider `mapstructure:"providers"`
	Model            string              `mapstructure:"model"`
	TelemetryDisable bool                `mapstructure:"telemetry_disable"`
	// APIKey is kept for backward compatibility with older configs.
	APIKey string `mapstructure:"api_key"`
}

const (
	secret = "01234567890123456789012345678901"

	// Provider-specific default models (2026 latest light/fast models)
	DefaultModelOpenAI     = "gpt-4.1-nano"       // OpenAI's fastest & cheapest model with 1M context
	DefaultModelAnthropic  = "claude-haiku-4-5"   // 4-5x faster than Sonnet 4.5 at fraction of cost
	DefaultModelOpenRouter = "x-ai/grok-4.1-fast" // xAI's best agentic tool calling model with 2M context
	DefaultModelGemini     = "gemini-3-flash"     // Pro-level intelligence at Flash speed, 3x faster than 2.5 Pro
	DefaultModelOllama     = "qwen3:8b"           // Best local lightweight model, 25 tokens/sec on laptops

	// DefaultModel is kept for backward compatibility
	DefaultModel = DefaultModelOpenAI
)

// DefaultModelForProvider returns the default model for a given provider.
func DefaultModelForProvider(provider string) string {
	switch provider {
	case "openai":
		return DefaultModelOpenAI
	case "anthropic":
		return DefaultModelAnthropic
	case "openrouter":
		return DefaultModelOpenRouter
	case "gemini":
		return DefaultModelGemini
	case "ollama":
		return DefaultModelOllama
	default:
		return DefaultModelOpenAI
	}
}

func cfgPath() string {
	if p := os.Getenv("CMD_CONFIG"); p != "" {
		return p
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "cmd", "config.yaml")
}

func ensureDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0o700)
}

func encrypt(plain string) (string, error) {
	key := []byte(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	out := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(out), nil
}

func decrypt(enc string) (string, error) {
	key := []byte(secret)
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce := data[:gcm.NonceSize()]
	plaintext, err := gcm.Open(nil, nonce, data[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// Load reads configuration from disk. Missing file is not an error.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(cfgPath())
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			if !os.IsNotExist(err) {
				return nil, err
			}
		}
	}
	var c Config
	_ = v.Unmarshal(&c)
	if c.Providers == nil {
		c.Providers = make(map[string]Provider)
	}
	for name, p := range c.Providers {
		if p.APIKey != "" {
			dec, err := decrypt(p.APIKey)
			if err != nil {
				return nil, err
			}
			p.APIKey = dec
			c.Providers[name] = p
		}
	}
	if c.APIKey != "" {
		dec, err := decrypt(c.APIKey)
		if err == nil {
			c.Providers["openai"] = Provider{APIKey: dec}
			if c.Provider == "" {
				c.Provider = "openai"
			}
		}
		c.APIKey = ""
	}
	if v.GetString("model") == "" {
		c.Model = DefaultModel
	}
	return &c, nil
}

// Save writes configuration to disk.
func Save(c *Config) error {
	path := cfgPath()
	if err := ensureDir(path); err != nil {
		return err
	}
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	prov := make(map[string]map[string]string)
	for name, p := range c.Providers {
		encKey := ""
		if p.APIKey != "" {
			var err error
			encKey, err = encrypt(p.APIKey)
			if err != nil {
				return err
			}
		}
		prov[name] = map[string]string{
			"api_key": encKey,
			"api_url": p.APIURL,
		}
	}
	v.Set("provider", c.Provider)
	v.Set("providers", prov)
	v.Set("model", c.Model)
	v.Set("telemetry_disable", c.TelemetryDisable)
	return v.WriteConfigAs(path)
}
