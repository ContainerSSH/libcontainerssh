package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.containerssh.io/libcontainerssh/config"
	"go.containerssh.io/libcontainerssh/internal/structutils"
)

func TestListenDefault(t *testing.T) {
	cfg := config.MetricsConfig{}
	structutils.Defaults(&cfg)
	assert.Equal(t, "0.0.0.0:9100", cfg.Listen)
}
