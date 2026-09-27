package config

import (
	"bufio"
	"errors"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/kartFr/Asset-Reuploader/internal/files"
)

var (
	mu            sync.RWMutex
	config        = make(map[string]string, 0)
	defaultConfig = map[string]string{
		"port":                "38073",
		"cookie_file":         "cookie.txt",
		"api_key":             "",
		"api_key_file":        "api-key.txt",
		"starts_per_minute":   "420",
		"max_concurrent":      "24",
		"max_parallel_chunks": "6",
		"upload_tries":        "3",
		// Per-type overrides (empty = fallback to global/hardcoded per-type defaults).
		"animation_starts_per_minute": "",
		"animation_max_concurrent":    "",
		"animation_upload_tries":      "",
		"sound_starts_per_minute":     "",
		"sound_max_concurrent":        "",
		"mesh_starts_per_minute":      "",
		"mesh_max_concurrent":         "",
	}
)

func init() {
	contents, err := files.Read("config.ini")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("failed reading config.ini, using defaults: %v", err)
	}
	if err != nil {
		contents = ""
	}

	parsed := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(contents))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		split := strings.SplitN(line, "=", 2)
		if len(split) != 2 {
			continue
		}

		key := strings.TrimSpace(split[0])
		if key == "" {
			continue
		}
		parsed[key] = strings.TrimSpace(split[1])
	}

	mu.Lock()
	for k, v := range parsed {
		config[k] = v
	}

	for i, v := range defaultConfig {
		if _, exists := config[i]; exists {
			continue
		}
		config[i] = v
	}

	keyFile := config["api_key_file"]
	mu.Unlock()
	data, err := files.Read(keyFile)
	switch {
	case err == nil && strings.TrimSpace(data) != "":
		mu.Lock()
		config["api_key"] = strings.TrimSpace(data)
		mu.Unlock()
	case err != nil && errors.Is(err, os.ErrNotExist):
		mu.RLock()
		k := strings.TrimSpace(config["api_key"])
		mu.RUnlock()
		if k != "" {
			if wErr := files.Write(keyFile, k); wErr != nil {
				log.Printf("could not migrate api key to %s: %v", keyFile, wErr)
			}
		}
	}

}

// PersistAPIKey writes the current api_key to api-key_file (Open Cloud key). Call after Set("api_key", ...).
func PersistAPIKey() error {
	mu.RLock()
	k := strings.TrimSpace(config["api_key"])
	keyFile := config["api_key_file"]
	mu.RUnlock()
	if k == "" {
		return nil
	}
	return files.Write(keyFile, k)
}

func Get(key string) string {
	mu.RLock()
	defer mu.RUnlock()
	return config[key]
}

func Set(key string, value string) {
	mu.Lock()
	defer mu.Unlock()
	config[key] = value
}

// CookieFile returns the cookie filename dynamically (do not cache at import time).
func CookieFile() string { return Get("cookie_file") }

// Port returns the HTTP port dynamically.
func Port() string { return Get("port") }

// APIKeyFile returns the api-key filename dynamically.
func APIKeyFile() string { return Get("api_key_file") }

// GetInt returns parsed int or fallback.
func GetInt(key string, fallback int) int {
	return parseIntOr(Get(key), fallback)
}

func parseIntOr(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n := 0
	neg := false
	for i, c := range s {
		if i == 0 && c == '-' {
			neg = true
			continue
		}
		if c < '0' || c > '9' {
			return fallback
		}
		n = n*10 + int(c-'0')
	}
	if neg {
		n = -n
	}
	return n
}

// StartsPerMinute returns rate tuning (default 420).
func StartsPerMinute() int { return GetInt("starts_per_minute", 420) }

// MaxConcurrent returns concurrency cap (default 24).
func MaxConcurrent() int { return GetInt("max_concurrent", 24) }

// MaxParallelChunks returns chunk parallelism (default 6).
func MaxParallelChunks() int { return GetInt("max_parallel_chunks", 6) }

// UploadTries returns upload retry count (default 3).
func UploadTries() int { return GetInt("upload_tries", 3) }

func Save() error {
	mu.RLock()
	keys := make([]string, 0, len(config))
	for key := range config {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var out strings.Builder
	for _, key := range keys {
		if key == "api_key" {
			continue
		}
		out.WriteString(key)
		out.WriteByte('=')
		out.WriteString(config[key])
		out.WriteByte('\n')
	}
	mu.RUnlock()
	return files.Write("config.ini", out.String())
}
