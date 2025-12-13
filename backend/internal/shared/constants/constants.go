package constants

// Game constants
const (
	// DefaultGameQuestionCount is the default number of questions per game session
	DefaultGameQuestionCount = 10

	// MaxGameQuestionCount is the maximum number of questions per game session
	MaxGameQuestionCount = 20

	// MinGameQuestionCount is the minimum number of questions per game session
	MinGameQuestionCount = 1
)

// API constants
const (
	// DefaultPageLimit is the default pagination limit
	DefaultPageLimit = 20

	// MaxPageLimit is the maximum pagination limit
	MaxPageLimit = 100

	// MinPageLimit is the minimum pagination limit
	MinPageLimit = 1
)

// Timeout constants (in milliseconds)
const (
	// DefaultAPITimeout is the default API request timeout
	DefaultAPITimeout = 30000 // 30 seconds

	// DictionarySearchTimeout is the timeout for dictionary search operations
	DictionarySearchTimeout = 1000 // 1 second per SC-005

	// QuestionGenerationTimeout is the timeout for question generation
	QuestionGenerationTimeout = 1000 // 1 second per SC-003
)
