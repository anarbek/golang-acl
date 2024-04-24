package models

import "fmt"

type Role struct {
	ID           int          `json:"id"`
	TenantID     int          `json:"tenantId"`
	RoleTypeId   int          `json:"roleTypeId"`
	Code         string       `json:"code"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	RolePolicies []RolePolicy `json:"rolePolicies"`
}

func (item *Role) GetID() string {
	IDStr := fmt.Sprintf("%+v", item.ID)
	return IDStr
}

func (item *Role) GetObjectStr() string {
	itemStr := fmt.Sprintf("%+v", item)
	return itemStr
}

func (item *Role) GetObjectName() string {
	return "Role"
}

type Policy struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RolePolicy struct {
	RoleID     int    `json:"roleId"`
	Role       Role   `json:"role"`
	Policy     Policy `json:"policy"`
	PolicyName string `json:"policyName"`
	PolicyID   int    `json:"policyId"`
	Read       bool   `json:"read"`
	Write      bool   `json:"write"`
}
