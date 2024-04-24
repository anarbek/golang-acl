package repositories

import "gokg/gomvc/models"

type IRoleRepo interface {
	RolesAll() []models.Role
	RolesByTenantID(TenantID int) []models.Role
	RolesByRoleID(RoleID int) []models.Role
	InsertRole(roleToInsert *models.Role, loggedInUser *models.User) error
	UpdateRole(roleToUpdate *models.Role, loggedInUser *models.User) error
	DeleteRole(id int, loggedInUser *models.User, fnGetRoles RolesGetFunc) error
}
