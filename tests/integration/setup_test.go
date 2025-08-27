package integration

import (
	"blockstracker_backend/config"
	"blockstracker_backend/handlers"
	"blockstracker_backend/internal/redis"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/internal/validators"
	"blockstracker_backend/middleware"
	"blockstracker_backend/pkg/logger"
	"net/http"

	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/gin-gonic/gin"
	packageredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TestDB *gorm.DB
var originalGooseDBString string

func TestMain(m *testing.M) {
	defer logger.Log.Sync()
	validators.RegisterCustomValidators()

	db, err := createTestDB()
	if err != nil {
		log.Fatal(err)
	}

	TestDB, err = connectToTestDB()
	if err != nil {
		log.Fatal(err)
	}

	testSqlDB, err := TestDB.DB()
	if err != nil {
		log.Fatalf("Failed to get the test database instance: %v", err)
	}
	// todo: we can get rid of below lines, if we don't see any impact of the below lines after running many tests
	// testSqlDB.SetMaxIdleConns(1)
	// testSqlDB.SetMaxOpenConns(1)

	err = runGooseMigrations()
	if err != nil {
		log.Fatal(err)
	}

	err = initializeRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	err = setupRouter()
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := teardown(testSqlDB, db); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func createTestDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=db port=5432 user=%s password=%s dbname=postgres sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := gorm.Open(postgres.Open(dsn), config.GormConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to PostgreSQL: %v", err)
	}

	testDbName, ok := os.LookupEnv("TEST_DB_NAME")
	if !ok {
		return nil, fmt.Errorf("Test db not found in environment variables")
	}

	dropDbSqlString := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", testDbName)
	if err := db.Exec(dropDbSqlString).Error; err != nil {
		return nil, fmt.Errorf("Failed to drop test DB: %v", err)
	}

	createDbSqlString := fmt.Sprintf("CREATE DATABASE %s;", testDbName)
	if err := db.Exec(createDbSqlString).Error; err != nil {
		return nil, fmt.Errorf("Failed to create test DB: %v", err)
	}

	return db, nil
}

func connectToTestDB() (*gorm.DB, error) {

	testDSN := fmt.Sprintf(
		"host=db port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	TestDB, err := gorm.Open(postgres.Open(testDSN), config.GormConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to test DB: %v", err)
	}

	return TestDB, nil
}

func runGooseMigrations() error {
	originalGooseDBString = os.Getenv("GOOSE_DBSTRING")
	defer os.Setenv("GOOSE_DBSTRING", originalGooseDBString)

	test_GOOSE_DBSTRING := fmt.Sprintf(
		"postgres://%s:%s@db/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	err := os.Setenv("GOOSE_DBSTRING", test_GOOSE_DBSTRING)
	if err != nil {
		return fmt.Errorf("Failed to set GOOSE_DBSTRING for testing: %v", err)
	}

	cmd := exec.Command("goose", "-dir", "../../migrations", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "DB_USER="+os.Getenv("DB_USER"), "DB_PASSWORD="+os.Getenv("DB_PASSWORD"))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to run migrations on test DB: %v", err)
	}

	return nil
}

func killConnections(db *sql.DB, dbName string) error {
	_, err := db.Exec(`
     SELECT pg_terminate_backend(pid)
     FROM pg_stat_activity
     WHERE datname = $1 AND pid <> pg_backend_pid();
 `, dbName)
	if err != nil {
		return fmt.Errorf("failed to kill connections to database: %w", err)
	}

	return nil
}

func teardown(testSqlDB *sql.DB, db *gorm.DB) error {
	if err := killConnections(testSqlDB, os.Getenv("TEST_DB_NAME")); err != nil {
		log.Print(err.Error())
		fmt.Print(err.Error()) // we can still run tear down without killing connections, so keep going
	}

	testSqlDB.Close()

	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Printf("Failed to close Redis connection: %v", err)
		}
	}

	dropDbSqlString := fmt.Sprintf("DROP DATABASE %s;", os.Getenv("TEST_DB_NAME"))
	if err := db.Exec(dropDbSqlString).Error; err != nil {
		return fmt.Errorf("Failed to drop test DB: %v", err)
	}

	return nil

}

var redisClient *packageredis.Client

func initializeRedisClient() error {
	redisConfig, err := config.LoadRedisConfig()
	if err != nil {
		return fmt.Errorf("Error loading redis config: %v", err)
	}

	redisClient, err = redis.NewRedisClient(redisConfig)
	if err != nil {
		return fmt.Errorf("Error creating redis client: %v", err)
	}
	return nil
}

var router *gin.Engine
var testAuthConfig *config.AuthConfig

func setupRouter() error {
	var err error
	gin.SetMode(gin.TestMode)

	userRepo := repositories.NewUserRepository(TestDB)
	testAuthConfig, err = config.LoadAuthConfig()
	if err != nil {
		return fmt.Errorf("Error loading auth config: %v", err)
	}

	taskRepo := repositories.NewTaskRepository(TestDB)
	tagRepo := repositories.NewTagRepository(TestDB)
	spaceRepo := repositories.NewSpaceRepository(TestDB)
	changeRepo := repositories.NewChangeRepository(TestDB)

	logger := zap.NewNop().Sugar()

	tokenRepository := repositories.NewTokenRepository(redisClient)

	authHandler := handlers.NewAuthHandler(userRepo, logger, testAuthConfig, tokenRepository)
	authMiddleware := middleware.NewAuthMiddleware(logger, testAuthConfig)
	taskHandler := handlers.NewTaskHandler(taskRepo, changeRepo, TestDB, logger)
	tagHandler := handlers.NewTagHandler(tagRepo, changeRepo, TestDB, logger)
	spaceHandler := handlers.NewSpaceHandler(spaceRepo, changeRepo, TestDB, logger)

	router = gin.Default()
	router.POST("/signup", authHandler.SignupUser)
	router.POST("/signin", authHandler.EmailSignIn)
	router.POST("/refresh", authHandler.RefreshToken)
	router.POST("/protected", authMiddleware.Handle, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	router.POST("/signout", authMiddleware.Handle, authHandler.Signout)

	router.Use(authMiddleware.Handle)

	taskGroup := router.Group("/tasks")
	taskGroup.POST("/", taskHandler.CreateTask)
	taskGroup.PUT("/:id", taskHandler.UpdateTask)
	taskGroup.POST("/repetitive", taskHandler.CreateRepetitiveTaskTemplate)
	taskGroup.PUT("/repetitive/:id", taskHandler.UpdateRepetitiveTaskTemplate)

	tagGroup := router.Group("/tags")
	tagGroup.POST("/", tagHandler.CreateTag)
	tagGroup.PUT("/:id", tagHandler.UpdateTag)

	spaceGroup := router.Group("/spaces")
	spaceGroup.POST("/", spaceHandler.CreateSpace)
	spaceGroup.PUT("/:id", spaceHandler.UpdateSpace)

	return nil
}
