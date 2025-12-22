package http

// SearchWordsRequest represents the query parameters for word search
type SearchWordsRequest struct {
	Query       string `form:"q" binding:"required"`
	LanguageID  int16  `form:"languageId" binding:"required"`
	Page        int    `form:"page"`
	PageSize    int    `form:"pageSize"`
	Limit       int    `form:"limit"`
	Offset      int    `form:"offset"`
}

// GetLevelsRequest represents the query parameters for getting levels
type GetLevelsRequest struct {
	LanguageID *int16 `form:"languageId"`
}

// GetWordDetailRequest represents the path parameter for getting word detail
type GetWordDetailRequest struct {
	WordID int64 `uri:"wordId" binding:"required"`
}

