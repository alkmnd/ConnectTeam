package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
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
	if err !=nil {
		return "", err
	}



	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) VerifyPhone(verifyPhone connectteam.VerifyPhone) (string, error) {
	return "1234", nil
}

func (s *AuthService) VerifyUser(verifyUser connectteam.VerifyUser) error {
	return s.repo.VerifyUser(verifyUser)
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}