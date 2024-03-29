package repositories

import (
	"fmt"
	"gokg/gomvc/models"
	"sync"
)

var Roles = []models.Role{
	{ID: models.ConstAdminInt, Name: models.RolesAdmin, Code: models.RolesAdmin},
	{ID: models.ConstUserInt, Name: models.RolesUser, Code: models.RolesUser},
	{ID: models.ConstSuperAdminInt, Name: models.RolesSuperadmin, Code: models.RolesSuperadmin},
	{ID: models.ConstTenantInt, Name: models.RolesTenant, Code: models.RolesTenant},
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

	{RoleID: 3, PolicyID: 1, Read: true, Write: true},
	{RoleID: 3, PolicyID: 2, Read: true, Write: true},

	{RoleID: 4, PolicyID: 1, Read: true, Write: true},
	{RoleID: 4, PolicyID: 2, Read: true, Write: true},
}

var Users = []models.User{
	{ID: 105, Name: "Superadmin", Email: "superadmin@example.com", RoleID: models.ConstSuperAdminInt, TenantID: 105},
	{ID: 106, Name: "User under superadmin", Email: "suser@example.com", RoleID: models.ConstUserInt, TenantID: 105},
	{ID: 107, Name: "Admin under superadmin", Email: "sadmin@example.com", RoleID: models.ConstAdminInt, TenantID: 105},
	{ID: 1, Name: "John Doe", Email: "john@example.com", RoleID: models.ConstAdminInt, TenantID: 1},
	{ID: 2, Name: "Jane Doe", Email: "jane@example.com", RoleID: models.ConstUserInt, TenantID: 1},
	{ID: 3, Name: "Bob Doe", Email: "bob@example.com", RoleID: models.ConstAdminInt, TenantID: 3},
	{ID: 4, Name: "Ken Doe", Email: "ken@example.com", RoleID: models.ConstUserInt, TenantID: 3},
	{ID: 205, Name: "Tenant", Email: "tenant@example.com", RoleID: models.ConstTenantInt, TenantID: 205},
	{ID: 206, Name: "User under tenant", Email: "tuser@example.com", RoleID: models.ConstUserInt, TenantID: 205},
	{ID: 207, Name: "Admin under tenant", Email: "tadmin@example.com", RoleID: models.ConstAdminInt, TenantID: 205},
}

type AclBase struct {
	acl *AclAbstract
}

func (aclBase *AclBase) Init() {
	aclBase.acl = NewAclAbstract()
}
func (aclBase *AclBase) UsersWithRoles() []models.User {
	return aclBase.acl.UsersWithRoles()
}
func (aclBase *AclBase) UsersAll() []models.User {
	return aclBase.acl.UsersAll()
}

func (aclBase *AclBase) Users(loggedInUser *models.User) []models.User {
	return aclBase.acl.Users(loggedInUser)
}

func (aclBase *AclBase) InsertUser(user, loggedInUser *models.User) error {
	// Check the role permissions
	if err := checkRolePermissions(user, loggedInUser); err != nil {
		return err
	}
	//loggedInUserRole := GetRole(loggedInUser.RoleID)
	switch loggedInUser.Role.Name {
	case models.RolesSuperadmin:
		// Superadmin can operate any user with any role
	case models.RolesTenant:
		// Tenant can only operate users with RoleNames: "Admin" and "User"
		user.TenantID = loggedInUser.ID
	case models.RolesAdmin:
		// Admin can only operate users with RoleName: "User"
		user.TenantID = loggedInUser.TenantID
	default:
		return fmt.Errorf("invalid role %v", loggedInUser.Role.Name)
	}

	return aclBase.acl.InsertUser(user, loggedInUser)
}

func (aclBase *AclBase) UpdateUser(user, loggedInUser *models.User) error {
	return aclBase.acl.UpdateUser(user, loggedInUser)
}

