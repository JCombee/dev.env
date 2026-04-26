package compose

import "strings"

// ServiceName derives the docker-compose service name from an image+tag pair.
// Dots and colons are replaced with dashes: mysql:8.0 → mysql-8-0.
func ServiceName(image, tag string) string {
	if tag == "" {
		tag = "latest"
	}
	return image + "-" + strings.ReplaceAll(tag, ".", "-")
}
