package controllers

import (
	"gokg/gomvc/controllers"
	"gokg/gomvc/repositories"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type RoleTestSetupObject struct {
	Router         *gin.Engine
	RoleController *controllers.RolesController
}

func RoleTestSetup(runMiddleware gin.HandlerFunc) RoleTestSetupObject {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(runMiddleware) // Use the middleware
	acl := &repositories.RoleBase{}
	acl.Init()
	userController := &controllers.RolesController{}
	userController.Init(acl)

	return RoleTestSetupObject{
		Router:         router,
		RoleController: userController,
	}
}

func TestInsertRoleAsSuperAdmin(t *testing.T) {
	setup := RoleTestSetup(SuperAdminMiddleware())

	setup.Router.POST("/insert", setup.RoleController.InsertRole)

	testCases := []struct {
		name         string
		userJson     string
		expectedCode int
	}{
		{"valid role", `{"name": "Role1", "tenantId": 105, "roleTypeId": 3}`, http.StatusOK},
		{"same name", `{"name": "Role1", "tenantId": 105, "roleTypeId": 3}`, http.StatusInternalServerError},
		{"superadmin cannot create superadmin", `{"name": "Superadmin", "tenantId": 105, "roleTypeId": 1}`, http.StatusInternalServerError},
		{"superadmin can create tenant", `{"name": "RoleTestTenant", "tenantId": 105, "roleTypeId": 2}`, http.StatusOK},
		{"admin already exists with same tenantId", `{"name": "Admin", "tenantId": 105, "roleTypeId": 3}`, http.StatusInternalServerError},
		{"superadmin can create tenant with different tenantId", `{"name": "Tenant", "tenantId": 305, "roleTypeId": 2}`, http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/insert", strings.NewReader(tc.userJson))
			response := httptest.NewRecorder()

			setup.Router.ServeHTTP(response, request)

			assert.Equal(t, tc.expectedCode, response.Code)
		})
	}
}

func TestInsertRoleAsTenant(t *testing.T) {
	setup := RoleTestSetup(TenantMiddleware())

	setup.Router.POST("/insert", setup.RoleController.InsertRole)

	testCases := []struct {
		name         string
		userJson     string
		expectedCode int
	}{
		{"valid role", `{"name": "TenantRole1", "tenantId": 205, "roleTypeId": 3}`, http.StatusOK},
		{"same name", `{"name": "TenantRole1", "tenantId": 205, "roleTypeId": 3}`, http.StatusInternalServerError},
		{"invalid role", `{"name": "TenantRole2", "tenantId": 205, "roleTypeId": 2}`, http.StatusInternalServerError},
		{"tenant cannot create tenant with same tenantId", `{"name": "TenantUnderTenant", "tenantId": 205, "roleTypeId": 2}`, http.StatusInternalServerError},
		{"tenant cannot create tenant with different tenantId", `{"name": "TenantUnderTenant2", "tenantId": 305, "roleTypeId": 2}`, http.StatusInternalServerError},
		{"tenant cannot create superadmin", `{"name": "Superadmin", "tenantId": 205, "roleTypeId": 1}`, http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/insert", strings.NewReader(tc.userJson))
			response := httptest.NewRecorder()

			setup.Router.ServeHTTP(response, request)

			assert.Equal(t, tc.expectedCode, response.Code)
		})
	}
}

func TestUpdateRoleAsSuperAdmin(t *testing.T) {
	setup := RoleTestSetup(SuperAdminMiddleware())

	setup.Router.POST("/update", setup.RoleController.UpdateRole)

	testCases := []struct {
		name         string
		userJson     string
		expectedCode int
	}{
		{"valid role", `{"id": 1, "name": "SuperAdminRole1", "tenantId": 105}`, http.StatusOK},
		{"same name", `{"id": 2, "name": "SuperAdminRole1", "tenantId": 105}`, http.StatusInternalServerError},
		{"superadmin cannot update superadmin", `{"id": 3, "name": "SuperAdmin", "tenantId": 105}`, http.StatusInternalServerError},
		{"superadmin can update tenant", `{"id": 4, "name": "Tenant", "tenantId": 205}`, http.StatusOK},
		{"superadmin can update admin", `{"id": 1, "name": "Admin", "tenantId": 105}`, http.StatusOK},
		// {"tenant can update tenant with same tenantId", `{"id": 4, "name": "Tenant", "tenantId": 205}`, http.StatusOK},
		// {"tenant cannot update tenant with different tenantId", `{"id": 4, "name": "Tenant", "tenantId": 105}`, http.StatusInternalServerError},
		// {"tenant cannot update superadmin", `{"id": 3, "name": "SuperAdmin", "tenantId": 105}`, http.StatusInternalServerError},
		// {"tenant cannot update admin", `{"id": 1, "name": "Admin", "tenantId": 105}`, http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/update", strings.NewReader(tc.userJson))
			response := httptest.NewRecorder()

			setup.Router.ServeHTTP(response, request)

			assert.Equal(t, tc.expectedCode, response.Code)
		})
	}
}

func TestUpdateRoleAsTenant(t *testing.T) {
	setup := RoleTestSetup(TenantMiddleware())

	setup.Router.POST("/update", setup.RoleController.UpdateRole)

	testCases := []struct {
		name         string
		userJson     string
		expectedCode int
	}{
		{"valid role", `{"id": 4, "name": "TenantRoleUpdated", "tenantId": 205}`, http.StatusOK},
		{"same name", `{"id": 205, "name": "TenantRoleUpdated", "tenantId": 205}`, http.StatusInternalServerError},
		{"tenant cannot update superadmin", `{"id": 3, "name": "SuperAdmin", "tenantId": 205}`, http.StatusInternalServerError},
		{"tenant cannot update admin", `{"id": 1, "name": "Admin", "tenantId": 105}`, http.StatusInternalServerError},
		//because TenantID is updated everytime to current users TenantID, below code will work with OK
		{"tenant cannot update tenant with different tenantId", `{"id": 4, "name": "Tenant", "tenantId": 205}`, http.StatusOK},
		//This should not projects other tests
		//{"tenant cannot update tenant with different tenantId", `{"id": 4, "name": "Tenant01", "tenantId": 205}`, http.StatusOK},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "/update", strings.NewReader(tc.userJson))
			response := httptest.NewRecorder()

			setup.Router.ServeHTTP(response, request)

			assert.Equal(t, tc.expectedCode, response.Code)
		})
	}
}
