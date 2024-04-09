package models

type Role struct {
	ID           int          `json:"id"`
	TenantID     int          `json:"tenantId"`
	RoleTypeId   int          `json:"roleTypeId"`
	Code         string       `json:"code"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	RolePolicies []RolePolicy `json:"rolePolicies"`
}

type Policy struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RolePolicy struct {
	RoleID   int    `json:"roleId"`
	Role     Role   `json:"role"`
	Policy   Policy `json:"policy"`
	PolicyID int    `json:"policyId"`
	Read     bool   `json:"read"`
	Write    bool   `json:"write"`
}
