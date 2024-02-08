package repository

import (
	connectteam "ConnectTeam"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type PlanPostgres struct {
	db *sqlx.DB
}

func NewPlanPostgres(db *sqlx.DB) *PlanPostgres {
	return &PlanPostgres{db: db}
}

func (r *PlanPostgres) GetUserPlan(userId int) (connectteam.UserPlan, error) {
	var userPlan connectteam.UserPlan
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", plansUsersTable)
	err := r.db.Get(&userPlan, query, userId)

	return userPlan, err
}

func (r *PlanPostgres) CreatePlan(request connectteam.UserPlan) (connectteam.UserPlan, error) {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, holder_id, expiry_date, duration, plan_access, confirmed, plan_type) VALUES ($1, $2, $3, $4, $5, $6, $7) 
	ON CONFLICT (user_id) DO UPDATE SET holder_id = $2, expiry_date = $3, duration = $4, plan_access = $5, confirmed = $6, plan_type = $7
	RETURNING *`, plansUsersTable)

	var userPlan connectteam.UserPlan
	row := r.db.QueryRow(query, request.UserId, request.UserId, time.Time{}, request.Duration, "holder", false, request.PlanType)
	if err := row.Scan(&userPlan.PlanType, &userPlan.UserId, &userPlan.HolderId, &userPlan.ExpiryDate, &userPlan.Duration, &userPlan.PlanAccess, &userPlan.Confirmed); err != nil {
		return request, err
	}

	println(userPlan.HolderId)


	return userPlan, nil
}

func (r *PlanPostgres) GetUsersPlans() ([] connectteam.UserPlan, error) {
	var plansUsers []connectteam.UserPlan

	query := fmt.Sprintf("SELECT user_id, holder_id, plan_type, plan_access, expiry_date, duration, confirmed FROM %s", plansUsersTable)
	err := r.db.Select(&plansUsers, query)
	return plansUsers, err
}