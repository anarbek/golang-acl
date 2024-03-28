package models

type User struct {
	TenantID int    `json:"tenantId"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
	RoleID   int    `json:"roleId"`
}
