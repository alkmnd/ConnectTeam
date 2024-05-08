package repository

import (
	connectteam "ConnectTeam/models"
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
	query := fmt.Sprintf("SELECT COALESCE(ah.sub_id, s.id) AS id,"+
		"s.plan_type,"+
		"s.holder_id,"+
		"s.expiry_date,"+
		"s.duration,"+
		"ah.access,"+
		"s.status,"+
		"s.invitation_code,"+
		"ah.user_id,"+
		"s.is_trial FROM %s s JOIN %s ah ON ah.sub_id = s.id WHERE ah.user_id=$1 AND s.status='active'", subscriptionsTable, usersSubsTable)
	err := r.db.Get(&userPlan, query, userId)

	return userPlan, err
}

func (r *PlanPostgres) UpgradePlan(planId uuid.UUID, planType string, invitationCode string) error {
	query := fmt.Sprintf(`UPDATE %s SET plan_type = $1, invitation_code = $2 WHERE id = $3`, subscriptionsTable)

	_, err := r.db.Exec(query, planType, invitationCode, planId)

	return err
}

func (r *PlanPostgres) SetExpiredStatus(id uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET status = 'expired', 
		expiry_date = NOW() WHERE id =$1`, subscriptionsTable)

	_, err := r.db.Exec(query, id)
	return err
}

func (r *PlanPostgres) SetExpiredStatusWithUserId(userId uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET status = 'expired', 
		expiry_date = NOW() WHERE holder_id = $1`, subscriptionsTable)

	_, err := r.db.Exec(query, userId)

	return err
}

func (r *PlanPostgres) AddUserToSubscription(userId uuid.UUID, planId uuid.UUID, access string) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, sub_id, access) values ($1, $2, $3)", usersSubsTable)
	_, err := r.db.Exec(query, userId, planId, access)
	return err
}

func (r *PlanPostgres) DeleteOnConfirmPlan(userId uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE holder_id = $1 AND status='on_confirm'", subscriptionsTable)
	_, err := r.db.Exec(query, userId)
	return err
}

func (r *PlanPostgres) GetUserSubscriptions(userId uuid.UUID) ([]connectteam.UserPlan, error) {
	var usersPlan []connectteam.UserPlan

	query := fmt.Sprintf(`SELECT id, holder_id, plan_type, 
		expiry_date, duration, status FROM %s WHERE holder_id=$1`, subscriptionsTable)
	err := r.db.Select(&usersPlan, query, userId)
	return usersPlan, err
}

func (r *PlanPostgres) GetPlanInvitationCode(code string) (id uuid.UUID, err error) {
	query := fmt.Sprintf(`SELECT holder_id FROM %s WHERE status='active' and invitation_code=$1 and holder_id=user_id LIMIT 1`, subscriptionsTable)
	err = r.db.Select(&id, query, code)
	if err != nil {
		return uuid.Nil, err
	}
	return id, err
}

func (r *PlanPostgres) GetMembers(id uuid.UUID) (users []connectteam.UserPublic, err error) {
	query := fmt.Sprintf(`SELECT u.id, u.email, u.first_name, u.second_name, u.profile_image FROM %s u 
    JOIN %s ah ON ah.user_id = u.id
	JOIN %s s ON s.id = ah.sub_id WHERE s.id=$1 AND status='active'`, usersTable, usersSubsTable, subscriptionsTable)
	err = r.db.Select(&users, query, id)
	return users, err
}

func (r *PlanPostgres) GetHolderWithInvitationCode(code string) (id uuid.UUID, err error) {
	query := fmt.Sprintf(`SELECT holder_id FROM %s WHERE status='active' and invitation_code=$1`, subscriptionsTable)
	err = r.db.Get(&id, query, code)

	return id, err
}

func (r *PlanPostgres) CreatePlan(request connectteam.UserPlan) (connectteam.UserPlan, error) {
	query := fmt.Sprintf(`INSERT INTO %s (holder_id, expiry_date, duration, 
		status, plan_type, invitation_code, is_trial) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING *`, subscriptionsTable)

	var userPlan connectteam.UserPlan
	row := r.db.QueryRow(query, request.HolderId, request.ExpiryDate, request.Duration, request.Status, request.PlanType, request.InvitationCode, request.IsTrial)
	if err := row.Scan(&userPlan.Id, &userPlan.PlanType, &userPlan.HolderId, &userPlan.ExpiryDate, &userPlan.Duration, &userPlan.Status, &userPlan.InvitationCode, &userPlan.IsTrial); err != nil {
		return request, err
	}
	return userPlan, nil
}

func (r *PlanPostgres) DeletePlan(id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", subscriptionsTable)
	_, err := r.db.Exec(query, id)

	return err
}

// change
func (r *PlanPostgres) GetUsersPlans() ([]connectteam.UserPlan, error) {
	var plansUsers []connectteam.UserPlan

	query := fmt.Sprintf(`SELECT id, holder_id, plan_type, expiry_date, duration, status FROM %s WHERE status='active' or status='on_confirm'`, subscriptionsTable)
	err := r.db.Select(&plansUsers, query)
	return plansUsers, err
}

func (r *PlanPostgres) SetConfirmed(id uuid.UUID) error {
	query := fmt.Sprintf(`UPDATE %s SET status = 'active', 
		expiry_date = NOW() + interval '1 day' * duration WHERE id = %d`, subscriptionsTable, id)

	_, err := r.db.Exec(query)

	return err
}

func (r *PlanPostgres) DeleteUserFromSub(userId uuid.UUID, planId uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND sub_id = $2", usersSubsTable)
	_, err := r.db.Exec(query, userId, planId)

	return err
}

func (r *PlanPostgres) GetPlan(planId uuid.UUID) (sub connectteam.Subscription, err error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id=$1`, subscriptionsTable)
	err = r.db.Get(&sub, query, planId)

	return sub, err
}
