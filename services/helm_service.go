package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"helm-viewer/models"

	"gopkg.in/yaml.v3"
)

// HELMService provides methods for working with YAML documents
type HELMService struct {
	dockerHubBaseURL string
}

// NewHELMService creates a new instance of HELMService
func NewHELMService() *HELMService {
	return &HELMService{
		dockerHubBaseURL: "https://hub.docker.com",
	}
}

// SetDockerHubBaseURL sets the base URL for Docker Hub API calls (used in testing)
func (s *HELMService) SetDockerHubBaseURL(url string) {
	s.dockerHubBaseURL = url
}

// LoadAndParseYAML loads and parses a YAML document from URL
func (s *HELMService) LoadAndParseYAML(url string) (any, error) {
	// Load YAML document from URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch YAML document: %w", err)
	}
	defer resp.Body.Close()

	// Read content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML content: %w", err)
	}

	// Parse YAML
	var yamlContent any
	if err := yaml.Unmarshal(body, &yamlContent); err != nil {
		return nil, fmt.Errorf("invalid YAML format: %w", err)
	}

	return yamlContent, nil
}

// FindContainerImages searches for container images in YAML structure
func (s *HELMService) FindContainerImages(content any) []models.ContainerImage {
	var images []models.ContainerImage

	switch v := content.(type) {
	case map[string]any:
		// Check for Helm chart image format
		if imageMap, ok := v["image"].(map[string]any); ok {
			repository, _ := imageMap["repository"].(string)
			tag, _ := imageMap["tag"].(string)

			if repository != "" {
				if tag == "" {
					tag = "latest"
				}

				images = append(images, models.ContainerImage{
					Name:      fmt.Sprintf("%s:%s", repository, tag),
					Container: "",
				})
			}
		}

		// Check for direct image string format
		if image, ok := v["image"].(string); ok {
			containerName := ""
			if name, ok := v["name"].(string); ok {
				containerName = name
			}

			images = append(images, models.ContainerImage{
				Name:      image,
				Container: containerName,
			})
		}

		// Recursively check all values in map
		for _, value := range v {
			images = append(images, s.FindContainerImages(value)...)
		}

	case []interface{}:
		// Recursively check all elements in slice
		for _, item := range v {
			images = append(images, s.FindContainerImages(item)...)
		}
	}

	return images
}

// formatSize converts bytes to human-readable format
func formatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	}
	return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
}

// GetDockerHubResponse makes a request to Docker Hub API and returns the response body
func (s *HELMService) GetDockerHubResponse(imageName string) ([]byte, error) {
	// Split image name into repository and tag
	parts := strings.Split(imageName, ":")
	repository := parts[0]
	tag := "latest"
	if len(parts) > 1 {
		tag = parts[1]
	}

	// Check if it's an official image
	var url string
	if !strings.Contains(repository, "/") {
		// Official image
		url = fmt.Sprintf("%s/v2/repositories/library/%s/tags/%s", s.dockerHubBaseURL, repository, tag)
	} else {
		// Custom image
		url = fmt.Sprintf("%s/v2/repositories/%s/tags/%s", s.dockerHubBaseURL, repository, tag)
	}

	// Make GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error requesting Docker Hub: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting image information: %s", resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return body, nil
}

// GetImageSize parses image size from Docker Hub API response
func (s *HELMService) GetImageSize(body []byte) (int64, error) {
	// Parse JSON response
	var result struct {
		FullSize int64 `json:"full_size"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("error parsing JSON: %w", err)
	}

	// Check if full_size field exists and is not zero
	if result.FullSize == 0 {
		return 0, fmt.Errorf("invalid response: full_size field is missing or zero")
	}

	return result.FullSize, nil
}

// GetImageLayers gets the number of layers for an image
// TODO: Implement layer counting by:
// 1. Parse the Docker Hub API response to get the digest
// 2. Use the digest to fetch the manifest from Docker Registry API
// 3. Parse the manifest to count the number of layers
// 4. Handle different manifest formats (v2 schema 1, v2 schema 2, OCI)
// 5. Consider caching the results to avoid repeated API calls
func (s *HELMService) GetImageLayers(body []byte) (int, error) {
	return 0, nil
}

// GetImageInfo gets image size from Docker Hub
func (s *HELMService) GetImageInfo(imageName string) (string, int, error) {
	body, err := s.GetDockerHubResponse(imageName)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get DockerHub response: %w", err)
	}

	size, err := s.GetImageSize(body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get image size: %w", err)
	}

	// Get number of layers
	layers, err := s.GetImageLayers(body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get image layers: %w", err)
	}

	return formatSize(size), layers, nil
}
