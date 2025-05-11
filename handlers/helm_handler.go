package handlers

import (
	"fmt"
	"net/http"

	"helm-viewer/models"

	"github.com/gin-gonic/gin"
)

// HELMService defines the interface for HELM operations
type HELMService interface {
	LoadAndParseYAML(url string) (any, error)
	FindContainerImages(yamlContent any) []models.ContainerImage
	GetImageInfo(imageName string) (string, int, error)
}

// HELMHandler handles requests related to YAML documents
type HELMHandler struct {
	helmService HELMService
}

// NewHELMHandler creates a new instance of HELMHandler
func NewHELMHandler(helmService HELMService) *HELMHandler {
	return &HELMHandler{
		helmService: helmService,
	}
}

// LoadHELM handles the request to load a YAML document from URL
func (h *HELMHandler) LoadHELM(c *gin.Context) {
	var request models.HELMRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.HELMResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	// Load and parse YAML
	yamlContent, err := h.helmService.LoadAndParseYAML(request.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.HELMResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Find container images
	images := h.helmService.FindContainerImages(yamlContent)

	// Get size for each image
	for i := range images {
		size, layers, err := h.helmService.GetImageInfo(images[i].Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.HELMResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to get size for image %s: %v", images[i].Name, err),
			})
			return
		}
		images[i].Size = size
		images[i].Layers = layers
	}

	c.JSON(http.StatusOK, models.ImagesResponse{
		Success: true,
		Images:  images,
	})
}
