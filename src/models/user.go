package models

import "fmt"

type User struct {
	TenantID int    `json:"tenantId"`
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
	RoleID   int    `json:"roleId"`
}

func (item *User) GetID() string {
	IDStr := fmt.Sprintf("%+v", item.ID)
	return IDStr
}

func (item *User) GetObjectStr() string {
	itemStr := fmt.Sprintf("%+v", item)
	return itemStr
}

func (item *User) GetObjectName() string {
	return "User"
}
