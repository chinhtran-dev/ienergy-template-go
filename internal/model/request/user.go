package request

import (
	"errors"
	"strings"
)

type UserRegisterRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (u *UserRegisterRequest) Validate() error {
	if u.Password != u.ConfirmPassword {
		return errors.New("Password and confirm password are not meet!") //nolint
	}
	if len(u.Password) > 150 && len(u.Password) < 8 {
		return errors.New("Password must be at least 8 characters") //nolint
	}
	if len(u.Email) == 0 {
		return errors.New("email is required!") //nolint
	}
	if len(u.FirstName) == 0 {
		return errors.New("First name is required!") //nolint
	}
	if len(u.LastName) == 0 {
		return errors.New("Last name is required!") //nolint
	}
	emailSplited := strings.Split(u.Email, "@")
	if len(emailSplited) != 2 {
		return errors.New("Invalid email address!") //nolint
	}

	return nil
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserLoginRequest) Validate() error {
	if len(u.Email) == 0 {
		return errors.New("email is required!") //nolint
	}
	emailSplited := strings.Split(u.Email, "@")
	if len(emailSplited) != 2 {
		return errors.New("Invalid email address!") //nolint
	}
	if len(u.Password) > 150 || len(u.Password) < 8 {
		return errors.New("Password must be at least 8 characters") //nolint
	}

	return nil
}
