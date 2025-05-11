package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHELMService(t *testing.T) {
	service := NewHELMService()
	require.NotNil(t, service)
}

func TestLoadAndParseYAML(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
image:
  repository: nginx
  tag: latest
containers:
  - name: web
    image: nginx:1.19
`))
	}))
	defer server.Close()

	service := NewHELMService()
	content, err := service.LoadAndParseYAML(server.URL)

	require.NoError(t, err)
	require.NotNil(t, content)

	// Verify the content structure
	contentMap, ok := content.(map[string]interface{})
	require.True(t, ok)

	// Check image section
	imageMap, ok := contentMap["image"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "nginx", imageMap["repository"])
	require.Equal(t, "latest", imageMap["tag"])
}

func TestFindContainerImages(t *testing.T) {
	svc := &HELMService{}

	t.Run("Helm style image map", func(t *testing.T) {
		yaml := map[string]interface{}{
			"image": map[string]interface{}{
				"repository": "nginx",
				"tag":        "1.21",
			},
		}
		imgs := svc.FindContainerImages(yaml)
		require.Len(t, imgs, 1)
		require.Equal(t, "nginx:1.21", imgs[0].Name)
	})

	t.Run("Helm style image map with no tag", func(t *testing.T) {
		yaml := map[string]interface{}{
			"image": map[string]interface{}{
				"repository": "nginx",
			},
		}
		imgs := svc.FindContainerImages(yaml)
		require.Len(t, imgs, 1)
		require.Equal(t, "nginx:latest", imgs[0].Name)
	})

	t.Run("Direct image string", func(t *testing.T) {
		yaml := map[string]interface{}{
			"image": "alpine:3.18",
			"name":  "test-container",
		}
		imgs := svc.FindContainerImages(yaml)
		require.Len(t, imgs, 1)
		require.Equal(t, "alpine:3.18", imgs[0].Name)
		require.Equal(t, "test-container", imgs[0].Container)
	})

	t.Run("Nested images", func(t *testing.T) {
		yaml := map[string]interface{}{
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"image": "busybox:1.36",
					},
				},
			},
		}
		imgs := svc.FindContainerImages(yaml)
		require.Len(t, imgs, 1)
		require.Equal(t, "busybox:1.36", imgs[0].Name)
	})
}

func TestFormatSize(t *testing.T) {
	testCases := []struct {
		size     int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := formatSize(tc.size)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestGetDockerHubResponse(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request URL
		if strings.Contains(r.URL.Path, "library/nginx") {
			response := map[string]interface{}{
				"full_size": 1024 * 1024 * 100, // 100MB
			}
			json.NewEncoder(w).Encode(response)
		} else if strings.Contains(r.URL.Path, "custom/nginx") {
			response := map[string]interface{}{
				"full_size": 1024 * 1024 * 50, // 50MB
			}
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := NewHELMService()
	service.SetDockerHubBaseURL(server.URL)

	// Test official image
	body, err := service.GetDockerHubResponse("nginx:latest")
	require.NoError(t, err)
	require.NotEmpty(t, body)

	// Test custom image
	body, err = service.GetDockerHubResponse("custom/nginx:latest")
	require.NoError(t, err)
	require.NotEmpty(t, body)
}

func TestGetImageSize(t *testing.T) {
	svc := &HELMService{}
	jsonBody := []byte(`{"full_size": 123456789}`)
	size, err := svc.GetImageSize(jsonBody)
	require.NoError(t, err)
	require.Equal(t, int64(123456789), size)

	t.Run("invalid json", func(t *testing.T) {
		_, err := svc.GetImageSize([]byte(`not json`))
		require.Error(t, err)
	})

	t.Run("zero size", func(t *testing.T) {
		_, err := svc.GetImageSize([]byte(`{"full_size":0}`))
		require.Error(t, err)
	})
}

func TestGetImageInfo(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"full_size": 104857600, // 100MB
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	service := NewHELMService()
	service.SetDockerHubBaseURL(server.URL)

	size, layers, err := service.GetImageInfo("nginx:latest")
	require.NoError(t, err)
	require.Equal(t, "100.00 MB", size)
	require.Equal(t, 0, layers) // Currently returns 0 as GetImageLayers is not implemented
}
