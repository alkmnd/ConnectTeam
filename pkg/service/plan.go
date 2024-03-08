package service

import (
	"ConnectTeam"
	"ConnectTeam/pkg/repository"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

type PlanService struct {
	repo repository.Plan
}

func NewPlanService(repo repository.Plan) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) GetUserActivePlan(userId int) (connectteam.UserPlan, error) {
	println(userId)
	return s.repo.GetUserActivePlan(userId)
}

func (s *PlanService) CheckIfSubscriptionExists(userId int) (bool, error) {
	var userPlans []connectteam.UserPlan

	userPlans, err := s.repo.GetUserSubscriptions(userId)
	if err != nil {
		return false, err
	}

	return len(userPlans) > 0, nil
}

func (s *PlanService) GetUserSubscriptions(userId int) ([]connectteam.UserPlan, error) {
	var userPlans []connectteam.UserPlan

	userPlans, err := s.repo.GetUserSubscriptions(userId)
	if err != nil {
		return userPlans, err
	}

	return userPlans, nil
}
func (s *PlanService) CreateTrialPlan(userId int) (userPlan connectteam.UserPlan, err error) {
	err = s.repo.DeleteOnConfirmPlan(userId)
	if err != nil {
		return userPlan, err
	}

	return s.repo.CreatePlan(connectteam.UserPlan{
		UserId:         userId,
		PlanType:       "basic",
		HolderId:       userId,
		Status:         connectteam.Active,
		Duration:       14,
		ExpiryDate:     time.Now().Add(14 * 24 * time.Hour),
		PlanAccess:     "holder",
		InvitationCode: "",
		IsTrial:        true,
	})
}

func (s *PlanService) AddUserToAdvanced(holderPlan connectteam.UserPlan, userId int) (userPlan connectteam.UserPlan, err error) {
	// добавить проверку на кол-во участников
	err = s.repo.DeleteOnConfirmPlan(userId)
	if err != nil {
		return userPlan, err
	}
	err = s.repo.SetExpiredStatusWithUserId(userId)
	if err != nil {
		return userPlan, err
	}
	return s.repo.CreatePlan(connectteam.UserPlan{
		UserId:         userId,
		PlanType:       "premium",
		HolderId:       holderPlan.UserId,
		Status:         connectteam.Active,
		Duration:       holderPlan.Duration,
		ExpiryDate:     holderPlan.ExpiryDate,
		PlanAccess:     "additional",
		InvitationCode: holderPlan.InvitationCode,
	})
}

func (s *PlanService) GetMembers(code string) ([]connectteam.UserPublic, error) {
	return s.repo.GetMembers(code)
}

func (s *PlanService) CreatePlan(request connectteam.UserPlan) (userPlan connectteam.UserPlan, err error) {
	err = s.repo.DeleteOnConfirmPlan(request.UserId)
	if err != nil {
		return userPlan, err
	}
	activePlan, _ := s.repo.GetUserActivePlan(request.UserId)
	if activePlan.Id != 0 {
		if err := s.repo.SetExpiredStatus(activePlan.Id); err != nil {
			return userPlan, err
		}
		return s.repo.CreatePlan(request)
	}

	if request.Duration <= 0 {
		return userPlan, errors.New("incorrect value of duration")
	}

	if request.PlanType == "premium" {
		request.InvitationCode, err = generateInviteCode()
		if err != nil {
			return userPlan, err
		}
	}
	return s.repo.CreatePlan(request)
}

func (s *PlanService) DeletePlan(id int) error {
	return s.repo.DeletePlan(id)
}

func (s *PlanService) SetPlanByAdmin(userId int, planType string, expiryDateString string) error {
	activePlan, _ := s.repo.GetUserActivePlan(userId)
	if activePlan.Id != 0 {
		if err := s.repo.SetExpiredStatus(activePlan.Id); err != nil {
			return err
		}
	}
	date, err := time.Parse(time.RFC3339, expiryDateString)
	expiryDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	if err != nil {
		return err
	}

	var invitationCode string
	if !expiryDate.After(time.Now()) {
		return errors.New("incorrect expiry date")
	}
	if planType == "premium" {
		invitationCode, err = generateInviteCode()
		if err != nil {
			return err
		}
	}
	var userPlan = connectteam.UserPlan{
		PlanType:       planType,
		Status:         "active",
		PlanAccess:     "holder",
		ExpiryDate:     expiryDate,
		UserId:         userId,
		HolderId:       userId,
		Duration:       int(time.Until(expiryDate).Hours() / 24),
		InvitationCode: invitationCode,
	}
	_, err = s.repo.CreatePlan(userPlan)

	return err

}

func (s *PlanService) GetUsersPlans() ([]connectteam.UserPlan, error) {
	return s.repo.GetUsersPlans()
}

func (s *PlanService) ConfirmPlan(id int) error {
	return s.repo.SetConfirmed(id)
}

func generateInviteCode() (string, error) {
	// Создание байтового среза для хранения случайных данных
	randomBytes := make([]byte, 16)

	// Заполнение среза случайными данными
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Кодирование случайных данных в base64
	inviteCode := base64.URLEncoding.EncodeToString(randomBytes)

	return inviteCode, nil
}

func (s *PlanService) GetHolderWithInvitationCode(code string) (id int, err error) {
	return s.repo.GetHolderWithInvitationCode(code)
}

func (s *PlanService) DeleteUserFromSub(id int) error {
	return s.repo.DeleteUserFromSub(id)
}
