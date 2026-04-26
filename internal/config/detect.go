package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DetectType inspects files in dir and returns the project type, or "" if unknown.
func DetectType(dir string) string {
	if isLaravel(dir) {
		return "laravel"
	}
	if isNode(dir) {
		return "node"
	}
	return ""
}

func isLaravel(dir string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "composer.json"))
	if err != nil {
		return false
	}
	var composer struct {
		Require map[string]string `json:"require"`
	}
	if err := json.Unmarshal(data, &composer); err != nil {
		return false
	}
	_, ok := composer.Require["laravel/framework"]
	return ok
}

func isNode(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "package.json"))
	return err == nil
}
