package integration

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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

type testEnvironment struct {
	router  *gin.Engine
	db      database.Database
	cleanup func()
}

func setupTestEnvironment(t *testing.T) *testEnvironment {
	// Create test configuration
	cfg := &config.Config{
		DB: config.DBConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "Admin",
			Password: "admin@123",
			DBName:   "ienergy_db",
			SSLMode:  "disable",
		},
		JWT: config.JWTConfig{
			Secret:         "test_secret_key",
			ExpirationTime: "86400", // 24 hours in seconds
		},
		Server: config.ServerCfg{
			Port:       "8080",
			Env:        "test",
			GINMode:    "debug",
			Production: false,
		},
	}

	// Initialize logger
	log := logger.NewLogger(cfg)

	// Create a test lifecycle
	lc := &testLifecycle{}

	// Initialize database
	db, err := database.NewDatabase(lc, cfg, log)
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
	authService := service.NewAuthService(userRepo, log, cfg)
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
		// Clean up test database
		err := db.GetDB().Exec("DROP TABLE IF EXISTS users CASCADE").Error
		require.NoError(t, err)
		sqlDB, err := db.GetDB().DB()
		require.NoError(t, err)
		sqlDB.Close()
	}

	return &testEnvironment{
		router:  router,
		db:      db,
		cleanup: cleanup,
	}
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

func TestAuthIntegration(t *testing.T) {
	// Setup test environment
	env := setupTestEnvironment(t)
	defer env.cleanup()

	// Define test cases
	testCases := []struct {
		name         string
		method       string
		path         string
		body         interface{}
		expectedCode int
		validateResp func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "successful registration",
			method: "POST",
			path:   "/auth/register",
			body: request.UserRegisterRequest{
				Email:           "test@example.com",
				Password:        "password1234",
				FirstName:       "John",
				LastName:        "Doe",
				ConfirmPassword: "password1234",
			},
			expectedCode: http.StatusOK,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.NotNil(t, resp.Data)

				// Convert data to map to access fields
				dataMap, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "test@example.com", dataMap["email"])
				assert.Equal(t, "John Doe", dataMap["full_name"])
			},
		},
		{
			name:   "successful login",
			method: "POST",
			path:   "/auth/login",
			body: request.UserLoginRequest{
				Email:    "test@example.com",
				Password: "password1234",
			},
			expectedCode: http.StatusOK,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.NotNil(t, resp.Data)

				// Convert data to map to access fields
				dataMap, ok := resp.Data.(map[string]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, dataMap["token"])
			},
		},
		{
			name:   "invalid credentials login",
			method: "POST",
			path:   "/auth/login",
			body: request.UserLoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			expectedCode: http.StatusUnauthorized,
			validateResp: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp wrapper.Response
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
				assert.Nil(t, resp.Data)
				assert.NotNil(t, resp.Message)
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Create request body
			var reqBody []byte
			if tc.body != nil {
				var err error
				reqBody, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}

			// Create request
			req, err := http.NewRequest(tc.method, tc.path, bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// If this is a user info request, we need to add the token from the login response
			if tc.path == "/user/info" {
				// First login to get token
				loginReq := request.UserLoginRequest{
					Email:    "test@example.com",
					Password: "password1234",
				}
				loginBody, err := json.Marshal(loginReq)
				require.NoError(t, err)

				loginHTTPReq, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBody))
				require.NoError(t, err)
				loginHTTPReq.Header.Set("Content-Type", "application/json")

				loginW := httptest.NewRecorder()
				env.router.ServeHTTP(loginW, loginHTTPReq)

				var loginResp wrapper.Response
				err = json.Unmarshal(loginW.Body.Bytes(), &loginResp)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, loginResp.StatusCode)
				require.NotNil(t, loginResp.Data)

				dataMap, ok := loginResp.Data.(map[string]interface{})
				require.True(t, ok)
				require.NotNil(t, dataMap["token"])

				token, ok := dataMap["token"].(string)
				require.True(t, ok)
				require.NotEmpty(t, token)

				// Add token to request
				req.Header.Set("Authorization", "Bearer "+token)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute request
			env.router.ServeHTTP(w, req)

			// Validate response
			assert.Equal(t, tc.expectedCode, w.Code)
			tc.validateResp(t, w)
		})
	}
}
