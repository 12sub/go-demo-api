package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Creating User Structure
type User struct{
	gorm.Model
	Username	string	`gorm:"unique"`
	Email	string	`gorm:"unique"`
	Password	string	
	Role	string	

}

// Hash password before starting
func (u *User) HashPassword() error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashPassword)
	return nil
}

// Compare hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(u.Password))
	return err == nil
}