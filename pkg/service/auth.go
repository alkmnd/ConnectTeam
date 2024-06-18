package service

import (
	"ConnectTeam/models"
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
	salt            = "hjqrhjqw124617ajfhajs"
	signingKey      = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL        = 24 * time.Hour
	refreshTokenTTL = 24 * time.Hour * 90
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
	Role   string    `json:"access"`
}

type refreshTokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.UserSignUpRequest) (uuid.UUID, error) {
	user.Password = s.GeneratePasswordHash(user.Password)
	dbCode, err := s.repo.GetVerificationCode(user.Email)
	if err != nil {
		return uuid.Nil, errors.New("wrong verification code")
	}
	if dbCode != user.VerificationCode {
		return uuid.Nil, errors.New("wrong verification code")
	}
	repoUser := models.User{
		Email:      user.Email,
		Password:   user.Password,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Access:     string(models.UserAccess),
	}
	return s.repo.CreateUser(repoUser)
}

// GenerateAccessToken generate token and returns user id, access and token.
func (s *AuthService) GenerateAccessToken(login, passwordHash string) (string, string, error, uuid.UUID) {
	var user models.User
	var err error

	// Check auth data.
	user, err = s.repo.GetUserWithEmail(login, passwordHash)
	if err != nil {
		return "", "", errors.New("invalid login data"), uuid.Nil
	}

	// Generate token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		user.Access,
	})

	// Sign token with signing key.
	signedString, err := token.SignedString([]byte(signingKey))

	if err != nil {
		return "", "", err, uuid.Nil
	}

	return user.Access, signedString, nil, user.Id
}

// GenerateRefreshToken generates refresh token and returns generated token.
func (s *AuthService) GenerateRefreshToken(userId uuid.UUID) (string, error) {

	var err error

	// Create a token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: userId,
	})

	// Sign token.
	signedString, err := token.SignedString([]byte(signingKey))

	if err != nil {
		return "", err
	}

	return signedString, nil
}

func (s *AuthService) GeneratePasswordHash(password string) string {
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

func (s *AuthService) VerifyEmail(verifyEmail models.VerifyEmail) error {

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

// ParseAccessToken parse token and returns user id and user access.
func (s *AuthService) ParseAccessToken(accessToken string) (uuid.UUID, string, error) {
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

// ParseRefreshToken parses token and returns user id.
func (s *AuthService) ParseRefreshToken(refreshToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*refreshTokenClaims)
	if !ok {
		return uuid.Nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, err
}

func (s *AuthService) DeleteVerificationCode(email string, code string) error {
	return s.repo.DeleteVerificationCode(email, code)
}
