package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"ienergy-template-go/config"
	"ienergy-template-go/internal/http/handler"
	"ienergy-template-go/internal/model/entity"
	"ienergy-template-go/internal/model/request"
	"ienergy-template-go/internal/repository"
	"ienergy-template-go/internal/service"
	"ienergy-template-go/pkg/database"
	"ienergy-template-go/pkg/logger"
	"ienergy-template-go/pkg/wrapper"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

// TestConfig holds test configuration
type TestConfig struct {
	DBConfig     config.DBConfig
	JWTConfig    config.JWTConfig
	ServerConfig config.ServerCfg
}

// TestEnvironment holds all test dependencies
type TestEnvironment struct {
	Router  *gin.Engine
	DB      database.Database
	Config  *TestConfig
	Cleanup func()
}

// TestUser holds test user data
type TestUser struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// DefaultTestUser returns a default test user
func DefaultTestUser() TestUser {
	return TestUser{
		Email:     "e2e_test@example.com",
		Password:  "password1234",
		FirstName: "E2E",
		LastName:  "Test",
	}
}

// DefaultTestConfig returns default test configuration
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		DBConfig: config.DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "Admin",
			Password: "admin@123",
			DBName:   "ienergy_db",
			SSLMode:  "disable",
		},
		JWTConfig: config.JWTConfig{
			Secret:         "test_secret_key",
			ExpirationTime: "86400",
		},
		ServerConfig: config.ServerCfg{
			Port:       "8080",
			Env:        "test",
			GINMode:    "debug",
			Production: false,
		},
	}
}

// setupTestEnvironment initializes the test environment
func setupTestEnvironment(t *testing.T) *TestEnvironment {
	cfg := DefaultTestConfig()

	// Initialize logger
	log := logger.NewLogger(&config.Config{
		DB:     cfg.DBConfig,
		JWT:    cfg.JWTConfig,
		Server: cfg.ServerConfig,
	})

	// Create test lifecycle
	lc := &testLifecycle{}

	// Initialize database
	db, err := database.NewDatabase(lc, &config.Config{
		DB:     cfg.DBConfig,
		JWT:    cfg.JWTConfig,
		Server: cfg.ServerConfig,
	}, log)
	require.NoError(t, err)

	// Ensure test database is clean
	err = db.GetDB().Exec("DROP TABLE IF EXISTS users CASCADE").Error
	require.NoError(t, err)

	// Run migrations
	err = db.GetDB().AutoMigrate(&entity.User{})
	require.NoError(t, err)

	// Create repositories
	userRepo := repository.NewUserRepo(db)

	// Create services
	authService := service.NewAuthService(userRepo, log, &config.Config{
		DB:     cfg.DBConfig,
		JWT:    cfg.JWTConfig,
		Server: cfg.ServerConfig,
	})
	userService := service.NewUserService(userRepo, db)

	// Create handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := gin.Default()
	router.POST("/auth/register", authHandler.Register())
	router.POST("/auth/login", authHandler.Login())
	router.GET("/user/info", userHandler.Info())

	// Cleanup function
	cleanup := func() {
		err := db.GetDB().Exec("DROP TABLE IF EXISTS users CASCADE").Error
		require.NoError(t, err)
		sqlDB, err := db.GetDB().DB()
		require.NoError(t, err)
		sqlDB.Close()
	}

	return &TestEnvironment{
		Router:  router,
		DB:      db,
		Config:  cfg,
		Cleanup: cleanup,
	}
}

// makeRequest helper function to make HTTP requests
func makeRequest(t *testing.T, router *gin.Engine, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	// Set default headers
	req.Header.Set("Content-Type", "application/json")

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// parseResponse helper function to parse HTTP responses
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) *wrapper.Response {
	var resp wrapper.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	return &resp
}

// TestAuthE2E runs the end-to-end authentication tests
func TestAuthE2E(t *testing.T) {
	// Setup test environment
	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Test user
	testUser := DefaultTestUser()

	// Test case 1: Register a new user
	t.Run("Register new user", func(t *testing.T) {
		registerData := request.UserRegisterRequest{
			Email:           testUser.Email,
			Password:        testUser.Password,
			FirstName:       testUser.FirstName,
			LastName:        testUser.LastName,
			ConfirmPassword: testUser.Password,
		}

		w := makeRequest(t, env.Router, "POST", "/auth/register", registerData, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, resp.Data)

		dataMap, ok := resp.Data.(map[string]interface{})
		require.True(t, ok, "Response data is not a map")
		assert.Equal(t, testUser.Email, dataMap["email"])
		assert.Equal(t, testUser.FirstName+" "+testUser.LastName, dataMap["full_name"])
	})

	// Test case 2: Login with registered user
	t.Run("Login with registered user", func(t *testing.T) {
		loginData := request.UserLoginRequest{
			Email:    testUser.Email,
			Password: testUser.Password,
		}

		w := makeRequest(t, env.Router, "POST", "/auth/login", loginData, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, resp.Data)

		dataMap, ok := resp.Data.(map[string]interface{})
		require.True(t, ok, "Response data is not a map")
		assert.NotEmpty(t, dataMap["token"])

		// Store token for next test
		token, ok := dataMap["token"].(string)
		require.True(t, ok, "Token is not a string")
		os.Setenv("TEST_TOKEN", token)
	})

	// Test case 3: Get user info with token
	t.Run("Get user info with token", func(t *testing.T) {
		token := os.Getenv("TEST_TOKEN")
		require.NotEmpty(t, token)

		headers := map[string]string{
			"Authorization": "Bearer " + token,
		}

		w := makeRequest(t, env.Router, "GET", "/user/info", nil, headers)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, resp.Data)

		_, ok := resp.Data.(map[string]interface{})
		require.True(t, ok, "Response data is not a map")
	})

	// Test case 4: Try to login with wrong password
	t.Run("Login with wrong password", func(t *testing.T) {
		loginData := request.UserLoginRequest{
			Email:    testUser.Email,
			Password: "wrongpassword",
		}

		w := makeRequest(t, env.Router, "POST", "/auth/login", loginData, nil)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Nil(t, resp.Data)
		assert.NotNil(t, resp.Message)
	})
}

// testLifecycle implements fx.Lifecycle for testing
type testLifecycle struct {
	hooks []fx.Hook
}

func (l *testLifecycle) Append(hook fx.Hook) {
	l.hooks = append(l.hooks, hook)
}

func (l *testLifecycle) Start(ctx context.Context) error {
	return nil
}

func (l *testLifecycle) Stop(ctx context.Context) error {
	for _, hook := range l.hooks {
		if err := hook.OnStop(ctx); err != nil {
			return err
		}
	}
	return nil
}
