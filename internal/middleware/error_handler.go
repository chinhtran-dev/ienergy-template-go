package middleware

import (
	"ienergy-template-go/pkg/errors"
	"ienergy-template-go/pkg/logger"
	"ienergy-template-go/pkg/wrapper"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

// ErrorHandler handles application errors
type ErrorHandler struct {
	logger *logger.StandardLogger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *logger.StandardLogger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// Handle is the main error handling middleware
func (h *ErrorHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				h.logger.Error("panic recovered", "error", err, "stack", string(debug.Stack()))
				c.JSON(http.StatusInternalServerError, wrapper.NewResponse(
					http.StatusInternalServerError,
					0,
					nil,
					"internal server error",
				))
			}
		}()

		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()
			if lastErr == nil {
				return
			}
			err := lastErr.Err

			h.logger.Error("request error", "error", err)

			// Handle different types of errors
			switch e := err.(type) {
			case *errors.AppError:
				c.JSON(e.Status, wrapper.NewResponse(
					e.Status,
					0,
					nil,
					e.Message,
				))
			case validator.ValidationErrors:
				if len(e) > 0 {
					c.JSON(http.StatusBadRequest, wrapper.NewResponse(
						http.StatusBadRequest,
						0,
						nil,
						"validation error",
					))
				} else {
					c.JSON(http.StatusBadRequest, wrapper.NewResponse(
						http.StatusBadRequest,
						0,
						nil,
						"invalid request",
					))
				}
			case *pq.Error:
				h.handleDatabaseError(c, e)
			default:
				c.JSON(http.StatusInternalServerError, wrapper.NewResponse(
					http.StatusInternalServerError,
					0,
					nil,
					err.Error(),
				))
			}
		}
	}
}

// handleDatabaseError handles database specific errors
func (h *ErrorHandler) handleDatabaseError(c *gin.Context, err *pq.Error) {
	if err == nil {
		c.JSON(http.StatusInternalServerError, wrapper.NewResponse(
			http.StatusInternalServerError,
			0,
			nil,
			"database error",
		))
		return
	}

	switch err.Code {
	case "23505": // unique_violation
		c.JSON(http.StatusConflict, wrapper.NewResponse(
			http.StatusConflict,
			0,
			nil,
			"resource already exists",
		))
	case "23503": // foreign_key_violation
		c.JSON(http.StatusBadRequest, wrapper.NewResponse(
			http.StatusBadRequest,
			0,
			nil,
			"invalid reference",
		))
	default:
		h.logger.Error("database error", "error", err)
		c.JSON(http.StatusInternalServerError, wrapper.NewResponse(
			http.StatusInternalServerError,
			0,
			nil,
			"database error",
		))
	}
}
