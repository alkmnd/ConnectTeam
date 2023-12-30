package connectteam
import (
	"errors"
)
 
var ErrNoRecord = errors.New("models: подходящей записи не найдено")
 
type User struct {
	Id int `json:"-" db:"id"`
  	Email string `json:"email" binding "required"` 
  	PhoneNumber string `json:"phone_number" binding "required"`
  	FirstName string `json:"first_name" binding "required"`
  	SecondName string `json:"second_name" binding "required"`
  	Password string `json:"password" binding "required"`
}