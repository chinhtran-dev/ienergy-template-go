package entity

import (
	"ienergy-template-go/internal/model/request"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName string    `gorm:"column:first_name;type:varchar(50)"`
	LastName  string    `gorm:"column:last_name;type:varchar(50)"`
	Email     string    `gorm:"column:email;type:varchar(50);index:email_idx,unique"`
	Password  string    `gorm:"column:password;type:varchar(150)"`
	BaseEntity
}

func (e *User) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}

func ToEntityModel(input request.UserRegisterRequest) User {
	return User{
		ID:        uuid.UUID{},
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
		BaseEntity: BaseEntity{
			CreatedBy: input.Email,
		}}
}
