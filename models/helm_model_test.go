package models

import (
	"encoding/json"
	"testing"
)

func TestYAMLRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request HELMRequest
		valid   bool
	}{
		{
			name: "Valid request with URL",
			request: HELMRequest{
				URL: "https://example.com/config.yaml",
			},
			valid: true,
		},
		{
			name: "Invalid request without URL",
			request: HELMRequest{
				URL: "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal and unmarshal to test JSON tags
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			var unmarshaled HELMRequest
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			if unmarshaled.URL != tt.request.URL {
				t.Errorf("URL = %v, want %v", unmarshaled.URL, tt.request.URL)
			}
		})
	}
}

func TestYAMLResponse_Marshaling(t *testing.T) {
	tests := []struct {
		name     string
		response HELMResponse
		wantJSON string
	}{
		{
			name: "Success response with data",
			response: HELMResponse{
				Success: true,
				Data:    map[string]interface{}{"key": "value"},
			},
			wantJSON: `{"success":true,"data":{"key":"value"}}`,
		},
		{
			name: "Error response",
			response: HELMResponse{
				Success: false,
				Error:   "test error",
			},
			wantJSON: `{"success":false,"error":"test error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Failed to marshal response: %v", err)
			}

			if string(data) != tt.wantJSON {
				t.Errorf("JSON = %v, want %v", string(data), tt.wantJSON)
			}
		})
	}
}

func TestContainerImage_Marshaling(t *testing.T) {
	image := ContainerImage{
		Name:      "test-image",
		Container: "test-container",
		Size:      "100MB",
		Layers:    5,
	}

	data, err := json.Marshal(image)
	if err != nil {
		t.Fatalf("Failed to marshal image: %v", err)
	}

	var unmarshaled ContainerImage
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal image: %v", err)
	}

	if unmarshaled.Name != image.Name {
		t.Errorf("Name = %v, want %v", unmarshaled.Name, image.Name)
	}
	if unmarshaled.Container != image.Container {
		t.Errorf("Container = %v, want %v", unmarshaled.Container, image.Container)
	}
	if unmarshaled.Size != image.Size {
		t.Errorf("Size = %v, want %v", unmarshaled.Size, image.Size)
	}
	if unmarshaled.Layers != image.Layers {
		t.Errorf("Layers = %v, want %v", unmarshaled.Layers, image.Layers)
	}
}

func TestImagesResponse_Marshaling(t *testing.T) {
	response := ImagesResponse{
		Success: true,
		Images: []ContainerImage{
			{
				Name:   "image1",
				Size:   "100MB",
				Layers: 3,
			},
			{
				Name:   "image2",
				Size:   "200MB",
				Layers: 4,
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var unmarshaled ImagesResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.Success != response.Success {
		t.Errorf("Success = %v, want %v", unmarshaled.Success, response.Success)
	}
	if len(unmarshaled.Images) != len(response.Images) {
		t.Errorf("Images length = %v, want %v", len(unmarshaled.Images), len(response.Images))
	}
}
