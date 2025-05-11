# Helm viewer

A server for working with HELM files, written in Go using the Gin framework.

## Description

This server provides an API for working with HELM files. It allows loading and analyzing HELM charts through a REST API, extracting container image information, and retrieving their metadata from Docker Hub.

## Features

- Loading HELM files from URL
- Parsing HELM charts
- Finding container images in HELM structure
- Retrieving image size information from Docker Hub
- Support for both official and custom Docker images

## Requirements

- Go 1.21 or higher
- Gin framework
- YAML v3

## Installation

1. Clone the repository:
```bash
git clone https://github.com/teterevlev/helm_viewer
cd helm_viewer
```

2. Install dependencies:
```bash
go mod download
```

## Building

You can build the binary in two ways:

1. Using `go build`:
```bash
go build -o helm-viewer
```

2. Using `go install`:
```bash
go install
```

The `go build` command will create a binary in the current directory, while `go install` will install it in your `$GOPATH/bin` directory.

## Running

```bash
go run main.go
```

The server will start on the port specified in the `PORT` environment variable (default: 8080).

## API Endpoints

### POST /api/helm/load

Load and analyze a HELM chart.

#### Request Body
```json
{
    "url": "https://raw.githubusercontent.com/helm/examples/refs/heads/main/charts/hello-world/values.yaml"
}
```

#### Response
```json
{
    "success": true,
    "images": [
        {
            "name": "nginx:latest",
            "container": "web",
            "size": "133.7 MB",
            "layers": 0
        }
    ]
}
```

#### Possible Errors
- 400 Bad Request - Invalid request format
- 500 Internal Server Error - Error loading or processing YAML

## Dependencies

Main dependencies:
- github.com/gin-gonic/gin v1.9.1 - Web framework
- gopkg.in/yaml.v3 v3.0.1 - YAML parsing

## Implementation Details

- Support for recursive image search in YAML structure
- Automatic image tag detection (defaults to 'latest')
- Human-readable image size formatting
- Integration with Docker Hub API for image metadata

## License

MIT License

Copyright (c) 2024

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
