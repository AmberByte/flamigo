package strategies_test

import (
	"testing"

	"github.com/amberbyte/flamigo/strategies"
	"github.com/stretchr/testify/assert"
)

func TestBuildCoreStrategyName(t *testing.T) {
	assert.Equal(t, "app::jobs:create", strategies.BuildCoreStrategyName("jobs", "create"))
}
