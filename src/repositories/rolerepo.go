package repositories

import (
	"fmt"
	"gokg/gomvc/models"
	"sync"
)

type RolesGetFunc func(loggedInUser *models.User) []models.Role

type RoleBase struct {
	roleAbstract IRoleRepo //*RoleAbstract
}

func (roleBase *RoleBase) Init(iRoleRepo IRoleRepo) {
	roleBase.roleAbstract = iRoleRepo //NewRoleAbstract(auditor)
}

func (roleBase *RoleBase) RolesAll() []models.Role {
	return roleBase.roleAbstract.RolesAll()
}

func (roleBase *RoleBase) GetPermissionsForLoggedinUser(loggedInUser *models.User) []string {
	currRole := GetRole(loggedInUser.RoleID)
	var permissions []string

	// Iterate over the RolePolicies of the current role
	for _, rolePolicy := range currRole.RolePolicies {
		// Get the policy associated with the RolePolicy
		policy := GetPolicy(rolePolicy.PolicyID)

		// Check the Read and Write permissions and add them to the permissions slice
		if rolePolicy.Read {
			permissions = append(permissions, policy.Code+".r")
		}
		if rolePolicy.Write {
			permissions = append(permissions, policy.Code+".w")
		}
	}

	// Return the permissions slice
	return permissions

}

func (roleBase *RoleBase) Roles(loggedInUser *models.User) []models.Role {
	fmt.Println("users requested!")
	loggedInUserRole := GetRole(loggedInUser.RoleID)
	switch loggedInUserRole.Name {
	case models.RolesSuperadmin:
		// Superadmin can see all users
		return roleBase.roleAbstract.RolesAll()

	case models.RolesTenant, models.RolesAdmin:
		// Tenant and Admin can only see users with the same TenantID
		return roleBase.roleAbstract.RolesByTenantID(loggedInUser.TenantID)

	default:
		// Other roles can only see themselves
		return roleBase.roleAbstract.RolesByRoleID(loggedInUser.ID)
	}
}

func (roleBase *RoleBase) GetRole(loggedInUser *models.User, id int) (models.Role, error) {
	roles := roleBase.Roles(loggedInUser)
	if len(roles) > 0 {
		for _, existingRole := range roles {
			if existingRole.ID == id {
				existingRole.RolePolicies = getRolePolicies(id)
				return existingRole, nil
			}
		}
	}
	return models.Role{}, LogErr("role not found: %v", id)
}

func (roleBase *RoleBase) InsertRole(roleToInsert *models.Role, loggedInUser *models.User) error {
	// Check if the user already exists
	for _, existingRole := range Roles {
		if existingRole.Name == roleToInsert.Name && existingRole.TenantID == loggedInUser.TenantID {
			return LogErr("role with Name %v already exists for current tenant", roleToInsert.Name)
		}
	}
	// Check the role permissions
	if err := checkCurrentRolePermissions(roleToInsert, loggedInUser); err != nil {
		return err
	}
	return roleBase.roleAbstract.InsertRole(roleToInsert, loggedInUser)
}

func (roleBase *RoleBase) UpdateRole(roleToUpdate *models.Role, loggedInUser *models.User) error {
	// Check if the new name is already taken
	for _, existingRole := range Roles {
		if existingRole.Name == roleToUpdate.Name && existingRole.ID != roleToUpdate.ID && existingRole.TenantID == roleToUpdate.TenantID {
			return LogErr("name %v is already taken for current tenant", roleToUpdate.Name)
		}
	}
	return roleBase.roleAbstract.UpdateRole(roleToUpdate, loggedInUser)
}

func (roleBase *RoleBase) DeleteRole(id int, loggedInUser *models.User) error {
	if loggedInUser.ID == id {
		return LogErr("Role has same ID, cannot be deleted")
	}
	fnGetRoles := roleBase.Roles
	return roleBase.roleAbstract.DeleteRole(id, loggedInUser, fnGetRoles)
}

type RoleAbstract struct {
	_roleCount int
	mu         sync.Mutex
	auditor    *AuditBase
}

func NewRoleAbstract(auditor *AuditBase) *RoleAbstract {
	roleAbstract := &RoleAbstract{}
	roleAbstract._roleCount = len(Roles)
	roleAbstract.auditor = auditor

	audits := auditor.auditInterface.AuditsAll()
	LogErr("audits len roleRepo: %v", len(audits))
	return roleAbstract
}

func (acl *RoleAbstract) RolesAll() []models.Role {
	fmt.Println("all users requested!")
	return Roles
}

func (acl *RoleAbstract) RolesByTenantID(TenantID int) []models.Role {
	fmt.Println("RolesByTenantID requested!")
	var filteredRoles []models.Role
	for _, user := range Roles {
		if user.TenantID == TenantID {
			filteredRoles = append(filteredRoles, user)
		}
	}
	return filteredRoles
}

func (acl *RoleAbstract) RolesByRoleID(RoleID int) []models.Role {
	for _, role := range Roles {
		if role.ID == RoleID {
			return []models.Role{role}
		}
	}
	return []models.Role{}
}

