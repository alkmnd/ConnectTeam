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

func (s *PlanService) GetUserPlan(user_id int) (connectteam.UserPlan, error) {
	return s.repo.GetUserPlan(user_id)
}

func (s *PlanService) CreatePlan(request connectteam.UserPlan) (connectteam.UserPlan, error) {
	request.ExpiryDate = time.Time{}
	return s.repo.CreatePlan(request)
}

func (s *PlanService) SetPlanByAdmin(user_id int, duration int, plan_type string) (error) {
	expiryDate := time.Now().Add(time.Duration(duration) * 24 * time.Hour)
	var userPlan = connectteam.UserPlan{
		PlanType: plan_type,
		Confirmed: true, 
		PlanAccess: "holder", 
		ExpiryDate: expiryDate, 
		UserId: user_id, 
		HolderId: user_id,
	}
	_, err := s.repo.CreatePlan(userPlan)

	if err != nil {
		return err 
	}

	return err

}

func (s *PlanService) GetUsersPlans() ([] connectteam.UserPlan, error) {
	return s.repo.GetUsersPlans()
}

func (s *PlanService) ConfirmPlan(id int) (error) {
	return s.repo.SetConfirmed(id)
}