func (aclBase *AclBase) DeleteUser(id int, loggedInUser *models.User) error {
	return aclBase.acl.DeleteUser(id, loggedInUser)
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

func (acl *AclAbstract) UsersAll() []models.User {
	fmt.Println("all users requested!")
	return Users
}

func (acl *AclAbstract) Users(loggedInUser *models.User) []models.User {
	fmt.Println("users requested!")
	loggedInUserRole := GetRole(loggedInUser.RoleID)
	switch loggedInUserRole.Name {
	case models.RolesSuperadmin:
		// Superadmin can see all users
		return Users

	case models.RolesTenant, models.RolesAdmin:
		// Tenant and Admin can only see users with the same TenantID
		var filteredUsers []models.User
		for _, user := range Users {
			if user.TenantID == loggedInUser.TenantID {
				filteredUsers = append(filteredUsers, user)
			}
		}
		return filteredUsers

	default:
		// Other roles can only see themselves
		for _, user := range Users {
			if user.ID == loggedInUser.ID {
				return []models.User{user}
			}
		}
	}

	// If no match, return an empty slice
	return []models.User{}
}

func checkRolePermissions(user, loggedInUser *models.User) error {
	userToUpdateRole := GetRole(user.RoleID)
	// Check the role of the loggedInUser and enforce the rules
	switch loggedInUser.Role.Name {
	case models.RolesSuperadmin:
		if userToUpdateRole.Name == models.RolesSuperadmin {
			return LogErr("superadmin cannot operate user with role %v", userToUpdateRole.Name)
		}
		return nil
		// Superadmin can operate any user with any role
	case models.RolesTenant:
		// Tenant can only operate users with RoleNames: "Admin" and "User"
		if userToUpdateRole.Name != models.RolesAdmin && userToUpdateRole.Name != models.RolesUser {
			return LogErr("tenant cannot operate user with role %v", userToUpdateRole.Name)
		}
	case models.RolesAdmin:
		// Admin can only operate users with RoleName: "User"
		if userToUpdateRole.Name != models.RolesUser {
			return LogErr("admin cannot operate user with role %v", userToUpdateRole.Name)
		}
	default:
		return LogErr("invalid role %v", userToUpdateRole.Name)
	}

	return nil
}

func (acl *AclAbstract) InsertUser(user, loggedInUser *models.User) error {
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
	user.TenantID = loggedInUser.TenantID
	acl._userCounter++

	// Append the new user to the Users slice
	Users = append(Users, *user)

	return nil
}

func LogErr(format string, a ...any) error {
	fmt.Println(format, a)
	return fmt.Errorf(format, a...)
}

func (acl *AclAbstract) UpdateUser(user, loggedInUser *models.User) error {
	// Check if the new name is already taken
	for _, existingUser := range Users {
		if existingUser.Name == user.Name && existingUser.ID != user.ID {
			return LogErr("name %v is already taken", user.Name)
		}
	}

	// Lock the mutex before accessing Users
	acl.mu.Lock()
	defer acl.mu.Unlock()

	// Find the user to update
	for i, existingUser := range Users {
		if existingUser.ID == user.ID {
			if err := checkRolePermissions(&existingUser, loggedInUser); err != nil {
				return err
			}
			// Update the user
			Users[i] = *user

			//dont change tenant
			Users[i].TenantID = existingUser.TenantID
			if err := checkRolePermissions(&Users[i], loggedInUser); err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("user with ID %d not found", user.ID)
}

func (acl *AclAbstract) DeleteUser(id int, loggedInUser *models.User) error {
	if loggedInUser.ID == id {
		return LogErr("User has same ID, cannot be deleted")
	}
	allUsers := acl.Users(loggedInUser)
	// Lock the mutex before accessing Users
	acl.mu.Lock()
	defer acl.mu.Unlock()

	// Find the user to delete
	for i, existingUser := range allUsers {
		if existingUser.ID == id {
			// Delete the user
			Users = append(Users[:i], Users[i+1:]...)
			return nil
		}
	}

	return LogErr("user with ID %d not found", id)
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
		//var rolePolicies []models.RolePolicy
		for _, rp := range RolePolicies {
			if rp.RoleID == role.ID {
				rp.Policy = policyMap[rp.PolicyID]
				//rolePolicies = append(rolePolicies, rp)
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
