package config

import (
	"os"
	"strconv"
)

type Config struct {
	Dev             bool
	Port            int
	DBPath          string
	SessionSecret   string
	Namespace       string
	AppNamespace    string
	RegistryURL     string
	RegistryUser    string
	RegistryPass    string
}

func Load() *Config {
	return &Config{
		Dev:           envBool("KUBEDECK_DEV", false),
		Port:          envInt("KUBEDECK_PORT", 8080),
		DBPath:        envStr("KUBEDECK_DB_PATH", "./kubedeck.db"),
		SessionSecret: envStr("KUBEDECK_SESSION_SECRET", ""),
		Namespace:     envStr("KUBEDECK_NAMESPACE", "kubedeck-system"),
		AppNamespace:  envStr("KUBEDECK_APP_NAMESPACE", "kubedeck-apps"),
		RegistryURL:   envStr("KUBEDECK_REGISTRY_URL", ""),
		RegistryUser:  envStr("KUBEDECK_REGISTRY_USERNAME", ""),
		RegistryPass:  envStr("KUBEDECK_REGISTRY_PASSWORD", ""),
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
