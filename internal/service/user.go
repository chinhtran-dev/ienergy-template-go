package service

import (
	"context"
	"fmt"
	"ienergy-template-go/internal/model/response"
	"ienergy-template-go/internal/repository"
	"ienergy-template-go/pkg/database"
	"ienergy-template-go/pkg/errors"
	"ienergy-template-go/pkg/util"

	"github.com/google/uuid"
)

type UserService interface {
	GetUserInfo(ctx context.Context) (user response.UserInfoResponse, err error)
}

type userService struct {
	userRepo repository.UserRepo
	db       database.Database
}

// GetUserInfo implements IUserService.
func (u *userService) GetUserInfo(ctx context.Context) (user response.UserInfoResponse, err error) {
	userID := util.UserIDFromCTX(ctx)
	if userID == uuid.Nil {
		return user, errors.NewBadRequestError("User ID is not found")
	}
	userEntity, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return user, err
	}

	user = response.UserInfoResponse{
		UserID:   userID,
		Email:    userEntity.Email,
		FullName: fmt.Sprintf("%s %s", userEntity.FirstName, userEntity.LastName),
	}

	return user, nil
}

func NewUserService(
	userRepo repository.UserRepo,
	db database.Database,
) UserService {
	return &userService{
		userRepo: userRepo,
		db:       db,
	}
}
