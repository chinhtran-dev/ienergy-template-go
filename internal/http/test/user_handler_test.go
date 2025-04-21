package handler_test

import (
	"context"
	"encoding/json"
	"ienergy-template-go/internal/http/handler"
	"ienergy-template-go/internal/model/response"
	"ienergy-template-go/pkg/wrapper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserInfo(ctx context.Context) (response.UserInfoResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(response.UserInfoResponse), args.Error(1)
}

func TestUserHandler_Info(t *testing.T) {
	mockUserService := new(MockUserService) // Giả sử bạn có Mock cho UserService
	handler := handler.NewUserHandler(mockUserService)

	// Giả sử mock GetUserInfo trả về một đối tượng hợp lệ
	mockUserService.On("GetUserInfo", mock.Anything).Return(response.UserInfoResponse{
		UserID:   uuid.New(),
		Email:    "test@example.com",
		FullName: "John Doe",
	}, nil).Once()

	// Tạo request và context
	req, _ := http.NewRequest("GET", "/user/info", nil)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	// Gọi phương thức handler.Info()
	handler.Info()(c)

	// Kiểm tra kết quả phản hồi
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Kiểm tra nội dung trả về
	var response wrapper.Response
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
	assert.Equal(t, 200, response.StatusCode)
}
