package connectteam
import (
	"errors"
)
 
var ErrNoRecord = errors.New("models: подходящей записи не найдено")
 
type User struct {
	Id int `json:"id" db:"id"`
  	Email string `json:"email" db:"email"` 
  	// PhoneNumber string `json:"phone_number" db:"phone_number"`
  	FirstName string `json:"first_name" db:"first_name" binding "required"`
  	SecondName string `json:"second_name" db:"second_name" binding "required"`
  	Password string `json:"password" binding "required"`
	Is_verified bool `json:"-" db:"is_verified"`
	Access string `json:"access" db:"access"`
}

type UserPublic struct {
	Id int `json:"id" db:"id"`
	Email string `json:"email" db:"email"` 
	// PhoneNumber string `json:"phone_number" db:"phone_number"`
	FirstName string `json:"first_name" db:"first_name" binding "required"`
	SecondName string `json:"second_name" db:"second_name" binding "required"`
	Description string `json:"description" db:"description"`
  	Access string `json:"access" db:"access"`
	CompanyName string `json:"company_name" db:"company_name"`
	CompanyInfo string `json:"company_info" db:"company_info"`
	CompanyURL string `json:"company_url" db:"company_url"`
	ProfileImage string `json:"profile_image" db:"profile_image"`
}

type UserPersonalInfo struct {
	FirstName string `json:"first_name" db:"first_name" binding "required"`
	SecondName string `json:"second_name" db:"second_name" binding "required"`
	Description string `json:"description" db:"description"`
}
type UserCompanyData struct {
	CompanyName string `json:"company_name" db:"company_name"`
	CompanyInfo string `json:"company_info" db:"company_info"`
	CompanyURL string `json:"company_url" db:"company_url"`
}

type VerifyPhone struct {
	PhoneNumber string `json:"phone_number" binding required`
}

type VerifyEmail struct {
	Email string `json:"email" binding "required"`
}

type VerifyUser struct {
	Id int `json:"id,string" binding "required"`
	Code string `json:"code" binding "required"`
}

