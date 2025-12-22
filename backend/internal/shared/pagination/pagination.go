package pagination

import (
	"strconv"

	"github.com/english-coach/backend/internal/shared/constants"
	"github.com/english-coach/backend/internal/shared/errors"
	"github.com/gin-gonic/gin"
)

// Params represents pagination parameters
type Params struct {
	Limit  int
	Offset int
	Page   int
	Size   int
}

// Metadata contains pagination metadata for responses
type Metadata struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

// ParseFromQuery parses pagination parameters from Gin query parameters
// Supports both page/pageSize and limit/offset approaches
// Priority: page/pageSize > limit/offset
func ParseFromQuery(c *gin.Context) (*Params, error) {
	params := &Params{
		Limit:  constants.DefaultPageLimit,
		Offset: 0,
		Page:   1,
		Size:   constants.DefaultPageLimit,
	}

	// Try to parse page/pageSize first (more user-friendly)
	pageStr := c.Query("page")
	sizeStr := c.Query("pageSize")

	if pageStr != "" || sizeStr != "" {
		// Parse page
		if pageStr != "" {
			page, err := strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				return nil, errors.ErrInvalidParameter.WithDetails("page must be a positive integer")
			}
			params.Page = page
		}

		// Parse pageSize
		if sizeStr != "" {
			size, err := strconv.Atoi(sizeStr)
			if err != nil || size < constants.MinPageLimit || size > constants.MaxPageLimit {
				return nil, errors.ErrInvalidParameter.WithDetails("pageSize must be between 1 and 100")
			}
			params.Size = size
		}

		// Convert page/pageSize to limit/offset
		params.Limit = params.Size
		params.Offset = (params.Page - 1) * params.Size
	} else {
		// Fallback to limit/offset
		limitStr := c.Query("limit")
		offsetStr := c.Query("offset")

		if limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil || limit < constants.MinPageLimit || limit > constants.MaxPageLimit {
				return nil, errors.ErrInvalidParameter.WithDetails("limit must be between 1 and 100")
			}
			params.Limit = limit
		}

		if offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err != nil || offset < 0 {
				return nil, errors.ErrInvalidParameter.WithDetails("offset must be non-negative")
			}
			params.Offset = offset
		}

		// Convert limit/offset to page/pageSize
		params.Size = params.Limit
		if params.Size > 0 {
			params.Page = (params.Offset / params.Size) + 1
		}
	}

	return params, nil
}

// CalculateMetadata calculates pagination metadata from params and total count
func CalculateMetadata(params *Params, total int64) *Metadata {
	totalPages := 0
	if params.Size > 0 {
		totalPages = int((total + int64(params.Size) - 1) / int64(params.Size))
	}

	hasNext := params.Page < totalPages
	hasPrev := params.Page > 1

	return &Metadata{
		Page:       params.Page,
		PageSize:   params.Size,
		Total:      total,
		TotalPages: totalPages,
		Limit:      params.Limit,
		Offset:     params.Offset,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}

// Validate validates pagination parameters
func Validate(params *Params) error {
	if params.Limit < constants.MinPageLimit || params.Limit > constants.MaxPageLimit {
		return errors.ErrInvalidParameter.WithDetails("limit must be between 1 and 100")
	}
	if params.Offset < 0 {
		return errors.ErrInvalidParameter.WithDetails("offset must be non-negative")
	}
	if params.Page < 1 {
		return errors.ErrInvalidParameter.WithDetails("page must be a positive integer")
	}
	if params.Size < constants.MinPageLimit || params.Size > constants.MaxPageLimit {
		return errors.ErrInvalidParameter.WithDetails("pageSize must be between 1 and 100")
	}
	return nil
}

