package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gokg/gomvc/controllers"
	"gokg/gomvc/models"
	"gokg/gomvc/repositories"
)

type TestSetup struct {
	Router         *gin.Engine
	UserController *controllers.UserController
}

func Setup(runMiddleware gin.HandlerFunc) TestSetup {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(runMiddleware) // Use the middleware
	acl := &repositories.AclBase{}
	acl.Init()
	userController := &controllers.UserController{}
	userController.Init(acl)

	return TestSetup{
		Router:         router,
		UserController: userController,
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Here you would retrieve your user. This is just an example.
		user := models.User{
			ID:       1,
			Name:     "Test User",
			TenantID: 1,
			RoleID:   models.ConstAdminInt,
			Role: models.Role{
				ID:   models.ConstAdminInt,
				Code: models.RolesAdmin,
				Name: models.RolesAdmin,
			},
			// Fill out the rest of the user fields...
		}

		// Store the user in the context
		c.Set("user", &user)

		// Continue with the next handler in the chain
		c.Next()
	}
}

func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Here you would retrieve your user. This is just an example.
		user := models.User{
			ID:       105,
			Name:     "Test SuperAdmin",
			TenantID: 105,
			RoleID:   models.ConstSuperAdminInt,
			Role: models.Role{
				ID:   models.ConstSuperAdminInt,
				Code: models.RolesSuperadmin,
				Name: models.RolesSuperadmin,
			},
			// Fill out the rest of the user fields...
		}

		// Store the user in the context
		c.Set("user", &user)

		// Continue with the next handler in the chain
		c.Next()
	}
}

func TestInsertUserAsAdmin(t *testing.T) {
	setup := Setup(AdminMiddleware())

	setup.Router.POST("/insert", setup.UserController.InsertUser)

	testCases := []struct {
		name         string
		userJson     string
		expectedCode int
	}{
		{"valid user", `{"name": "Alice", "email": "alice@example.com", "roleId": 2}`, http.StatusOK},
		{"same name", `{"name": "Alice", "email": "alice@example.com", "roleId": 2}`, http.StatusInternalServerError},
		{"admin cannot create admin role", `{"name": "Alice2", "email": "alice@example.com", "roleId": 1}`, http.StatusInternalServerError},
		{"admin cannot create superadmin", `{"name": "Alice3", "email": "alice@example.com", "roleId": 3}`, http.StatusInternalServerError},
		{"admin cannot create tenant", `{"name": "Alice4", "email": "alice@example.com", "roleId": 4}`, http.StatusInternalServerError},
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

func TestInsertUserAsSuperAdmin(t *testing.T) {
	setup := Setup(SuperAdminMiddleware())

	setup.Router.POST("/insert", setup.UserController.InsertUser)

	testCases := []struct {
		name         string
		userJson     string
		expectedCode int
	}{
		{"valid user", `{"name": "AliceS", "email": "alice@example.com", "roleId": 2}`, http.StatusOK},
		{"same name", `{"name": "AliceS", "email": "alice@example.com", "roleId": 2}`, http.StatusInternalServerError},
		{"superadmin can create admin role", `{"name": "Alice2", "email": "alice@example.com", "roleId": 1}`, http.StatusOK},
		{"same name error should be", `{"name": "Alice2", "email": "alice@example.com", "roleId": 4}`, http.StatusInternalServerError},
		{"superadmin can create tenant role", `{"name": "Alice12", "email": "alice@example.com", "roleId": 4}`, http.StatusOK},
		{"superadmin cannot create superadmin", `{"name": "Alice3", "email": "alice@example.com", "roleId": 3}`, http.StatusInternalServerError},
		{"superadmin can create tenant", `{"name": "Alice4", "email": "alice@example.com", "roleId": 4}`, http.StatusOK},
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

func TestUpdateUserAsAdmin(t *testing.T) {
	setup := Setup(AdminMiddleware())

	setup.Router.POST("/update", setup.UserController.UpdateUser)

	testCases := []struct {
		name         string
		userJson     string
		expectedCode int
	}{
		{"admin cannot operate with admin", `{"id": 1, "name": "Bob", "email": "bob@example.com", "roleId": 2}`, http.StatusInternalServerError},
		{"valid", `{"id": 2, "name": "Jane1", "email": "jane1@example.com", "roleId": 2}`, http.StatusOK},
		{"admin cannot change user role to admin", `{"id": 2, "name": "Jane1", "email": "jane1@example.com", "roleId": 1}`, http.StatusInternalServerError},
		{"might not take same name", `{"id": 2, "name": "Bob"}`, http.StatusInternalServerError},
		{"name already taken", `{"id": 1, "name": "Jane Doe"}`, http.StatusInternalServerError},
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

func TestUpdateUserAsSuperAdmin(t *testing.T) {
	setup := Setup(SuperAdminMiddleware())

	setup.Router.POST("/update", setup.UserController.UpdateUser)

	testCases := []struct {
		name         string
		userJson     string
		expectedCode int
	}{
		{"superadmin cannot operate with superadmin", `{"id": 1, "name": "Bob", "email": "bob@example.com", "roleId": 3}`, http.StatusInternalServerError},
		{"superadmin can operate with admin", `{"id": 3, "name": "Bob1", "email": "sadmin1@example.com", "roleId": 2}`, http.StatusOK},
		{"valid", `{"id": 106, "name": "suser", "email": "suser1@example.com", "roleId": 2}`, http.StatusOK},
		{"superadmin can change user role to admin", `{"id": 106, "name": "sadmin2", "email": "jane1@example.com", "roleId": 1}`, http.StatusOK},
		{"superadmin cannot change user role to superadmin", `{"id": 106, "name": "sadmin3", "email": "jane1@example.com", "roleId": 3}`, http.StatusInternalServerError},
		{"might not take same name", `{"id": 1, "name": "Bob"}`, http.StatusInternalServerError},
		{"name already taken", `{"id": 1, "name": "Jane Doe"}`, http.StatusInternalServerError},
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

func TestDeleteUserAsAdmin(t *testing.T) {
	setup := Setup(AdminMiddleware())

	setup.Router.DELETE("/delete/:id", setup.UserController.DeleteUser)

	testCases := []struct {
		name         string
		userId       string
		expectedCode int
	}{
		{"valid user", "2", http.StatusOK},
		{"invalid user", "abc", http.StatusBadRequest},
		{"nonexistent user", "9999", http.StatusInternalServerError},
		{"same id", "1", http.StatusInternalServerError},
		{"other admin", "3", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest("DELETE", "/delete/"+tc.userId, nil)
			response := httptest.NewRecorder()

			setup.Router.ServeHTTP(response, request)

			assert.Equal(t, tc.expectedCode, response.Code)
		})
	}
}

func TestDeleteUserAsSuperAdmin(t *testing.T) {
	setup := Setup(SuperAdminMiddleware())

	setup.Router.DELETE("/delete/:id", setup.UserController.DeleteUser)

	testCases := []struct {
		name         string
		userId       string
		expectedCode int
	}{
		{"valid user", "107", http.StatusOK},
		{"invalid user", "abc", http.StatusBadRequest},
		{"nonexistent user", "9999", http.StatusInternalServerError},
		{"same id", "105", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest("DELETE", "/delete/"+tc.userId, nil)
			response := httptest.NewRecorder()

			setup.Router.ServeHTTP(response, request)

			assert.Equal(t, tc.expectedCode, response.Code)
		})
	}
}
