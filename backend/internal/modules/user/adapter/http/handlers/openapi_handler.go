package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// OpenAPIHandler handles OpenAPI specification serving
type OpenAPIHandler struct {
	logger     *zap.Logger
	specPath   string
	specLoaded bool
	specJSON   []byte
	baseDir    string
}

// NewOpenAPIHandler creates a new OpenAPI handler
func NewOpenAPIHandler(logger *zap.Logger, specPath string) *OpenAPIHandler {
	return &OpenAPIHandler{
		logger:   logger,
		specPath: specPath,
	}
}

// findSpecFile finds the OpenAPI spec file
func (h *OpenAPIHandler) findSpecFile() (string, error) {
	// Try to read from the provided path
	if _, err := os.Stat(h.specPath); err == nil {
		return h.specPath, nil
	}

	// Try relative path from current working directory
	wd, _ := os.Getwd()
	relativePath := filepath.Join(wd, "docs", "openapi", "openapi.yaml")
	if _, err := os.Stat(relativePath); err == nil {
		return relativePath, nil
	}

	// Try from backend directory
	backendPath := filepath.Join(wd, "backend", "docs", "openapi", "openapi.yaml")
	if _, err := os.Stat(backendPath); err == nil {
		return backendPath, nil
	}

	return "", fmt.Errorf("OpenAPI spec file not found")
}

// loadYAMLFile loads a YAML file and returns its content as a map
func (h *OpenAPIHandler) loadYAMLFile(filePath string) (map[string]interface{}, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// resolveExternalRef resolves an external $ref reference (file-based)
func (h *OpenAPIHandler) resolveExternalRef(ref string, baseDir string) (interface{}, error) {
	// Handle external file references (e.g., './paths/user.yaml#/paths/~1auth~1register')
	parts := strings.SplitN(ref, "#", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid external $ref format: %s", ref)
	}

	filePath := strings.TrimPrefix(parts[0], "./")
	filePath = filepath.Join(baseDir, filePath)
	jsonPath := parts[1]

	// Load the referenced file
	data, err := h.loadYAMLFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load file %s: %w", filePath, err)
	}

	// Navigate the JSON path
	pathParts := strings.Split(strings.TrimPrefix(jsonPath, "/"), "/")
	var result interface{} = data

	for _, part := range pathParts {
		if part == "" {
			continue
		}
		// Handle JSON pointer encoding (~1 for /, ~0 for ~)
		part = strings.ReplaceAll(part, "~1", "/")
		part = strings.ReplaceAll(part, "~0", "~")

		if m, ok := result.(map[string]interface{}); ok {
			result, ok = m[part]
			if !ok {
				return nil, fmt.Errorf("path not found: %s in %s", part, jsonPath)
			}
		} else {
			return nil, fmt.Errorf("invalid path: %s", jsonPath)
		}
	}

	return result, nil
}

// resolveInternalRef resolves an internal $ref reference within the merged spec
func (h *OpenAPIHandler) resolveInternalRef(ref string, spec map[string]interface{}) (interface{}, error) {
	if !strings.HasPrefix(ref, "#/") {
		return nil, fmt.Errorf("not an internal reference: %s", ref)
	}

	// Navigate the JSON path in the spec
	pathParts := strings.Split(strings.TrimPrefix(ref, "#/"), "/")
	var result interface{} = spec

	for _, part := range pathParts {
		if part == "" {
			continue
		}

		if m, ok := result.(map[string]interface{}); ok {
			result, ok = m[part]
			if !ok {
				return nil, fmt.Errorf("internal reference path not found: %s in %s", part, ref)
			}
		} else {
			return nil, fmt.Errorf("invalid internal reference path: %s", ref)
		}
	}

	return result, nil
}

// resolveExternalRefsInValue recursively resolves all external $ref references in a value
func (h *OpenAPIHandler) resolveExternalRefsInValue(value interface{}, baseDir string) (interface{}, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		// Check if this is a $ref
		if ref, ok := v["$ref"].(string); ok {
			// Only resolve external refs (not starting with #)
			if !strings.HasPrefix(ref, "#/") {
				resolved, err := h.resolveExternalRef(ref, baseDir)
				if err != nil {
					return nil, err
				}
				// Recursively resolve external refs in the resolved value
				return h.resolveExternalRefsInValue(resolved, baseDir)
			}
			// Internal refs will be resolved in the second pass
			return v, nil
		}

		// Recursively resolve external refs in all values
		resolved := make(map[string]interface{})
		for k, val := range v {
			resolvedVal, err := h.resolveExternalRefsInValue(val, baseDir)
			if err != nil {
				return nil, err
			}
			resolved[k] = resolvedVal
		}
		return resolved, nil

	case []interface{}:
		resolved := make([]interface{}, len(v))
		for i, val := range v {
			resolvedVal, err := h.resolveExternalRefsInValue(val, baseDir)
			if err != nil {
				return nil, err
			}
			resolved[i] = resolvedVal
		}
		return resolved, nil

	default:
		return value, nil
	}
}

