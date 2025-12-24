package di

import (
	"context"

	config "github.com/english-coach/backend/configs"
	dictadapter "github.com/english-coach/backend/internal/modules/dictionary/adapter/http"
	dictrepo "github.com/english-coach/backend/internal/modules/dictionary/infra/persistence/postgres"
	dictusecase "github.com/english-coach/backend/internal/modules/dictionary/usecase/get_word_detail"
	gameadapter "github.com/english-coach/backend/internal/modules/game/adapter/http"
	gamerepo "github.com/english-coach/backend/internal/modules/game/infra/persistence/postgres"
	gamecreatesession "github.com/english-coach/backend/internal/modules/game/usecase/create_session"
	gamesubmitanswer "github.com/english-coach/backend/internal/modules/game/usecase/submit_answer"
	useradapter "github.com/english-coach/backend/internal/modules/user/adapter/http"
	userrepo "github.com/english-coach/backend/internal/modules/user/infra/persistence/postgres"
	usergetprofile "github.com/english-coach/backend/internal/modules/user/usecase/get_profile"
	userlogin "github.com/english-coach/backend/internal/modules/user/usecase/login"
	userregister "github.com/english-coach/backend/internal/modules/user/usecase/register"
	userupdateprofile "github.com/english-coach/backend/internal/modules/user/usecase/update_profile"
	"github.com/english-coach/backend/internal/platform/db"
	"github.com/english-coach/backend/internal/shared/auth"
	"github.com/english-coach/backend/internal/shared/logger"
	"github.com/english-coach/backend/internal/transport/http/handler"
	"github.com/english-coach/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Container holds all application dependencies
type Container struct {
	// Core
	Config *config.Config
	Logger *logger.Logger
	DB     *pgxpool.Pool

	// Auth
	JWTManager *auth.JWTManager

	// Repositories
	DictionaryRepo *dictrepo.DictionaryRepository
	GameRepo       *gamerepo.GameRepository
	UserRepo       *userrepo.UserRepository

	// Use Cases
	GetWordDetailUC     *dictusecase.Handler
	CreateGameSessionUC *gamecreatesession.Handler
	SubmitAnswerUC      *gamesubmitanswer.Handler
	RegisterUC          *userregister.Handler
	LoginUC             *userlogin.Handler
	GetProfileUC        *usergetprofile.Handler
	UpdateProfileUC     *userupdateprofile.Handler

	// Handlers
	DictionaryHandler *dictadapter.Handler
	GameHandler       *gameadapter.Handler
	UserHandler       *useradapter.Handler
	OpenAPIHandler    *handler.OpenAPIHandler

	// Middleware
	CORSMiddleware   gin.HandlerFunc
	ErrorMiddleware  gin.HandlerFunc
	LoggerMiddleware gin.HandlerFunc
	AuthMiddleware   gin.HandlerFunc
}

// NewContainer creates a new dependency injection container
func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{
		Config: cfg,
	}

	// Initialize logger
	appLogger, err := logger.NewLogger(cfg.App.Env)
	if err != nil {
		return nil, err
	}
	container.Logger = appLogger

	// Initialize database
	ctx := context.Background()
	pool, err := db.NewPostgres(ctx, db.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Database:        cfg.Database.Database,
		SSLMode:         cfg.Database.SSLMode,
		MaxConns:        cfg.Database.MaxConns,
		MinConns:        cfg.Database.MinConns,
		MaxConnLifetime: cfg.Database.MaxConnLifetime,
		MaxConnIdleTime: cfg.Database.MaxConnIdleTime,
	})
	if err != nil {
		return nil, err
	}
	container.DB = pool
	appLogger.Info("Database connection established")

	// Log CORS configuration
	appLogger.Info("CORS configuration loaded",
		logger.Strings("allowed_origins", cfg.CORS.AllowedOrigins),
	)

	// Initialize JWT manager
	container.JWTManager = auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)

	// Initialize repositories
	container.DictionaryRepo = dictrepo.NewDictionaryRepository(pool)
	container.GameRepo = gamerepo.NewGameRepository(pool)
	container.UserRepo = userrepo.NewUserRepository(pool)

	// Initialize use cases
	container.GetWordDetailUC = dictusecase.NewHandler(
		container.DictionaryRepo.WordRepository(),
		container.DictionaryRepo.SenseRepository(),
		container.DictionaryRepo.LanguageRepository(),
		container.DictionaryRepo.LevelRepository(),
		container.DictionaryRepo.PartOfSpeechRepository(),
		pool,
		appLogger,
	)

	container.CreateGameSessionUC = gamecreatesession.NewHandler(
		container.GameRepo.GameSessionRepo(),
		container.GameRepo.GameQuestionRepo(),
		container.DictionaryRepo.WordRepository(),
		appLogger,
	)

	container.SubmitAnswerUC = gamesubmitanswer.NewHandler(
		container.GameRepo.GameAnswerRepo(),
		container.GameRepo.GameQuestionRepo(),
		container.GameRepo.GameSessionRepo(),
		appLogger,
	)

	container.RegisterUC = userregister.NewHandler(
		container.UserRepo.UserRepository(),
		appLogger,
	)

	container.LoginUC = userlogin.NewHandler(
		container.UserRepo.UserRepository(),
		container.JWTManager,
		appLogger,
	)

	container.GetProfileUC = usergetprofile.NewHandler(
		container.UserRepo.UserProfileRepository(),
		appLogger,
	)

	container.UpdateProfileUC = userupdateprofile.NewHandler(
		container.UserRepo.UserProfileRepository(),
		appLogger,
	)

	// Initialize handlers
	container.DictionaryHandler = dictadapter.NewHandler(
		container.DictionaryRepo.LanguageRepository(),
		container.DictionaryRepo.TopicRepository(),
		container.DictionaryRepo.LevelRepository(),
		container.DictionaryRepo.WordRepository(),
		container.GetWordDetailUC,
		appLogger,
	)

	container.GameHandler = gameadapter.NewHandler(
		container.CreateGameSessionUC,
		container.SubmitAnswerUC,
		container.GameRepo.GameQuestionRepo(),
		container.GameRepo.GameSessionRepo(),
		appLogger,
	)

	container.UserHandler = useradapter.NewHandler(
		container.RegisterUC,
		container.LoginUC,
		container.GetProfileUC,
		container.UpdateProfileUC,
		container.UserRepo.UserRepository(),
		container.UserRepo.UserProfileRepository(),
		appLogger,
	)

	container.OpenAPIHandler = handler.NewOpenAPIHandler(
		appLogger,
		"docs/openapi/openapi.yaml",
	)

	// Initialize middleware
	container.CORSMiddleware = middleware.CORS(cfg.CORS.AllowedOrigins)
	container.ErrorMiddleware = middleware.ErrorHandler(appLogger)
	container.LoggerMiddleware = middleware.LoggerMiddleware(appLogger)
	container.AuthMiddleware = middleware.AuthMiddleware(container.JWTManager)

	return container, nil
}

// Close closes all resources in the container
func (c *Container) Close() error {
	if c.DB != nil {
		c.DB.Close()
	}
	if c.Logger != nil {
		c.Logger.Sync()
	}
	return nil
}
