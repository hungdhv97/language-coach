package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/english-coach/backend/internal/domain/game/model"
	"github.com/english-coach/backend/internal/domain/game/port"
)

// GameRepository implements game repository interfaces
type GameRepository struct {
	pool *pgxpool.Pool
}

// NewGameRepository creates a new game repository
func NewGameRepository(pool *pgxpool.Pool) *GameRepository {
	return &GameRepository{
		pool: pool,
	}
}

// gameSessionRepo is a wrapper that implements GameSessionRepository
type gameSessionRepo struct {
	*GameRepository
}

// gameQuestionRepo is a wrapper that implements GameQuestionRepository
type gameQuestionRepo struct {
	*GameRepository
}

// gameAnswerRepo is a wrapper that implements GameAnswerRepository
type gameAnswerRepo struct {
	*GameRepository
}

// GameSessionRepository implementation

// Create creates a new game session
func (r *gameSessionRepo) Create(ctx context.Context, session *model.GameSession) error {
	query := `
		INSERT INTO vocab_game_sessions (
			user_id, mode, source_language_id, target_language_id,
			topic_id, level_id, total_questions, correct_questions,
			started_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, started_at
	`
	
	var startedAt time.Time
	err := r.pool.QueryRow(ctx, query,
		session.UserID,
		session.Mode,
		session.SourceLanguageID,
		session.TargetLanguageID,
		session.TopicID,
		session.LevelID,
		session.TotalQuestions,
		session.CorrectQuestions,
		time.Now(),
	).Scan(&session.ID, &startedAt)
	
	if err != nil {
		return err
	}
	
	session.StartedAt = startedAt
	return nil
}

// FindByID returns a game session by ID
func (r *gameSessionRepo) FindByID(ctx context.Context, id int64) (*model.GameSession, error) {
	query := `
		SELECT id, user_id, mode, source_language_id, target_language_id,
		       topic_id, level_id, total_questions, correct_questions,
		       started_at, ended_at
		FROM vocab_game_sessions
		WHERE id = $1
	`
	
	var session model.GameSession
	var endedAt *time.Time
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.Mode,
		&session.SourceLanguageID,
		&session.TargetLanguageID,
		&session.TopicID,
		&session.LevelID,
		&session.TotalQuestions,
		&session.CorrectQuestions,
		&session.StartedAt,
		&endedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	session.EndedAt = endedAt
	return &session, nil
}

// Update updates a game session
func (r *gameSessionRepo) Update(ctx context.Context, session *model.GameSession) error {
	query := `
		UPDATE vocab_game_sessions
		SET total_questions = $2,
		    correct_questions = $3,
		    ended_at = $4
		WHERE id = $1
	`
	
	_, err := r.pool.Exec(ctx, query,
		session.ID,
		session.TotalQuestions,
		session.CorrectQuestions,
		session.EndedAt,
	)
	
	return err
}

