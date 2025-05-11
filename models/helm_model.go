package models

// YAMLDocument represents a loaded YAML document
type YAMLDocument struct {
	Content interface{} `json:"content"`
}

// HELMRequest represents a request to load YAML
type HELMRequest struct {
	URL string `json:"url" binding:"required"`
}

// HELMResponse represents the server response
type HELMResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ContainerImage represents container image information
type ContainerImage struct {
	Name      string `json:"name"`
	Container string `json:"container,omitempty"`
	Size      string `json:"size,omitempty"`
	Layers    int    `json:"layers"`
}

// ImagesResponse represents the response containing container images
type ImagesResponse struct {
	Success bool             `json:"success"`
	Images  []ContainerImage `json:"images"`
}
