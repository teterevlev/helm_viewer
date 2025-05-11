package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := SetupRouter()
	require.NotNil(t, r)

	// Test if the router has the expected routes
	routes := r.Routes()

	// Check if we have the expected number of routes
	require.Greater(t, len(routes), 0)

	// Check if our specific endpoint exists
	found := false
	for _, route := range routes {
		if route.Path == "/api/helm/load" && route.Method == "POST" {
			found = true
			break
		}
	}
	require.True(t, found, "Expected POST /api/helm/load endpoint not found")
}

func TestHELMEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := SetupRouter()

	// Test the endpoint with a POST request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/helm/load", nil)
	r.ServeHTTP(w, req)

	// We expect a 400 Bad Request since we're not sending any data
	require.Equal(t, http.StatusBadRequest, w.Code)
}
