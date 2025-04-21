package repository

import (
	"context"
	"ienergy-template-go/internal/model/entity"
	"ienergy-template-go/pkg/database"
	"ienergy-template-go/pkg/errors"

	logger "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (resp entity.User, error error)
	GetUserByEmail(ctx context.Context, email string) (resp entity.User, error error)
	UserRegister(ctx context.Context, userInfo entity.User) (resp entity.User, error error)
	ValidateUser(userInfo entity.User) (userID uuid.UUID, error error)
	UpdateUser(ctx context.Context, userInfo entity.User) error
	DeleteUser(ctx context.Context, userInfo entity.User) error
	VerifyUserEmail(ctx context.Context, email string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db database.Database) UserRepo {
	return &userRepo{
		db: db.GetDB(),
	}
}

// DeleteUser implements IUserRepo.
func (u *userRepo) DeleteUser(ctx context.Context, userInfo entity.User) error {
	err := u.db.Delete(&entity.User{
		ID: userInfo.ID,
	}).Error
	if err != nil {
		return errors.NewNotFoundError("User not found")
	}
	return nil
}

// GetUserByEmail implements IUserRepo.
func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (resp entity.User, error error) {
	err := u.db.
		Where("email = ?", email).
		First(&resp).Error
	if err != nil {
		return resp, errors.NewInternalServerError("Database error: " + err.Error())
	}
	if resp.ID == uuid.Nil {
		return resp, errors.NewNotFoundError("User not found")
	}
	return
}

// GetUserByID implements IUserRepo.
func (u *userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (resp entity.User, error error) {
	err := u.db.
		Where("id = ?", userID).
		Find(&resp).Error
	if err != nil {
		return resp, errors.NewInternalServerError("Database error: " + err.Error())
	}
	if resp.ID == uuid.Nil {
		return resp, errors.NewNotFoundError("User not found")
	}
	return
}

// UpdateUser implements IUserRepo.
func (u *userRepo) UpdateUser(ctx context.Context, userInfo entity.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternalServerError("Database error: " + err.Error())
	}
	userInfo.Password = string(hashedPassword)
	err = u.db.Save(&userInfo).Error
	if err != nil {
		return errors.NewInternalServerError("Database error: " + err.Error())
	}
	return nil
}

// UserRegister implements IUserRepo.
func (u *userRepo) UserRegister(ctx context.Context, userInfo entity.User) (resp entity.User, error error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return userInfo, errors.NewInternalServerError("Database error: " + err.Error())
	}
	userInfo.Password = string(hashedPassword)
	dbExecute := u.db.Create(&userInfo)
	if dbExecute.Error != nil {
		logger.WithContext(ctx).
			WithField("UserRegister-Input", userInfo).
			WithError(err).
			Error()
		return userInfo, errors.NewInternalServerError("Database error: " + dbExecute.Error.Error())
	}

	return userInfo, nil
}

// ValidateUser implements IUserRepo.
func (u *userRepo) ValidateUser(userInfo entity.User) (userID uuid.UUID, error error) {
	var userInfoDB entity.User
	dbQuery := u.db.
		Where("email = ?", userInfo.Email).
		Find(&userInfoDB)
	if dbQuery.Error != nil {
		return userID, errors.NewInternalServerError("Database error: " + dbQuery.Error.Error())
	}

	if bcrypt.CompareHashAndPassword([]byte(userInfoDB.Password), []byte(userInfo.Password)) == nil {
		userID = userInfoDB.ID
		return
	}

	return userID, errors.NewUnauthorizedError("Invalid email or password")
}

// VerifyUserEmail implements IUserRepo.
func (u *userRepo) VerifyUserEmail(ctx context.Context, email string) error {
	var resp []entity.User
	err := u.db.
		WithContext(ctx).
		Where("email = ?", email).
		Find(&resp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return errors.NewInternalServerError("Database error: " + err.Error())
	}
	if len(resp) != 0 {
		return errors.NewConflictError("Email already exists")
	}
	return nil
}
