package service

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/repository"
	"errors"
	"log"
)

type UserService struct {
	repo repository.UserInterface

}
func NewUserService(repo repository.UserInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserById(id int) (connectteam.UserPublic, error) {
	user, err := s.repo.GetUserById(id)
	return user, err 
}

func (s *UserService) ChangeAccessWithId(id int, access string) (error) {
	if err := s.repo.ChangeAccessWithId(id, access); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUsersList() ([]connectteam.UserPublic, error) {
	return s.repo.GetUsersList()
}

func (s *UserService) ChangePassword(old_password string, new_password string, id int) (error) {
	db_password, err := s.repo.GetPassword(id)
	if err != nil {
		return err
	}

	if db_password != generatePasswordHash(old_password) {
		return errors.New("Wrong old password")
	}

	return s.repo.ChangePassword(generatePasswordHash(new_password), id)
}

func (s *UserService) CheckEmailOnChange(id int, email string, password string) (error) {
	ifEmailExist, err := s.repo.CheckIfExist(email)
	if err != nil {
		return err
	}

	if ifEmailExist {
		return errors.New("Email is already taken")
	}

	db_password, err := s.repo.GetPassword(id)
	if err != nil {
		return err
	}

	if db_password != generatePasswordHash(password) {
		return errors.New("Wrong password")
	}


	


	code, err := CreateVerificationCode(id, email)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	err = s.repo.CreateVerificationCode(id, code)

	if err != nil {
		log.Printf("smtp error: %s", err)
		return errors.New("Error while generating code")
	}

	log.Printf("verification code: %s", code)

	return nil
}

func (s *UserService) ChangeEmail(id int, newEmail string, code string) (error) {
	if newEmail == "" {
		return errors.New("Invalid email")
	}
	db_code, err := s.repo.GetVerificationCode(id)
	if err != nil {
		return errors.New("Verification code is not sent")
	}

	if code != db_code {
		return errors.New("Wrong verification code")
	}

	err = s.repo.DeleteVerificationCode(id, code)
	if err != nil {
		return errors.New("No such row")
	}

	return s.repo.ChangeEmail(newEmail, id)
}

func (s *UserService) DeleteVerificationCode(id int, code string) (error) {
	return s.repo.DeleteVerificationCode(id, code)
}

func (s *UserService) ChangePersonalData(id int, user connectteam.UserPersonalInfo) (error) {
	if user.FirstName != "" {
		err := s.repo.ChangeUserFirstName(id, user.FirstName) 
		if err != nil {
			return err
		}
	}

	if user.SecondName != "" {
		err := s.repo.ChangeUserSecondName(id, user.SecondName) 
		if err != nil {
			return err
		}
	}

	return s.repo.ChangeUserDescription(id, user.Description) 
}