package connectteam
import (
	"errors"
)
 
var ErrNoRecord = errors.New("models: подходящей записи не найдено")
 
type User struct {
	Id int `json:"-" db:"id"`
  	Email string `json:"email"` 
  	PhoneNumber string `json:"phone_number"`
  	FirstName string `json:"first_name" binding "required"`
  	SecondName string `json:"second_name" binding "required"`
  	Password string `json:"password" binding "required"`
	Is_verified bool `json:"-" db:"is_verified"`
}

type VerifyPhone struct {
	PhoneNumber string `json:"phone_number" binding required`
}

type VerifyEmail struct {
	Email string `json:"email" binding required`
}

type VerifyUser struct {
	Id int `json:"id,string" binding required`
}