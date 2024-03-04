package service

import (
	"ConnectTeam"
	"ConnectTeam/pkg/repository"
	"time"
)

type PlanService struct {
	repo repository.Plan
}

func NewPlanService(repo repository.Plan) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) GetUserActivePlan(userId int) (connectteam.UserPlan, error) {
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
		UserId:     userId,
		PlanType:   "trial",
		HolderId:   userId,
		Status:     connectteam.Active,
		Duration:   14,
		ExpiryDate: time.Now().Add(14 * 24 * time.Hour),
		PlanAccess: "holder",
	})
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

	request.ExpiryDate = time.Time{}
	request.Status = connectteam.OnConfirm
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

	var userPlan = connectteam.UserPlan{
		PlanType:   planType,
		Status:     "active",
		PlanAccess: "holder",
		ExpiryDate: expiryDate,
		UserId:     userId,
		HolderId:   userId,
		Duration:   int(time.Until(expiryDate).Hours() / 24),
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

// get trial