// EndSession marks a session as ended
func (r *gameSessionRepo) EndSession(ctx context.Context, sessionID int64, endedAt interface{}) error {
	var endTime time.Time
	if endedAt != nil {
		if t, ok := endedAt.(time.Time); ok {
			endTime = t
		} else {
			endTime = time.Now()
		}
	} else {
		endTime = time.Now()
	}
	
	query := `UPDATE vocab_game_sessions SET ended_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, sessionID, endTime)
	return err
}

// GameQuestionRepository implementation

// CreateBatch creates multiple questions and their options in a transaction
func (r *gameQuestionRepo) CreateBatch(ctx context.Context, questions []*model.GameQuestion, options []*model.GameQuestionOption) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert questions
	questionQuery := `
		INSERT INTO vocab_game_questions (
			session_id, question_order, question_type,
			source_word_id, source_sense_id, correct_target_word_id,
			source_language_id, target_language_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	for _, question := range questions {
		var createdAt time.Time
		err := tx.QueryRow(ctx, questionQuery,
			question.SessionID,
			question.QuestionOrder,
			question.QuestionType,
			question.SourceWordID,
			question.SourceSenseID,
			question.CorrectTargetWordID,
			question.SourceLanguageID,
			question.TargetLanguageID,
			time.Now(),
		).Scan(&question.ID, &createdAt)
		if err != nil {
			return err
		}
		question.CreatedAt = createdAt
	}

	// Insert options
	optionQuery := `
		INSERT INTO vocab_game_question_options (
			question_id, option_label, target_word_id, is_correct
		) VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	for _, option := range options {
		err := tx.QueryRow(ctx, optionQuery,
			option.QuestionID,
			option.OptionLabel,
			option.TargetWordID,
			option.IsCorrect,
		).Scan(&option.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// FindBySessionID returns all questions for a session
func (r *gameQuestionRepo) FindBySessionID(ctx context.Context, sessionID int64) ([]*model.GameQuestion, []*model.GameQuestionOption, error) {
	// Fetch questions
	questionQuery := `
		SELECT id, session_id, question_order, question_type,
		       source_word_id, source_sense_id, correct_target_word_id,
		       source_language_id, target_language_id, created_at
		FROM vocab_game_questions
		WHERE session_id = $1
		ORDER BY question_order
	`
	rows, err := r.pool.Query(ctx, questionQuery, sessionID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	questions := make([]*model.GameQuestion, 0)
	questionIDs := make([]int64, 0)
	for rows.Next() {
		var q model.GameQuestion
		var sourceSenseID *int64
		if err := rows.Scan(
			&q.ID,
			&q.SessionID,
			&q.QuestionOrder,
			&q.QuestionType,
			&q.SourceWordID,
			&sourceSenseID,
			&q.CorrectTargetWordID,
			&q.SourceLanguageID,
			&q.TargetLanguageID,
			&q.CreatedAt,
		); err != nil {
			return nil, nil, err
		}
		q.SourceSenseID = sourceSenseID
		questions = append(questions, &q)
		questionIDs = append(questionIDs, q.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	if len(questionIDs) == 0 {
		return questions, []*model.GameQuestionOption{}, nil
	}

	// Fetch options
	optionQuery := `
		SELECT id, question_id, option_label, target_word_id, is_correct
		FROM vocab_game_question_options
		WHERE question_id = ANY($1)
		ORDER BY question_id, option_label
	`
	optionRows, err := r.pool.Query(ctx, optionQuery, questionIDs)
	if err != nil {
		return nil, nil, err
	}
	defer optionRows.Close()

	options := make([]*model.GameQuestionOption, 0)
	for optionRows.Next() {
		var opt model.GameQuestionOption
		if err := optionRows.Scan(
			&opt.ID,
			&opt.QuestionID,
			&opt.OptionLabel,
			&opt.TargetWordID,
			&opt.IsCorrect,
		); err != nil {
			return nil, nil, err
		}
		options = append(options, &opt)
	}

	if err := optionRows.Err(); err != nil {
		return nil, nil, err
	}

	return questions, options, nil
}

// FindByID returns a question by ID with its options
func (r *gameQuestionRepo) FindByID(ctx context.Context, questionID int64) (*model.GameQuestion, []*model.GameQuestionOption, error) {
	// Fetch question
	questionQuery := `
		SELECT id, session_id, question_order, question_type,
		       source_word_id, source_sense_id, correct_target_word_id,
		       source_language_id, target_language_id, created_at
		FROM vocab_game_questions
		WHERE id = $1
	`
	var q model.GameQuestion
	var sourceSenseID *int64
	err := r.pool.QueryRow(ctx, questionQuery, questionID).Scan(
		&q.ID,
		&q.SessionID,
		&q.QuestionOrder,
		&q.QuestionType,
		&q.SourceWordID,
		&sourceSenseID,
		&q.CorrectTargetWordID,
		&q.SourceLanguageID,
		&q.TargetLanguageID,
		&q.CreatedAt,
	)
	if err != nil {
		return nil, nil, err
	}
	q.SourceSenseID = sourceSenseID

	// Fetch options
	optionQuery := `
		SELECT id, question_id, option_label, target_word_id, is_correct
		FROM vocab_game_question_options
		WHERE question_id = $1
		ORDER BY option_label
	`
	rows, err := r.pool.Query(ctx, optionQuery, questionID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	options := make([]*model.GameQuestionOption, 0)
	for rows.Next() {
		var opt model.GameQuestionOption
		if err := rows.Scan(
			&opt.ID,
			&opt.QuestionID,
			&opt.OptionLabel,
			&opt.TargetWordID,
			&opt.IsCorrect,
		); err != nil {
			return nil, nil, err
		}
		options = append(options, &opt)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return &q, options, nil
}

// GameAnswerRepository implementation

// Create creates a new answer
func (r *gameAnswerRepo) Create(ctx context.Context, answer *model.GameAnswer) error {
	query := `
		INSERT INTO vocab_game_question_answers (
			question_id, session_id, user_id,
			selected_option_id, is_correct, response_time_ms, answered_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, answered_at
	`
	var answeredAt time.Time
	err := r.pool.QueryRow(ctx, query,
		answer.QuestionID,
		answer.SessionID,
		answer.UserID,
		answer.SelectedOptionID,
		answer.IsCorrect,
		answer.ResponseTimeMs,
		time.Now(),
	).Scan(&answer.ID, &answeredAt)
	if err != nil {
		return err
	}
	answer.AnsweredAt = answeredAt
	return nil
}

// FindByQuestionID returns the answer for a specific question
func (r *gameAnswerRepo) FindByQuestionID(ctx context.Context, questionID, sessionID, userID int64) (*model.GameAnswer, error) {
	query := `
		SELECT id, question_id, session_id, user_id,
		       selected_option_id, is_correct, response_time_ms, answered_at
		FROM vocab_game_question_answers
		WHERE question_id = $1 AND session_id = $2 AND user_id = $3
		LIMIT 1
	`
	var answer model.GameAnswer
	var selectedOptionID *int64
	var responseTimeMs *int
	err := r.pool.QueryRow(ctx, query, questionID, sessionID, userID).Scan(
		&answer.ID,
		&answer.QuestionID,
		&answer.SessionID,
		&answer.UserID,
		&selectedOptionID,
		&answer.IsCorrect,
		&responseTimeMs,
		&answer.AnsweredAt,
	)
	if err != nil {
		return nil, err
	}
	answer.SelectedOptionID = selectedOptionID
	answer.ResponseTimeMs = responseTimeMs
	return &answer, nil
}

// FindBySessionID returns all answers for a session
func (r *gameAnswerRepo) FindBySessionID(ctx context.Context, sessionID, userID int64) ([]*model.GameAnswer, error) {
	query := `
		SELECT id, question_id, session_id, user_id,
		       selected_option_id, is_correct, response_time_ms, answered_at
		FROM vocab_game_question_answers
		WHERE session_id = $1 AND user_id = $2
		ORDER BY answered_at
	`
	rows, err := r.pool.Query(ctx, query, sessionID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	answers := make([]*model.GameAnswer, 0)
	for rows.Next() {
		var answer model.GameAnswer
		var selectedOptionID *int64
		var responseTimeMs *int
		if err := rows.Scan(
			&answer.ID,
			&answer.QuestionID,
			&answer.SessionID,
			&answer.UserID,
			&selectedOptionID,
			&answer.IsCorrect,
			&responseTimeMs,
			&answer.AnsweredAt,
		); err != nil {
			return nil, err
		}
		answer.SelectedOptionID = selectedOptionID
		answer.ResponseTimeMs = responseTimeMs
		answers = append(answers, &answer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}

// Ensure wrapper types implement the interfaces
var (
	_ port.GameSessionRepository = (*gameSessionRepo)(nil)
	_ port.GameQuestionRepository = (*gameQuestionRepo)(nil)
	_ port.GameAnswerRepository = (*gameAnswerRepo)(nil)
)

// GameSessionRepo returns a GameSessionRepository implementation
func (r *GameRepository) GameSessionRepo() port.GameSessionRepository {
	return &gameSessionRepo{GameRepository: r}
}

// GameQuestionRepo returns a GameQuestionRepository implementation
func (r *GameRepository) GameQuestionRepo() port.GameQuestionRepository {
	return &gameQuestionRepo{GameRepository: r}
}

// GameAnswerRepo returns a GameAnswerRepository implementation
func (r *GameRepository) GameAnswerRepo() port.GameAnswerRepository {
	return &gameAnswerRepo{GameRepository: r}
}

