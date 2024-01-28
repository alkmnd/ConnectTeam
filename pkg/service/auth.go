package service

import (
	"ConnectTeam"
	"ConnectTeam/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)
const (
	salt = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
	Role string `json:"access"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user connectteam.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(login, password string, isEmail bool) (string, error) {
	var user connectteam.User
	var err error 
	if isEmail {
		user, err = s.repo.GetUserWithEmail(login, generatePasswordHash(password))
	} else {
		user, err = s.repo.GetUserWithPhone(login, generatePasswordHash(password))
	}
	if err != nil {
		return "", errors.New("Invalid login data")
	}
	if !user.Is_verified {
		println("meow")
		return "", errors.New("User is not verified")
	}

	if err != nil {
		return "", err
	}



	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		user.Access,
	})
	return token.SignedString([]byte(signingKey))
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

func (s *AuthService) VerifyPhone(verifyPhone connectteam.VerifyPhone) (string, error) {
	return "1234", nil
}




func (s *AuthService) VerifyEmail(verifyEmail connectteam.VerifyEmail) (int, error) {
	from := os.Getenv("EMAIL")
   	password := os.Getenv("EMAIL_PASSWORD")
	to := verifyEmail.Email
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
		return 0, errors.New("The recipient address is not a valid")
	}


	id, err := s.repo.GetIdWithEmail(verifyEmail.Email)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return 0, errors.New("No user with such email")
	}

	err = s.repo.CreateVerificationCode(id, confirmationCode)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return 0, errors.New("Error while generating code")
	}

	log.Printf("verification code: %s", confirmationCode)


   return id, err
}

func (s *AuthService) VerifyUser(verifyUser connectteam.VerifyUser) error {
	code, err := s.repo.GetVerificationCode(verifyUser.Id)
	if err != nil {
		return errors.New("Wrong verification code")
	}

	if code != verifyUser.Code {

		return errors.New("Wrong verification code")
	}
	
	return s.repo.VerifyUser(verifyUser)
}

func (s *AuthService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("Token claims are not of type *tokenClaims")
	}

	return claims.UserId, claims.Role, nil
}