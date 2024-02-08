package service

import (
	"ConnectTeam"
	"ConnectTeam/pkg/repository"
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
	return s.repo.CreatePlan(request)
}

func (s *PlanService) GetUsersPlans() ([] connectteam.UserPlan, error) {
	return s.repo.GetUsersPlans()
}

