package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"helm-viewer/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockHELMService is a mock implementation of HELMServiceInterface
type MockHELMService struct {
	mock.Mock
}

func (m *MockHELMService) LoadAndParseYAML(url string) (interface{}, error) {
	args := m.Called(url)
	return args.Get(0), args.Error(1)
}

func (m *MockHELMService) FindContainerImages(content interface{}) []models.ContainerImage {
	args := m.Called(content)
	return args.Get(0).([]models.ContainerImage)
}

func (m *MockHELMService) GetImageInfo(imageName string) (string, int, error) {
	args := m.Called(imageName)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockHELMService) SetDockerHubBaseURL(url string) {
	m.Called(url)
}

func setupTestRouter(handler *HELMHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/load-helm", handler.LoadHELM)
	return router
}

func TestLoadHELM_Success(t *testing.T) {
	// Setup
	mockService := new(MockHELMService)
	handler := NewHELMHandler(mockService)
	router := setupTestRouter(handler)

	// Mock data
	yamlContent := map[string]interface{}{
		"image": map[string]interface{}{
			"repository": "nginx",
			"tag":        "latest",
		},
	}

	images := []models.ContainerImage{
		{
			Name:      "nginx:latest",
			Container: "",
		},
	}

	// Setup expectations
	mockService.On("LoadAndParseYAML", "http://example.com/chart.yaml").Return(yamlContent, nil)
	mockService.On("FindContainerImages", yamlContent).Return(images)
	mockService.On("GetImageInfo", "nginx:latest").Return("100MB", 5, nil)

	// Create request
	reqBody := models.HELMRequest{URL: "http://example.com/chart.yaml"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/load-helm", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	require.Equal(t, http.StatusOK, w.Code)

	var response models.ImagesResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.True(t, response.Success)
	require.Len(t, response.Images, 1)
	require.Equal(t, "nginx:latest", response.Images[0].Name)
	require.Equal(t, "100MB", response.Images[0].Size)
	require.Equal(t, 5, response.Images[0].Layers)

	mockService.AssertExpectations(t)
}

func TestLoadHELM_InvalidRequest(t *testing.T) {
	// Setup
	mockService := new(MockHELMService)
	handler := NewHELMHandler(mockService)
	router := setupTestRouter(handler)

	// Create invalid request
	req := httptest.NewRequest(http.MethodPost, "/load-helm", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response models.HELMResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.False(t, response.Success)
	require.Equal(t, "Invalid request format", response.Error)
}

func TestLoadHELM_YAMLLoadError(t *testing.T) {
	// Setup
	mockService := new(MockHELMService)
	handler := NewHELMHandler(mockService)
	router := setupTestRouter(handler)

	// Setup expectations
	mockService.On("LoadAndParseYAML", "http://example.com/chart.yaml").Return(nil, assert.AnError)

	// Create request
	reqBody := models.HELMRequest{URL: "http://example.com/chart.yaml"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/load-helm", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	require.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.HELMResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.False(t, response.Success)
	require.Equal(t, assert.AnError.Error(), response.Error)
}

func TestLoadHELM_ImageInfoError(t *testing.T) {
	// Setup
	mockService := new(MockHELMService)
	handler := NewHELMHandler(mockService)
	router := setupTestRouter(handler)

	// Mock data
	yamlContent := map[string]interface{}{
		"image": map[string]interface{}{
			"repository": "nginx",
			"tag":        "latest",
		},
	}

	images := []models.ContainerImage{
		{
			Name:      "nginx:latest",
			Container: "",
		},
	}

	// Setup expectations
	mockService.On("LoadAndParseYAML", "http://example.com/chart.yaml").Return(yamlContent, nil)
	mockService.On("FindContainerImages", yamlContent).Return(images)
	mockService.On("GetImageInfo", "nginx:latest").Return("", 0, assert.AnError)

	// Create request
	reqBody := models.HELMRequest{URL: "http://example.com/chart.yaml"}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/load-helm", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	require.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.HELMResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.False(t, response.Success)
	require.Contains(t, response.Error, "Failed to get size for image nginx:latest")
}
