package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	"log"
	"math/big"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserById(id uuid.UUID) (connectteam.UserPublic, error) {
	user, err := s.repo.GetUserById(id)
	return user, err
}

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbolBytes    = "!@#$%^&*()_-+=<>?"
	digitBytes     = "0123456789"
	minPasswordLen = 8
)

func generatePassword() (string, error) {

	password := randomCharacter(letterBytes)

	password += randomCharacter(symbolBytes)

	password += randomCharacters(digitBytes + letterBytes + symbolBytes)

	return password, nil
}

func randomCharacter(characters string) string {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
	if err != nil {
		panic(err)
	}
	return string(characters[index.Int64()])
}

func randomCharacters(characters string) string {
	var result string
	for i := 0; i < minPasswordLen-2; i++ {
		result += randomCharacter(characters)
	}
	return result
}
func (s *UserService) RestorePassword(email string) error {

	ifExists, err := s.repo.CheckIfExist(email)
	if err != nil {
		return err
	}

	if !ifExists {
		return errors.New("user with such email is not exist")
	}

	password, err := generatePassword()
	log.Printf("password: %s", password)
	if err != nil {
		return err
	}

	println(email)

	if err := s.repo.UpdatePasswordWithEmail(generatePasswordHash(password), email); err != nil {
		return err
	}

	msg := "Восстановление пароля\n\n" + "Ваш новый пароль: " + password
	if err := SendMessage(email, msg); err != nil {
		return err
	}

	return nil
}

func (s *UserService) RestorePasswordAuthorized(id uuid.UUID) error {
	var userCredentials connectteam.UserCredentials
	userCredentials, err := s.repo.GetUserCredentials(id)
	if err != nil {
		return err
	}

	password, err := generatePassword()
	log.Printf("password: %s", password)
	if err != nil {
		return err
	}

	if err = s.repo.UpdatePasswordWithId(generatePasswordHash(password), id); err != nil {
		return err
	}

	msg := "Восстановление пароля\n\n" + "Ваш новый пароль: " + password
	if err := SendMessage(userCredentials.Email, msg); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateAccessWithId(id uuid.UUID, access connectteam.AccessLevel) error {
	if err := s.repo.UpdateAccessWithId(id, string(access)); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUsersList() ([]connectteam.UserPublic, error) {
	return s.repo.GetUsersList()
}

func (s *UserService) UpdatePassword(oldPassword string, newPassword string, id uuid.UUID) error {
	dbPassword, err := s.repo.GetPassword(id)
	if err != nil {
		return err
	}

	if dbPassword != generatePasswordHash(oldPassword) {
		return errors.New("wrong old password")
	}

	return s.repo.UpdatePasswordWithId(generatePasswordHash(newPassword), id)
}

func (s *UserService) CheckEmailOnChange(id uuid.UUID, email string, password string) error {
	ifEmailExist, err := s.repo.CheckIfExist(email)
	if err != nil {
		return err
	}

	if ifEmailExist {
		return errors.New("email is already taken")
	}

	dbPassword, err := s.repo.GetPassword(id)
	if err != nil {
		return errors.New("invalid password")
	}

	if dbPassword != generatePasswordHash(password) {

		return errors.New("wrong password")

	}

	code, err := CreateVerificationCode(email)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	err = s.repo.CreateVerificationCode(email, code)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return errors.New("error while generating code")
	}

	log.Printf("verification code: %s", code)

	return nil
}

func (s *UserService) UpdateEmail(id uuid.UUID, newEmail string, code string) error {
	if newEmail == "" {
		return errors.New("invalid email")
	}
	dbCode, err := s.repo.GetVerificationCode(newEmail)
	if err != nil {
		return errors.New("verification code is not sent")
	}

	if code != dbCode {
		return errors.New("wrong verification code")
	}

	err = s.repo.DeleteVerificationCode(newEmail, code)
	if err != nil {
		return errors.New("no such row")
	}

	return s.repo.UpdateEmail(newEmail, id)
}

func (s *UserService) DeleteVerificationCode(email string, code string) error {
	return s.repo.DeleteVerificationCode(email, code)
}

func (s *UserService) UpdatePersonalData(id uuid.UUID, user connectteam.UserPersonalInfo) error {
	if user.FirstName != "" {
		err := s.repo.UpdateUserFirstName(id, user.FirstName)
		if err != nil {
			return err
		}
	}

	if user.SecondName != "" {
		err := s.repo.UpdateUserSecondName(id, user.SecondName)
		if err != nil {
			return err
		}
	}

	return s.repo.UpdateUserDescription(id, user.Description)
}

func (s *UserService) UpdateCompanyData(id uuid.UUID, company connectteam.UserCompanyData) error {
	err := s.repo.UpdateCompanyName(id, company.CompanyName)
	if err != nil {
		return err
	}

	err = s.repo.UpdateCompanyInfo(id, company.CompanyInfo)
	if err != nil {
		return err
	}

	err = s.repo.UpdateCompanyURL(id, company.CompanyURL)
	if err != nil {
		return err
	}

	return err
}

func (s *UserService) GetUserPlan(userId uuid.UUID) (connectteam.UserPlan, error) {
	return s.repo.GetUserPlan(userId)
}
