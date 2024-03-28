package repositories

import (
	"fmt"
	"gokg/gomvc/models"
	"sync"
)

var Roles = []models.Role{
	{ID: 1, Name: "Admin", Code: "Admin"},
	{ID: 2, Name: "User", Code: "User"},
}

var Policies = []models.Policy{
	{ID: 1, Name: "UserManagement", Code: "UserManagement"},
	{ID: 2, Name: "SubjectManagement", Code: "SubjectManagement"},
}

var RolePolicies = []models.RolePolicy{
	{RoleID: 1, PolicyID: 1, Read: true, Write: true},
	{RoleID: 1, PolicyID: 2, Read: true, Write: true},
	{RoleID: 2, PolicyID: 1, Read: true, Write: false},
	{RoleID: 2, PolicyID: 2, Read: true, Write: true},
}

var Users = []models.User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", RoleID: 1},
	{ID: 2, Name: "Jane Doe", Email: "jane@example.com", RoleID: 2},
}

type AclAbstract struct {
	_userCounter int
	mu           sync.Mutex
}

func NewAclAbstract() *AclAbstract {
	aclAbstract := &AclAbstract{}
	aclAbstract._userCounter = len(Users)
	return aclAbstract
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

func (acl *AclAbstract) InsertUser(user *models.User) error {
	// Check if the user already exists
	for _, existingUser := range Users {
		if existingUser.Name == user.Name {
			return fmt.Errorf("user with Name %v already exists", user.Name)
		}
	}

	// Lock the mutex before accessing _userCounter
	acl.mu.Lock()
	defer acl.mu.Unlock() // Move the defer statement here

	user.ID = acl._userCounter + 1
	acl._userCounter++

	// Append the new user to the Users slice
	Users = append(Users, *user)

	return nil
}

func (acl *AclAbstract) UpdateUser(user *models.User) error {
	// Check if the new name is already taken
	for _, existingUser := range Users {
		if existingUser.Name == user.Name && existingUser.ID != user.ID {
			return fmt.Errorf("name %v is already taken", user.Name)
		}
	}

	// Lock the mutex before accessing Users
	acl.mu.Lock()
	defer acl.mu.Unlock()

	// Find the user to update
	for i, existingUser := range Users {
		if existingUser.ID == user.ID {
			// Update the user
			Users[i] = *user
			return nil
		}
	}

	return fmt.Errorf("user with ID %d not found", user.ID)
}

func (acl *AclAbstract) DeleteUser(id int) error {
	// Lock the mutex before accessing Users
	acl.mu.Lock()
	defer acl.mu.Unlock()

	// Find the user to delete
	for i, existingUser := range Users {
		if existingUser.ID == id {
			// Delete the user
			Users = append(Users[:i], Users[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("user with ID %d not found", id)
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