// Only SuperAdmin and Tenant can work with roles
func checkCurrentRolePermissions(role *models.Role, loggedInUser *models.User) error {
	roleToUpdateRole := role //GetRole(role.ID)
	if role.ID > 0 {
		tempRole := GetRole(role.ID)
		roleToUpdateRole = &tempRole
	}
	loggedInRole := GetRole(loggedInUser.RoleID)
	// Check the role of the loggedInUser and enforce the rules
	switch loggedInRole.RoleTypeId {
	case models.ConstRoleTypeSuperAdminInt:
		// if roleToUpdateRole.RoleTypeId == models.ConstRoleTypeSuperAdminInt {
		// 	return LogErr("superadmin cannot operate user with role %v", roleToUpdateRole.Name)
		// }
		return nil
	case models.ConstRoleTypeTenantInt:
		switch roleToUpdateRole.RoleTypeId {
		case models.ConstRoleTypeSuperAdminInt:
			return LogErr("Tenant cannot work with superadmin")
		case models.ConstRoleTypeTenantInt:
			if role.ID != loggedInUser.RoleID {
				return LogErr("Tenant can work only with its own user roles")
			}
		default:
			if role.TenantID != loggedInUser.TenantID {
				return LogErr("tenant cannot operate user with role %v of other tenant", roleToUpdateRole.Name)
			}
		}
	default:
		return LogErr("invalid role %v", roleToUpdateRole.Name)
	}

	return nil
}

func (acl *RoleAbstract) InsertRole(roleToInsert *models.Role, loggedInUser *models.User) error {

	// Lock the mutex before accessing _userCounter
	acl.mu.Lock()
	defer acl.mu.Unlock() // Move the defer statement here

	roleToInsert.ID = acl._roleCount + 1
	roleToInsert.TenantID = loggedInUser.TenantID
	acl._roleCount++

	// Append the new user to the Users slice
	Roles = append(Roles, *roleToInsert)

	// Insert the role policies
	for _, newPolicy := range roleToInsert.RolePolicies {
		// Create a copy of newPolicy
		newPolicyCopy := models.RolePolicy{
			RoleID:   roleToInsert.ID,
			PolicyID: newPolicy.PolicyID,
			Read:     newPolicy.Read,
			Write:    newPolicy.Write,
		}

		// Add new policy
		RolePolicies = append(RolePolicies, newPolicyCopy)
	}

	acl.auditor.auditInterface.CreateInsertEvent(loggedInUser.ID, roleToInsert)
	return nil
}

func (acl *RoleAbstract) UpdateRole(roleToUpdate *models.Role, loggedInUser *models.User) error {

	// Lock the mutex before accessing Users
	acl.mu.Lock()
	defer acl.mu.Unlock()
	var roleOld *models.Role
	// Find the user to update
	for i, existingRole := range Roles {
		if existingRole.ID == roleToUpdate.ID {
			roleOld = &existingRole
			if err := checkCurrentRolePermissions(&existingRole, loggedInUser); err != nil {
				return err
			}
			// Update the role
			Roles[i] = *roleToUpdate
			//dont change tenant
			Roles[i].TenantID = existingRole.TenantID
			Roles[i].RoleTypeId = existingRole.RoleTypeId

			// Update the role policies
			for _, newPolicy := range roleToUpdate.RolePolicies {
				currPolicyFound := false
				for j, existingPolicy := range RolePolicies {
					if existingPolicy.RoleID == newPolicy.RoleID && existingPolicy.PolicyID == newPolicy.PolicyID {
						// Update the policy
						RolePolicies[j].Read = newPolicy.Read
						RolePolicies[j].Write = newPolicy.Write
						currPolicyFound = true
						break
					}
				}
				if !currPolicyFound {
					// Create a copy of newPolicy
					newPolicyCopy := models.RolePolicy{
						RoleID:   newPolicy.RoleID,
						PolicyID: newPolicy.PolicyID,
						Read:     newPolicy.Read,
						Write:    newPolicy.Write,
					}

					// Add new policy
					RolePolicies = append(RolePolicies, newPolicyCopy)
				}
			}

			if err := checkCurrentRolePermissions(&Roles[i], loggedInUser); err != nil {
				return err
			}

			acl.auditor.auditInterface.CreateUpdateEvent(loggedInUser.ID, roleOld, roleToUpdate)
			return nil
		}
	}

	return LogErr("role with ID %d not found", roleToUpdate.ID)
}

func (acl *RoleAbstract) DeleteRole(id int, loggedInUser *models.User, fnGetRoles RolesGetFunc) error {

	allRoles := fnGetRoles(loggedInUser)
	// Lock the mutex before accessing Users
	acl.mu.Lock()
	defer acl.mu.Unlock()

	// Find the user to delete
	for i, existingRole := range allRoles {
		if existingRole.ID == id {
			// Delete the user
			Users = append(Users[:i], Users[i+1:]...)
			acl.auditor.auditInterface.CreateDeleteEvent(loggedInUser.ID, &existingRole)
			return nil
		}
	}

	return LogErr("user with ID %d not found", id)
}

func GetPolicy(policyID int) models.Policy {
	// Get the policy from the Policies slice
	var policy models.Policy
	for _, p := range Policies {
		if p.ID == policyID {
			policy = p
			break
		}
	}
	return policy
}
