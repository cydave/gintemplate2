package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Init() (*viper.Viper, error) {
	cfg := viper.New()
	cfg.SetConfigType("yaml")
	cfg.SetDefault("environment", "production")
	cfg.SetDefault("server.host", "127.0.0.1")
	cfg.SetDefault("server.port", 3000)
	cfg.SetDefault("server.proto", "http")
	cfg.SetDefault("server.domain", "example.com")

	cfg.SetDefault("cookie.secret", randomSecretKey(32))
	cfg.SetDefault("cookie.path", "/")
	cfg.SetDefault("cookie.domain", "127.0.0.1")
	cfg.SetDefault("cookie.max_age", 0)
	cfg.SetDefault("cookie.secure", false)
	cfg.SetDefault("cookie.http_only", true)

	config = cfg
	return config, nil
}

func InitFrom(path string) (*viper.Viper, error) {
	cfg, err := Init()
	if err != nil {
		return nil, err
	}
	fh, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration from %s: %w", path, err)
	}
	defer fh.Close()
	err = cfg.MergeConfig(fh)
	if err != nil {
		return nil, fmt.Errorf("failed to merge configuration from %s: %w", path, err)
	}
	config = cfg
	return config, nil
}

func Get() *viper.Viper {
	return config
}

func GetSessionSecret() []byte {
	cfg := Get()
	s, err := hex.DecodeString(cfg.GetString("cookie.secret"))
	if err != nil {
		panic(err)
	}
	return s
}

func randomSecretKey(size int) string {
	secretKey := make([]byte, size)
	_, err := rand.Read(secretKey)
	if err != nil {
		return "not-so-random-secret-key"
	}
	return hex.EncodeToString(secretKey)
}