// resolveInternalRefsInValue recursively resolves all internal $ref references in a value
func (h *OpenAPIHandler) resolveInternalRefsInValue(value interface{}, spec map[string]interface{}) (interface{}, error) {
	switch v := value.(type) {
	case map[string]interface{}:
		// Check if this is a $ref
		if ref, ok := v["$ref"].(string); ok {
			// Only resolve internal refs (starting with #)
			if strings.HasPrefix(ref, "#/") {
				resolved, err := h.resolveInternalRef(ref, spec)
				if err != nil {
					return nil, err
				}
				// Recursively resolve internal refs in the resolved value
				return h.resolveInternalRefsInValue(resolved, spec)
			}
			// External refs should have been resolved in the first pass
			return v, nil
		}

		// Recursively resolve internal refs in all values
		resolved := make(map[string]interface{})
		for k, val := range v {
			resolvedVal, err := h.resolveInternalRefsInValue(val, spec)
			if err != nil {
				return nil, err
			}
			resolved[k] = resolvedVal
		}
		return resolved, nil

	case []interface{}:
		resolved := make([]interface{}, len(v))
		for i, val := range v {
			resolvedVal, err := h.resolveInternalRefsInValue(val, spec)
			if err != nil {
				return nil, err
			}
			resolved[i] = resolvedVal
		}
		return resolved, nil

	default:
		return value, nil
	}
}

// loadSpec loads and resolves the OpenAPI specification
func (h *OpenAPIHandler) loadSpec() error {
	if h.specLoaded {
		return nil
	}

	// Find the spec file
	specFile, err := h.findSpecFile()
	if err != nil {
		return err
	}

	// Get base directory for resolving relative paths
	h.baseDir = filepath.Dir(specFile)

	// Load main spec file
	specData, err := h.loadYAMLFile(specFile)
	if err != nil {
		return fmt.Errorf("failed to load spec file: %w", err)
	}

	// First pass: Resolve all external $ref references (file-based)
	resolvedExternal, err := h.resolveExternalRefsInValue(specData, h.baseDir)
	if err != nil {
		return fmt.Errorf("failed to resolve external $ref references: %w", err)
	}

	// Convert to map for internal reference resolution
	specMap, ok := resolvedExternal.(map[string]interface{})
	if !ok {
		return fmt.Errorf("spec is not a map after external resolution")
	}

	// Second pass: Resolve all internal $ref references (within merged spec)
	resolvedSpec, err := h.resolveInternalRefsInValue(specMap, specMap)
	if err != nil {
		return fmt.Errorf("failed to resolve internal $ref references: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(resolvedSpec)
	if err != nil {
		return fmt.Errorf("failed to marshal spec to JSON: %w", err)
	}

	h.specJSON = jsonData
	h.specLoaded = true

	return nil
}

// getOpenAPISpecJSON returns the resolved OpenAPI spec as JSON
func (h *OpenAPIHandler) getOpenAPISpecJSON() ([]byte, error) {
	if err := h.loadSpec(); err != nil {
		return nil, err
	}
	return h.specJSON, nil
}

// GetSwaggerUI serves the Swagger UI HTML page with embedded spec
func (h *OpenAPIHandler) GetSwaggerUI(c *gin.Context) {
	// Load and resolve the spec
	specJSON, err := h.getOpenAPISpecJSON()
	if err != nil {
		h.logger.Error("failed to load OpenAPI spec",
			zap.Error(err),
			zap.String("path", h.specPath),
		)
		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Không thể tải đặc tả OpenAPI",
			nil,
		)
		return
	}

	// Escape the JSON for embedding in HTML/JavaScript
	specJSONEscaped := strings.ReplaceAll(string(specJSON), "`", "\\`")
	specJSONEscaped = strings.ReplaceAll(specJSONEscaped, "</script>", "<\\/script>")

	// Swagger UI HTML with embedded spec
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>English Coach API - Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.17.14/favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.17.14/favicon-16x16.png" sizes="16x16" />
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *, *:before, *:after {
      box-sizing: inherit;
    }
    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      const spec = %s;
      const ui = SwaggerUIBundle({
        spec: spec,
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        validatorUrl: null,
        tryItOutEnabled: true
      });
    };
  </script>
</body>
</html>`, specJSONEscaped)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
