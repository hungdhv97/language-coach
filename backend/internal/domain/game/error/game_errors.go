package error

import (
	"github.com/english-coach/backend/internal/shared/errors"
)

// Game domain error codes
const (
	CodeInsufficientWords      = "INSUFFICIENT_WORDS"
	CodeSessionNotFound        = "SESSION_NOT_FOUND"
	CodeSessionEnded           = "SESSION_ENDED"
	CodeQuestionNotFound       = "QUESTION_NOT_FOUND"
	CodeQuestionNotInSession   = "QUESTION_NOT_IN_SESSION"
	CodeOptionNotFound         = "OPTION_NOT_FOUND"
	CodeAnswerAlreadySubmitted = "ANSWER_ALREADY_SUBMITTED"
	CodeInvalidMode            = "INVALID_MODE"
	CodeSessionNotOwned        = "SESSION_NOT_OWNED"
)

// Game domain errors
var (
	ErrInsufficientWords = errors.NewDomainError(
		CodeInsufficientWords,
		"Không đủ từ vựng để tạo phiên chơi. Vui lòng chọn chủ đề hoặc cấp độ khác",
		errors.StatusBadRequest,
	)

	ErrSessionNotFound = errors.NewDomainError(
		CodeSessionNotFound,
		"Không tìm thấy phiên chơi",
		errors.StatusNotFound,
	)

	ErrSessionEnded = errors.NewDomainError(
		CodeSessionEnded,
		"Phiên chơi đã kết thúc",
		errors.StatusBadRequest,
	)

	ErrQuestionNotFound = errors.NewDomainError(
		CodeQuestionNotFound,
		"Không tìm thấy câu hỏi",
		errors.StatusNotFound,
	)

	ErrQuestionNotInSession = errors.NewDomainError(
		CodeQuestionNotInSession,
		"Câu hỏi không thuộc về phiên chơi này",
		errors.StatusBadRequest,
	)

	ErrOptionNotFound = errors.NewDomainError(
		CodeOptionNotFound,
		"Không tìm thấy lựa chọn đã chọn",
		errors.StatusNotFound,
	)

	ErrAnswerAlreadySubmitted = errors.NewDomainError(
		CodeAnswerAlreadySubmitted,
		"Đã gửi câu trả lời cho câu hỏi này",
		errors.StatusBadRequest,
	)

	ErrInvalidMode = errors.NewDomainError(
		CodeInvalidMode,
		"Chế độ không hợp lệ",
		errors.StatusBadRequest,
	)

	ErrSessionNotOwned = errors.NewDomainError(
		CodeSessionNotOwned,
		"Phiên chơi không thuộc về người dùng này",
		errors.StatusForbidden,
	)
)
