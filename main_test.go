package main

import (
	"testing"

	"helm-viewer/config"
	"helm-viewer/router"

	"github.com/stretchr/testify/require"
)

func TestConfigInitialization(t *testing.T) {
	cfg := config.NewConfig()
	require.NotNil(t, cfg)
	require.NotEmpty(t, cfg.Port)
}

func TestRouterSetup(t *testing.T) {
	r := router.SetupRouter()
	require.NotNil(t, r)
}

func TestServerStartErrorHandling(t *testing.T) {
	r := router.SetupRouter()
	err := r.Run("invalid-port")
	require.Error(t, err)
}
