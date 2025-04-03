package strategies

import "strings"

// BuildCoreStrategyName builds a flamigo strategy name
// Name is prefixed by app:: automatically
// Example: BuildCoreStrategyName("jobs", "create") -> "app::jobs:create"
func BuildCoreStrategyName(parts ...string) string {
	return "app::" + strings.Join(parts, ":")
}
