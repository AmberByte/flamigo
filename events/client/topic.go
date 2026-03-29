package client

import "strings"

func BuildClientTopic(parts ...string) string {
	return strings.Join(parts, ":")
}
