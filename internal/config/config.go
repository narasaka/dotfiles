package config

import (
	"os"
	"strconv"
)

type Config struct {
	Dev           bool
	Port          int
	DBPath        string
	SessionSecret string
	Namespace     string
	AppNamespace  string
	RegistryURL   string
	RegistryUser  string
	RegistryPass  string
	BuildKitAddr  string
}

func Load() *Config {
	return &Config{
		Dev:           envBool("KUBEPLOY_DEV", false),
		Port:          envInt("KUBEPLOY_PORT", 8080),
		DBPath:        envStr("KUBEPLOY_DB_PATH", "./kubeploy.db"),
		SessionSecret: envStr("KUBEPLOY_SESSION_SECRET", ""),
		Namespace:     envStr("KUBEPLOY_NAMESPACE", "kubeploy-system"),
		AppNamespace:  envStr("KUBEPLOY_APP_NAMESPACE", "kubeploy-apps"),
		RegistryURL:   envStr("KUBEPLOY_REGISTRY_URL", ""),
		RegistryUser:  envStr("KUBEPLOY_REGISTRY_USERNAME", ""),
		RegistryPass:  envStr("KUBEPLOY_REGISTRY_PASSWORD", ""),
		BuildKitAddr:  envStr("KUBEPLOY_BUILDKIT_ADDR", "tcp://kubeploy-buildkitd:1234"),
	}
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
