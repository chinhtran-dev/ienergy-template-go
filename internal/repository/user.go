package repository

import (
	"context"
	"errors"
	"fmt"
	"ienergy-template-go/internal/model/entity"
	"ienergy-template-go/pkg/database"

	logger "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (resp entity.User, err error)
	GetUserByEmail(ctx context.Context, email string) (resp entity.User, err error)
	UserRegister(ctx context.Context, userInfo entity.User) (resp entity.User, err error)
	ValidateUser(userInfo entity.User) (userID uuid.UUID, err error)
	UpdateUser(ctx context.Context, userInfo entity.User) error
	DeleteUser(ctx context.Context, userInfo entity.User) error
	VerifyUserEmail(ctx context.Context, email string) (err error)
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
	return u.db.Delete(&entity.User{
		ID: userInfo.ID,
	}).Error
}

// GetUserByEmail implements IUserRepo.
func (u *userRepo) GetUserByEmail(ctx context.Context, email string) (resp entity.User, err error) {
	err = u.db.
		Where("email = ?", email).
		First(&resp).Error
	if resp.ID == uuid.Nil {
		return resp, fmt.Errorf("Email address not found: %v", email)
	}
	return
}

// GetUserByID implements IUserRepo.
func (u *userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (resp entity.User, err error) {
	err = u.db.
		Where("id = ?", userID).
		Find(&resp).Error
	return
}

// UpdateUser implements IUserRepo.
func (u *userRepo) UpdateUser(ctx context.Context, userInfo entity.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	userInfo.Password = string(hashedPassword)
	return u.db.Save(&userInfo).Error
}

// UserRegister implements IUserRepo.
func (u *userRepo) UserRegister(ctx context.Context, userInfo entity.User) (resp entity.User, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	userInfo.Password = string(hashedPassword)
	dbExecute := u.db.Create(&userInfo)
	if dbExecute.Error != nil {
		logger.WithContext(ctx).
			WithField("UserRegister-Input", userInfo).
			WithError(err).
			Error()
		return userInfo, dbExecute.Error
	}

	return userInfo, nil
}

// ValidateUser implements IUserRepo.
func (u *userRepo) ValidateUser(userInfo entity.User) (userID uuid.UUID, err error) {
	var userInfoDB entity.User
	dbQuery := u.db.
		Where("email = ?", userInfo.Email).
		Find(&userInfoDB)
	if dbQuery.Error != nil {
		err = dbQuery.Error
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(userInfoDB.Password), []byte(userInfo.Password)) == nil {
		userID = userInfoDB.ID
		return
	}

	return
}

// VerifyUserEmail implements IUserRepo.
func (u *userRepo) VerifyUserEmail(ctx context.Context, email string) (err error) {
	var resp []entity.User
	err = u.db.
		Where("email = ?", email).
		Find(&resp).Error
	if err != nil {
		if err == gorm.ErrEmptySlice {
			return nil
		}
		return
	}
	if len(resp) != 0 {
		return errors.New("Email already exists")
	}
	return
}
