package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type User struct {
	Id         uuid.UUID `json:"id" db:"id" swagger:"_"`
	Email      string    `json:"email" db:"email" binding:"required"`
	FirstName  string    `json:"first_name" db:"first_name" binding:"required"`
	SecondName string    `json:"second_name" db:"second_name" binding:"required"`
	Password   string    `json:"password" binding:"required"`
	Access     string    `json:"access" db:"access" swagger:"_"`
}

type UserSignUpRequest struct {
	Id               uuid.UUID `json:"id" db:"id" swagger:"_"`
	Email            string    `json:"email" db:"email" binding:"required"`
	FirstName        string    `json:"first_name" db:"first_name" binding:"required"`
	SecondName       string    `json:"second_name" db:"second_name" binding:"required"`
	Password         string    `json:"password" binding:"required"`
	VerificationCode string    `json:"verification_code" binding:"required"`
	Access           string    `json:"access" db:"access" swagger:"_"`
}

type UserCredentials struct {
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
}

type UserResult struct {
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	Value  uuid.UUID `json:"value" db:"value"`
}

type UserPublic struct {
	Id    uuid.UUID `json:"id" db:"id"`
	Email string    `json:"email" db:"email"`
	// PhoneNumber string `json:"phone_number" db:"phone_number"`
	FirstName    string `json:"first_name" db:"first_name" binding:"required"`
	SecondName   string `json:"second_name" db:"second_name" binding:"required"`
	Description  string `json:"description" db:"description"`
	Access       string `json:"access" db:"access"`
	CompanyName  string `json:"company_name" db:"company_name"`
	CompanyInfo  string `json:"company_info" db:"company_info"`
	CompanyURL   string `json:"company_url" db:"company_url"`
	CompanyLogo  string `json:"company_logo" db:"company_logo"`
	ProfileImage string `json:"profile_image" db:"profile_image"`
}

type UserPersonalInfo struct {
	FirstName   string `json:"first_name" db:"first_name" binding:"required"`
	SecondName  string `json:"second_name" db:"second_name" binding:"required"`
	Description string `json:"description" db:"description"`
}
type UserCompanyData struct {
	CompanyName string `json:"company_name" db:"company_name"`
	CompanyInfo string `json:"company_info" db:"company_info"`
	CompanyURL  string `json:"company_url" db:"company_url"`
}

type PlanType string

func (pt *PlanType) Scan(value interface{}) error {
	if value == nil {
		*pt = ""
		return nil
	}
	stringValue, ok := value.([]byte)
	if !ok {
		return errors.New("unexpected type for PlanType")
	}
	*pt = PlanType(string(stringValue))
	return nil
}

type UserPlan struct {
	Id             uuid.UUID `json:"id,omitempty" db:"id"`
	PlanType       string    `json:"plan_type,omitempty" db:"plan_type"`
	HolderId       uuid.UUID `json:"holder_id,omitempty" db:"holder_id"`
	UserId         uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	ExpiryDate     time.Time `json:"expiry_date,omitempty" db:"expiry_date"`
	Duration       int       `json:"duration,omitempty" db:"duration"`
	PlanAccess     string    `json:"access,omitempty" db:"access"`
	Status         string    `json:"status,omitempty" db:"status"`
	InvitationCode string    `json:"invitation_code,omitempty" db:"invitation_code"`
	IsTrial        bool      `json:"is_trial,omitempty" db:"is_trial"`
}

type Subscription struct {
	Id             uuid.UUID `json:"id,omitempty" db:"id"`
	PlanType       string    `json:"plan_type,omitempty" db:"plan_type"`
	HolderId       uuid.UUID `json:"holder_id,omitempty" db:"holder_id"`
	UserId         uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	ExpiryDate     time.Time `json:"expiry_date,omitempty" db:"expiry_date"`
	Duration       int       `json:"duration,omitempty" db:"duration"`
	Status         string    `json:"status,omitempty" db:"status"`
	InvitationCode string    `json:"invitation_code,omitempty" db:"invitation_code"`
	IsTrial        bool      `json:"is_trial,omitempty" db:"is_trial"`
}

const (
	OnConfirm = "on_confirm"
	Active    = "active"
	Expired   = "expired"
	Trial     = "trial"
)

type VerifyPhone struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyEmail struct {
	Email string `json:"email" binding:"required"`
}

type VerifyUser struct {
	Id   uuid.UUID `json:"id,string" binding:"required"`
	Code string    `json:"code" binding:"required"`
}

type PlanRequest struct {
	Id          uuid.UUID     `json:"-" db:"id"`
	UserId      int           `json:"-" db:"user_id"`
	Duration    time.Duration `json:"duration" db:"duration"`
	RequestDate time.Time     `json:"-" db:"request_date"`
	PlanType    string        `json:"plan_type" db:"plan_type"`
}

type AccessLevel string

const (
	Admin      AccessLevel = "admin"
	SuperAdmin AccessLevel = "super_admin"
	UserAccess AccessLevel = "user"
)

type Topic struct {
	Id    uuid.UUID `json:"id" db:"id"`
	Title string    `json:"title" db:"title" binding:"required"`
}

type Question struct {
	Id      uuid.UUID `json:"id" db:"id"`
	TopicId uuid.UUID `json:"topic_id" db:"topic_id"`
	Content string    `json:"content" db:"content"`
}

type Game struct {
	Id             uuid.UUID `json:"id" db:"id"`
	CreatorId      uuid.UUID `json:"creator_id" db:"creator_id"`
	InvitationCode string    `json:"invitation_code" db:"invitation_code"`
	Name           string    `json:"name" db:"name"`
	StartDate      time.Time `json:"start_date" db:"start_date"`
	Status         string    `json:"status" db:"status"`
	MaxSize        int       `json:"max_size" db:"max_size"`
}
