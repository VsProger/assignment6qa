package models

import (
	"forum/pkg/validator"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Status         int
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type UserSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (u UserSignupForm) FormToUser() User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	return User{
		Name:           u.Name,
		Email:          u.Email,
		HashedPassword: hashedPassword,
	}
}
