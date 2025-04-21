# Testing Guideline

## Table of Contents
1. [Testing Principles](#testing-principles)
2. [Testing Types](#testing-types)
3. [Test Structure](#test-structure)
4. [Test Naming Conventions](#test-naming-conventions)
5. [Test Organization](#test-organization)
6. [Test Coverage](#test-coverage)
7. [Best Practices](#best-practices)
8. [Common Testing Patterns](#common-testing-patterns)
9. [Testing Tools](#testing-tools)

## Testing Principles

1. **Test Early, Test Often**: Write tests alongside code development
2. **Test Independence**: Each test should be independent and not rely on other tests
3. **Test Readability**: Tests should be easy to read and understand
4. **Test Maintainability**: Tests should be easy to maintain and update
5. **Test Coverage**: Aim for meaningful coverage, not just percentage

## Testing Types

### 1. Unit Tests
- Test individual functions and methods
- Isolate dependencies using mocks
- Focus on business logic
- Location: `internal/*/test/`

### 2. Integration Tests
- Test interaction between components
- Use real database connections
- Test API endpoints
- Location: `tests/integration/`

### 3. HTTP Tests
- Test HTTP handlers and middleware
- Use `httptest` package
- Test request/response flow
- Location: `internal/http/test/`

## Test Structure

### Basic Test Structure
```go
func TestFunctionName(t *testing.T) {
    // Setup
    t.Parallel()
    
    // Test cases
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "success case",
            input:   InputType{},
            want:    OutputType{},
            wantErr: false,
        },
        {
            name:    "error case",
            input:   InputType{},
            want:    OutputType{},
            wantErr: true,
        },
    }

    // Run tests
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### HTTP Test Structure
```go
func TestHandlerName(t *testing.T) {
    // Setup
    t.Parallel()
    router := setupRouter()
    
    // Test cases
    tests := []struct {
        name       string
        method     string
        path       string
        body       interface{}
        wantStatus int
        wantBody   string
    }{
        {
            name:       "success case",
            method:     "POST",
            path:       "/api/v1/endpoint",
            body:       requestBody{},
            wantStatus: http.StatusOK,
            wantBody:   `{"message":"success"}`,
        },
    }

    // Run tests
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Test Naming Conventions

1. **Test Function Names**:
   - Format: `Test[FunctionName]`
   - Example: `TestUserService_Create`

2. **Test Case Names**:
   - Format: `[scenario]_[expected_result]`
   - Example: `valid_input_returns_success`

3. **Test File Names**:
   - Format: `[original_file]_test.go`
   - Example: `user_service_test.go`

## Test Organization

### Directory Structure
```
.
├── internal/
│   ├── service/
│   │   └── test/
│   │       └── user_service_test.go
│   ├── repository/
│   │   └── test/
│   │       └── user_repository_test.go
│   └── http/
│       └── test/
│           └── user_handler_test.go
└── tests/
    └── integration/
        └── user_integration_test.go
```

### Test Categories
1. **Happy Path Tests**: Test successful scenarios
2. **Error Path Tests**: Test error handling
3. **Edge Case Tests**: Test boundary conditions
4. **Integration Tests**: Test component interactions

## Test Coverage

### Coverage Goals
- Unit Tests: 80% minimum
- Integration Tests: 60% minimum
- Critical Paths: 100%

### Running Coverage
```bash
# Run tests with coverage
go test -v -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out -o coverage.html
```

## Best Practices

1. **Use Table-Driven Tests**
   - Group related test cases
   - Easy to add new cases
   - Clear test structure

2. **Mock External Dependencies**
   - Use interfaces for dependencies
   - Create mock implementations
   - Control test environment

3. **Clean Test Data**
   - Use test fixtures
   - Clean up after tests
   - Use transactions for database tests

4. **Test Error Cases**
   - Test invalid inputs
   - Test error conditions
   - Test error messages

5. **Use Test Helpers**
   - Create helper functions
   - Reduce code duplication
   - Improve test readability

## Common Testing Patterns

### 1. Database Testing
```go
func TestWithDB(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()
    
    // Test implementation
}
```

### 2. HTTP Testing
```go
func TestHTTPHandler(t *testing.T) {
    router := setupTestRouter()
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/path", nil)
    
    router.ServeHTTP(w, req)
    
    // Assertions
}
```

### 3. Mock Testing
```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) FindByID(id string) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}
```

## Testing Tools

1. **Standard Library**
   - `testing`: Core testing package
   - `httptest`: HTTP testing utilities
   - `testify`: Testing utilities and assertions

2. **Third-Party Tools**
   - `mockery`: Generate mock implementations
   - `golangci-lint`: Code quality and style checking
   - `go-cmp`: Advanced comparison utilities

3. **IDE Integration**
   - VS Code Go Test Explorer
   - GoLand Test Runner
   - Debugging test cases

## Example Test Implementation

### Unit Test Example
```go
func TestUserService_Create(t *testing.T) {
    t.Parallel()
    
    tests := []struct {
        name    string
        user    *User
        wantErr bool
    }{
        {
            name: "success",
            user: &User{
                Email:    "test@example.com",
                Password: "password123",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            user: &User{
                Email:    "invalid-email",
                Password: "password123",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            mockRepo := new(MockRepository)
            service := NewUserService(mockRepo)
            
            // Test
            err := service.Create(tt.user)
            
            // Assert
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### HTTP Test Example
```go
func TestUserHandler_Create(t *testing.T) {
    t.Parallel()
    
    tests := []struct {
        name       string
        body       string
        wantStatus int
        wantBody   string
    }{
        {
            name:       "success",
            body:       `{"email":"test@example.com","password":"password123"}`,
            wantStatus: http.StatusCreated,
            wantBody:   `{"id":"123","email":"test@example.com"}`,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            router := setupTestRouter()
            w := httptest.NewRecorder()
            req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(tt.body))
            
            // Test
            router.ServeHTTP(w, req)
            
            // Assert
            assert.Equal(t, tt.wantStatus, w.Code)
            assert.JSONEq(t, tt.wantBody, w.Body.String())
        })
    }
} 