package service

import (
	"ConnectTeam"
	"ConnectTeam/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
	Role   string    `json:"access"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user connectteam.UserSignUpRequest) (uuid.UUID, error) {
	user.Password = generatePasswordHash(user.Password)
	dbCode, err := s.repo.GetVerificationCode(user.Email)
	if err != nil {
		return uuid.Nil, errors.New("wrong verification code")
	}
	if dbCode != user.VerificationCode {
		return uuid.Nil, errors.New("wrong verification code")
	}
	repoUser := connectteam.User{
		Email:      user.Email,
		Password:   user.Password,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Access:     string(connectteam.UserAccess),
	}
	return s.repo.CreateUser(repoUser)
}

func (s *AuthService) GenerateToken(login, password string, isEmail bool) (string, string, error) {
	var user connectteam.User
	var err error
	if isEmail {
		user, err = s.repo.GetUserWithEmail(login, generatePasswordHash(password))
	} else {
		return "", "", nil
	}
	if err != nil {
		return "", "", errors.New("invalid login data")
	}

	if err != nil {
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		user.Access,
	})
	signedString, err := token.SignedString([]byte(signingKey))

	if err != nil {
		return "", "", err
	}

	return user.Access, signedString, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func generateConfirmationCode() string {
	code := ""

	for i := 0; i < 4; i++ {
		code += fmt.Sprint(randomCrypto())
	}

	return code

}

func CreateVerificationCode(email string) (string, error) {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	to := email
	confirmationCode := generateConfirmationCode()

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Код подтверждения\n\n" +
		confirmationCode
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return "", errors.New("the recipient address is not a valid")
	}
	return confirmationCode, nil

}

func (s *AuthService) VerifyEmail(verifyEmail connectteam.VerifyEmail) error {

	//id, err := s.repo.GetIdWithEmail(verifyEmail.Email)
	//
	//if err != nil {
	//	log.Printf("smtp error: %s", err)
	//	return 0, errors.New("no user with such email")
	//}

	confirmationCode, err := CreateVerificationCode(verifyEmail.Email)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	err = s.repo.CreateVerificationCode(verifyEmail.Email, confirmationCode)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return errors.New("error while generating code")
	}

	log.Printf("verification code: %s", confirmationCode)

	return err
}

//func (s *AuthService) VerifyUser(verifyUser connectteam.VerifyUser) error {
//	code, err := s.repo.GetVerificationCode(verifyUser.Id)
//	if err != nil {
//		return errors.New("wrong verification code")
//	}
//
//	if code != verifyUser.Code {
//
//		return errors.New("wrong verification code")
//	}
//
//	err = s.repo.DeleteVerificationCode(verifyUser.Id, verifyUser.Code)
//	if err != nil {
//		return errors.New("no such row")
//	}
//
//	return s.repo.VerifyUser(verifyUser)
//}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.Nil, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, claims.Role, nil
}

func (s *AuthService) DeleteVerificationCode(email string, code string) error {
	return s.repo.DeleteVerificationCode(email, code)
}
