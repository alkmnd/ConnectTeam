package repository

import (
	connectteam "ConnectTeam"
	"fmt"
	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
)

type PlanPostgres struct {
	db *sqlx.DB
}

func NewPlanPostgres(db *sqlx.DB) *PlanPostgres {
	return &PlanPostgres{db: db}
}

func (r *PlanPostgres) GetUserActivePlan(userId uuid.UUID) (connectteam.UserPlan, error) {
	var userPlan connectteam.UserPlan
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 AND (status='active' or status='on_confirm')", plansUsersTable)
	err := r.db.Get(&userPlan, query, userId)

	return userPlan, err
}
func (r *PlanPostgres) SetExpiredStatus(id uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET status = 'expired', 
		expiry_date = NOW() WHERE id = %d`, plansUsersTable, id)

	_, err := r.db.Exec(query)
	return err
}

func (r *PlanPostgres) SetExpiredStatusWithUserId(userId uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET status = 'expired', 
		expiry_date = NOW() WHERE user_id = %d`, plansUsersTable, userId)

	_, err := r.db.Exec(query)

	return err
}

func (r *PlanPostgres) DeleteOnConfirmPlan(userId uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND status='on_confirm'", plansUsersTable)
	_, err := r.db.Exec(query, userId)
	return err
}

func (r *PlanPostgres) GetUserSubscriptions(userId uuid.UUID) ([]connectteam.UserPlan, error) {
	var usersPlan []connectteam.UserPlan

	query := fmt.Sprintf(`SELECT id, user_id, holder_id, plan_type, 
		plan_access, expiry_date, duration, status FROM %s WHERE user_id=$1`, plansUsersTable)
	err := r.db.Select(&usersPlan, query, userId)
	return usersPlan, err
}

func (r *PlanPostgres) GetPlanInvitationCode(code string) (id uuid.UUID, err error) {
	query := fmt.Sprintf(`SELECT holder_id FROM %s WHERE status='active' and invitation_code=$1 and holder_id=user_id LIMIT 1`, plansUsersTable)
	err = r.db.Select(&id, query, code)
	if err != nil {
		return uuid.Nil, err
	}
	return id, err
}

//func (r *PlanPostgres) UpdateData() error {
//
//

func (r *PlanPostgres) GetMembers(code string) (users []connectteam.UserPublic, err error) {
	query := fmt.Sprintf(`SELECT u.id, u.email, u.first_name, u.second_name, u.profile_image FROM %s u
	JOIN %s p ON p.user_id = u.id WHERE invitation_code=$1 AND status='active'`, usersTable, plansUsersTable)
	err = r.db.Select(&users, query, code)
	return users, err
}

func (r *PlanPostgres) GetHolderWithInvitationCode(code string) (id uuid.UUID, err error) {
	query := fmt.Sprintf(`SELECT holder_id FROM %s WHERE status='active' and invitation_code=$1 and holder_id=user_id LIMIT 1`, plansUsersTable)
	err = r.db.Get(&id, query, code)

	return id, err
}

func (r *PlanPostgres) CreatePlan(request connectteam.UserPlan) (connectteam.UserPlan, error) {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, holder_id, expiry_date, duration, plan_access, 
		status, plan_type, invitation_code, is_trial) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING *`, plansUsersTable)

	var userPlan connectteam.UserPlan
	row := r.db.QueryRow(query, request.UserId, request.HolderId, request.ExpiryDate, request.Duration, request.PlanAccess, request.Status, request.PlanType, request.InvitationCode, request.IsTrial)
	if err := row.Scan(&userPlan.Id, &userPlan.PlanType, &userPlan.UserId, &userPlan.HolderId, &userPlan.ExpiryDate, &userPlan.Duration, &userPlan.PlanAccess, &userPlan.Status, &userPlan.InvitationCode, &userPlan.IsTrial); err != nil {
		return request, err
	}
	return userPlan, nil
}

func (r *PlanPostgres) DeletePlan(id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", plansUsersTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r *PlanPostgres) GetUsersPlans() ([]connectteam.UserPlan, error) {
	var plansUsers []connectteam.UserPlan

	query := fmt.Sprintf(`SELECT id, user_id, holder_id, plan_type, 
		plan_access, expiry_date, duration, status FROM %s WHERE status='active' or status='on_confirm'`, plansUsersTable)
	err := r.db.Select(&plansUsers, query)
	return plansUsers, err
}

func (r *PlanPostgres) SetConfirmed(id uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET status = 'active', 
		expiry_date = NOW() + interval '1 day' * duration WHERE id = %d`, plansUsersTable, id)

	_, err := r.db.Exec(query)

	return err
}

func (r *PlanPostgres) DeleteUserFromSub(id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", plansUsersTable)
	_, err := r.db.Exec(query, id)

	return err
}
