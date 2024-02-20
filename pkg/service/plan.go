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

func (s *PlanService) GetUserPlan(userId int) (connectteam.UserPlan, error) {
	return s.repo.GetUserPlan(userId)
}

func (s *PlanService) CreatePlan(request connectteam.UserPlan) (connectteam.UserPlan, error) {
	request.ExpiryDate = time.Time{}
	println(request.ExpiryDate.String())
	return s.repo.CreatePlan(request)
}

func (s *PlanService) DeletePlan(id int) error {
	return s.repo.DeletePlan(id)
}

func (s *PlanService) SetPlanByAdmin(userId int, planType string, expiryDateString string) error {
	date, err := time.Parse(time.RFC3339, expiryDateString)
	expiryDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	println(int(time.Until(expiryDate) / 24))

	if err != nil {
		return err
	}

	var userPlan = connectteam.UserPlan{
		PlanType:   planType,
		Confirmed:  true,
		PlanAccess: "holder",
		ExpiryDate: expiryDate,
		UserId:     userId,
		HolderId:   userId,
		Duration: int(time.Until(expiryDate).Hours() / 24),
	}
	_, err = s.repo.CreatePlan(userPlan)

	if err != nil {
		return err
	}

	return err

}

func (s *PlanService) GetUsersPlans() ([]connectteam.UserPlan, error) {
	return s.repo.GetUsersPlans()
}

func (s *PlanService) ConfirmPlan(id int) error {
	return s.repo.SetConfirmed(id)
}
