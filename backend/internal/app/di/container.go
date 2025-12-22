package di

import (
	"context"

	config "github.com/english-coach/backend/configs"
	dictusecase "github.com/english-coach/backend/internal/modules/dictionary/usecase/get_word_detail"
	dictrepo "github.com/english-coach/backend/internal/modules/dictionary/infra/persistence/postgres"
	gamesvc "github.com/english-coach/backend/internal/modules/game/service"
	gamecreatesession "github.com/english-coach/backend/internal/modules/game/usecase/create_session"
	gamesubmitanswer "github.com/english-coach/backend/internal/modules/game/usecase/submit_answer"
	gamerepo "github.com/english-coach/backend/internal/modules/game/infra/persistence/postgres"
	userrepo "github.com/english-coach/backend/internal/modules/user/infra/persistence/postgres"
	usergetprofile "github.com/english-coach/backend/internal/modules/user/usecase/get_profile"
	userlogin "github.com/english-coach/backend/internal/modules/user/usecase/login"
	userregister "github.com/english-coach/backend/internal/modules/user/usecase/register"
	userupdateprofile "github.com/english-coach/backend/internal/modules/user/usecase/update_profile"
	"github.com/english-coach/backend/internal/platform/db"
	"github.com/english-coach/backend/internal/shared/auth"
	"github.com/english-coach/backend/internal/shared/logger"
	"github.com/english-coach/backend/internal/transport/http"
	"github.com/english-coach/backend/internal/transport/http/handler"
	"github.com/english-coach/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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

	// Services
	QuestionGeneratorService *gamesvc.QuestionGeneratorService

	// Use Cases
	GetWordDetailUC     *dictusecase.Handler
	CreateGameSessionUC *gamecreatesession.Handler
	SubmitAnswerUC      *gamesubmitanswer.Handler
	RegisterUC          *userregister.Handler
	LoginUC             *userlogin.Handler
	GetProfileUC        *usergetprofile.Handler
	UpdateProfileUC     *userupdateprofile.Handler

	// Handlers
	DictionaryHandler *handler.DictionaryHandler
	GameHandler       *handler.GameHandler
	UserHandler       *handler.UserHandler
	OpenAPIHandler    *handler.OpenAPIHandler

	// HTTP Server
	HTTPServer *http.Server

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
	appLogger.Logger.Info("Database connection established")

	// Log CORS configuration
	appLogger.Logger.Info("CORS configuration loaded",
		zap.Strings("allowed_origins", cfg.CORS.AllowedOrigins),
	)

	// Initialize JWT manager
	container.JWTManager = auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)

	// Initialize repositories
	container.DictionaryRepo = dictrepo.NewDictionaryRepository(pool)
	container.GameRepo = gamerepo.NewGameRepository(pool)
	container.UserRepo = userrepo.NewUserRepository(pool)

	// Initialize services
	container.QuestionGeneratorService = gamesvc.NewQuestionGeneratorService(
		container.DictionaryRepo.WordRepository(),
		appLogger.Logger,
	)

	// Initialize use cases
	container.GetWordDetailUC = dictusecase.NewHandler(
		container.DictionaryRepo.WordRepository(),
		container.DictionaryRepo.SenseRepository(),
		container.DictionaryRepo.LanguageRepository(),
		container.DictionaryRepo.LevelRepository(),
		container.DictionaryRepo.PartOfSpeechRepository(),
		pool,
		appLogger.Logger,
	)

	container.CreateGameSessionUC = gamecreatesession.NewHandler(
		container.GameRepo.GameSessionRepo(),
		container.GameRepo.GameQuestionRepo(),
		container.QuestionGeneratorService,
		appLogger.Logger,
	)

	container.SubmitAnswerUC = gamesubmitanswer.NewHandler(
		container.GameRepo.GameAnswerRepo(),
		container.GameRepo.GameQuestionRepo(),
		container.GameRepo.GameSessionRepo(),
		appLogger.Logger,
	)

	container.RegisterUC = userregister.NewHandler(
		container.UserRepo.UserRepository(),
		appLogger.Logger,
	)

	container.LoginUC = userlogin.NewHandler(
		container.UserRepo.UserRepository(),
		container.JWTManager,
		appLogger.Logger,
	)

	container.GetProfileUC = usergetprofile.NewHandler(
		container.UserRepo.UserProfileRepository(),
		appLogger.Logger,
	)

	container.UpdateProfileUC = userupdateprofile.NewHandler(
		container.UserRepo.UserProfileRepository(),
		appLogger.Logger,
	)

	// Initialize handlers
	container.DictionaryHandler = handler.NewDictionaryHandler(
		container.DictionaryRepo.LanguageRepository(),
		container.DictionaryRepo.TopicRepository(),
		container.DictionaryRepo.LevelRepository(),
		container.DictionaryRepo.WordRepository(),
		container.GetWordDetailUC,
		appLogger.Logger,
	)

	container.GameHandler = handler.NewGameHandler(
		container.CreateGameSessionUC,
		container.SubmitAnswerUC,
		container.GameRepo.GameQuestionRepo(),
		container.GameRepo.GameSessionRepo(),
		appLogger.Logger,
	)

	container.UserHandler = handler.NewUserHandler(
		container.RegisterUC,
		container.LoginUC,
		container.GetProfileUC,
		container.UpdateProfileUC,
		container.UserRepo.UserRepository(),
		container.UserRepo.UserProfileRepository(),
		appLogger.Logger,
	)

	container.OpenAPIHandler = handler.NewOpenAPIHandler(
		appLogger.Logger,
		"docs/openapi/openapi.yaml",
	)

	// Initialize middleware
	container.CORSMiddleware = middleware.CORS(cfg.CORS.AllowedOrigins)
	container.ErrorMiddleware = middleware.ErrorHandler(appLogger.Logger)
	container.LoggerMiddleware = middleware.LoggerMiddleware(appLogger.Logger)
	container.AuthMiddleware = middleware.AuthMiddleware(container.JWTManager)

	// Initialize HTTP server
	container.HTTPServer = http.NewServer(
		http.Config{
			Port:            cfg.Server.Port,
			ReadTimeout:     cfg.Server.ReadTimeout,
			WriteTimeout:    cfg.Server.WriteTimeout,
			IdleTimeout:     cfg.Server.IdleTimeout,
			ShutdownTimeout: cfg.Server.ShutdownTimeout,
		},
		appLogger.Logger,
		container.CORSMiddleware,
		container.ErrorMiddleware,
		container.LoggerMiddleware,
	)

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
