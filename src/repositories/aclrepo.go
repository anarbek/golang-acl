package repositories

import (
	"example/hello/models"
	"fmt"
)

var Roles = []models.Role{
	{ID: 1, Name: "User", Code: "User"},
	{ID: 2, Name: "Admin", Code: "Admin"},
}

var Policies = []models.Policy{
	{ID: 1, Name: "SubjectManagement", Code: "SubjectManagement"},
	{ID: 2, Name: "UserManagement", Code: "UserManagement"},
}

var RolePolicies = []models.RolePolicy{
	{RoleID: 1, PolicyID: 1, Read: true, Write: true},
	{RoleID: 1, PolicyID: 2, Read: true, Write: true},
	{RoleID: 2, PolicyID: 1, Read: true, Write: false},
	{RoleID: 2, PolicyID: 2, Read: false, Write: false},
}

var Users = []models.User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", RoleID: 1},
	{ID: 2, Name: "Jane Doe", Email: "jane@example.com", RoleID: 2},
}

type AclAbstract struct {
}

func (acl *AclAbstract) UsersWithRoles() []models.User {
	// Create a map of roles for quick lookup
	roleMap := GetRoleMap()
	// Create a new slice of users with roles
	usersWithRoles := make([]models.User, len(Users))
	for i, user := range Users {
		user.Role = roleMap[user.RoleID]
		usersWithRoles[i] = user
	}

	return usersWithRoles
}

func (acl *AclAbstract) Users() []models.User {
	fmt.Println("users requested!")
	return Users
}

func GetPolicyMap() map[int]models.Policy {
	policyMap := make(map[int]models.Policy)
	for _, policy := range Policies {
		policyMap[policy.ID] = policy
	}
	return policyMap
}

func GetRoleMap() map[int]models.Role {
	// Create a map of policies for quick lookup
	policyMap := GetPolicyMap()
	// Create a map of roles for quick lookup
	roleMap := make(map[int]models.Role)
	for _, role := range Roles {
		// Create a slice of RolePolicy for this role
		var rolePolicies []models.RolePolicy
		for _, rp := range RolePolicies {
			if rp.RoleID == role.ID {
				rp.Policy = policyMap[rp.PolicyID]
				rolePolicies = append(rolePolicies, rp)
			}
		}
		// Assign the RolePolicies slice to the RolePolicies field of the role
		role.RolePolicies = getRolePoliciesInner(role.ID, policyMap)
		// Add the role to the roleMap
		roleMap[role.ID] = role
	}

	return roleMap
}

func getRolePolicies(roleID int) []models.RolePolicy {
	policyMap := GetPolicyMap()
	return getRolePoliciesInner(roleID, policyMap)
}

// This function takes a RoleID and returns a slice of RolePolicies for that role
func getRolePoliciesInner(roleID int, policyMap map[int]models.Policy) []models.RolePolicy {
	var rolePolicies []models.RolePolicy
	for _, rp := range RolePolicies {
		if rp.RoleID == roleID {
			rp.Policy = policyMap[rp.PolicyID]
			rolePolicies = append(rolePolicies, rp)
		}
	}
	return rolePolicies
}

func GetRole(roleID int) models.Role {
	// Get the user's role from the Roles slice
	var role models.Role
	for _, r := range Roles {
		if r.ID == roleID {
			role = r
			break
		}
	}

	role.RolePolicies = getRolePolicies(roleID)
	return role
}
