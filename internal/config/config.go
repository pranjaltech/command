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
type Config struct {
	APIKey string `mapstructure:"api_key"`
}

const secret = "01234567890123456789012345678901"

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
	if c.APIKey != "" {
		dec, err := decrypt(c.APIKey)
		if err != nil {
			return nil, err
		}
		c.APIKey = dec
	}
	return &c, nil
}

// Save writes configuration to disk.
func Save(c *Config) error {
	enc, err := encrypt(c.APIKey)
	if err != nil {
		return err
	}
	path := cfgPath()
	if err := ensureDir(path); err != nil {
		return err
	}
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.Set("api_key", enc)
	return v.WriteConfigAs(path)
}
