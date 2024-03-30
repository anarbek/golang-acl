package models

type User struct {
	TenantID int    `json:"tenantId"`
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
	RoleID   int    `json:"roleId"`
}
