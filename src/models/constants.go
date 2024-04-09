package models

const Acl_read string = "read"
const Acl_write string = "write"

const RolesSuperadmin string = "Superadmin"
const RolesTenant string = "Tenant"
const RolesAdmin string = "Admin"
const RolesUser string = "User"

const (
	ConstAdminInt      = 1
	ConstUserInt       = 2
	ConstSuperAdminInt = 3
	ConstTenantInt     = 4
)

const (
	ConstRoleTypeSuperAdminInt = 1
	ConstRoleTypeTenantInt     = 2
	ConstRoleTypeOtherInt      = 3
)
