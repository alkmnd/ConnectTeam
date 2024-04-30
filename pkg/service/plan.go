package service

import (
	"ConnectTeam"
	"ConnectTeam/pkg/repository"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"time"
)

type PlanService struct {
	planRepo  repository.Plan
	yooClient repository.Payment
}

func NewPlanService(repo repository.Plan, client repository.Payment) *PlanService {
	return &PlanService{planRepo: repo, yooClient: client}
}

func (s *PlanService) GetUserActivePlan(userId uuid.UUID) (connectteam.UserPlan, error) {
	return s.planRepo.GetUserActivePlan(userId)
}

func (s *PlanService) CheckIfSubscriptionExists(userId uuid.UUID) (bool, error) {
	var userPlans []connectteam.UserPlan

	userPlans, err := s.planRepo.GetUserSubscriptions(userId)
	if err != nil {
		return false, err
	}

	return len(userPlans) > 0, nil
}

func (s *PlanService) GetUserSubscriptions(userId uuid.UUID) ([]connectteam.UserPlan, error) {
	var userPlans []connectteam.UserPlan

	userPlans, err := s.planRepo.GetUserSubscriptions(userId)
	if err != nil {
		return userPlans, err
	}

	return userPlans, nil
}

func (s *PlanService) CreateTrialPlan(userId uuid.UUID) (userPlan connectteam.UserPlan, err error) {
	err = s.planRepo.DeleteOnConfirmPlan(userId)
	if err != nil {
		return userPlan, err
	}

	return s.planRepo.CreatePlan(connectteam.UserPlan{
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

func (s *PlanService) AddUserToAdvanced(holderPlan connectteam.UserPlan, userId uuid.UUID) (userPlan connectteam.UserPlan, err error) {

	err = s.planRepo.SetExpiredStatusWithUserId(userId)
	if err != nil {
		return userPlan, err
	}
	err = s.planRepo.AddUserToSubscription(userId, holderPlan.Id, "additional")

	if err != nil {
		return userPlan, err
	}

	return s.planRepo.GetUserActivePlan(userId)
}

func (s *PlanService) GetMembers(id uuid.UUID) ([]connectteam.UserPublic, error) {
	return s.planRepo.GetMembers(id)
}

func (s *PlanService) UpgradePlan(orderId string, planId uuid.UUID, userId uuid.UUID) error {
	activePlan, _ := s.planRepo.GetUserActivePlan(userId)
	if activePlan.HolderId != userId {
		return errors.New("permission denied")
	}
	if activePlan.Id != planId {
		return errors.New("incorrect plan id")
	}
	payment, err := s.yooClient.GetPayment(orderId)
	if err != nil {
		return err
	}
	if !payment.Paid {
		return errors.New("order id not paid")
	}
	if payment.Metadata.UserId != userId.String() {
		return errors.New("permission denied")
	}
	if payment.Metadata.PlanType != "basic" &&
		payment.Metadata.PlanType != "advanced" &&
		payment.Metadata.PlanType != "premium" {
		return errors.New("incorrect plan type")
	}

	if payment.Metadata.PlanType == "premium" {
		activePlan.InvitationCode, err = generateInviteCode()
		if err != nil {
			return err
		}
	}

	return s.planRepo.UpgradePlan(activePlan.Id, payment.Metadata.PlanType, activePlan.InvitationCode)

}

func (s *PlanService) CreatePlan(orderId string, userId uuid.UUID) (userPlan connectteam.UserPlan, err error) {
	activePlan, _ := s.planRepo.GetUserActivePlan(userId)
	if activePlan.Id != uuid.Nil {
		return userPlan, errors.New("user already has subscription")
	}

	payment, err := s.yooClient.GetPayment(orderId)
	if !payment.Paid {
		return userPlan, errors.New("order id not paid")
	}
	if payment.Metadata.UserId != userId.String() {
		return userPlan, errors.New("permission denied")
	}
	if err != nil {
		return userPlan, err
	}
	if payment.Metadata.PlanType != "basic" &&
		payment.Metadata.PlanType != "advanced" &&
		payment.Metadata.PlanType != "premium" {
		return userPlan, errors.New("incorrect plan type")
	}

	if payment.Metadata.PlanType == "premium" {
		userPlan.InvitationCode, err = generateInviteCode()
		if err != nil {
			return userPlan, err
		}
	}

	userPlan, err = s.planRepo.CreatePlan(connectteam.UserPlan{
		PlanType:       payment.Metadata.PlanType,
		HolderId:       userId,
		ExpiryDate:     time.Now().Add(time.Hour * 24 * 30),
		Duration:       30,
		Status:         connectteam.Active,
		InvitationCode: userPlan.InvitationCode,
		IsTrial:        false,
	})

	if err != nil {
		return userPlan, err
	}

	err = s.planRepo.AddUserToSubscription(userId, userPlan.Id, "holder")

	return userPlan, err
}

func (s *PlanService) DeletePlan(id uuid.UUID) error {
	return s.planRepo.DeletePlan(id)
}

func (s *PlanService) SetPlanByAdmin(userId uuid.UUID, planType string, expiryDateString string) error {
	activePlan, _ := s.planRepo.GetUserActivePlan(userId)
	if activePlan.Id != uuid.Nil {
		if err := s.planRepo.SetExpiredStatus(activePlan.Id); err != nil {
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
		HolderId:       userId,
		Duration:       int(time.Until(expiryDate).Hours() / 24),
		InvitationCode: invitationCode,
	}
	_, err = s.planRepo.CreatePlan(userPlan)

	return err

}

func (s *PlanService) GetUsersPlans() ([]connectteam.UserPlan, error) {
	return s.planRepo.GetUsersPlans()
}

func (s *PlanService) ConfirmPlan(id uuid.UUID) error {
	return s.planRepo.SetConfirmed(id)
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

func (s *PlanService) GetHolderWithInvitationCode(code string) (id uuid.UUID, err error) {
	return s.planRepo.GetHolderWithInvitationCode(code)
}

func (s *PlanService) DeleteUserFromSub(userId uuid.UUID, planId uuid.UUID) error {
	return s.planRepo.DeleteUserFromSub(userId, planId)
}
