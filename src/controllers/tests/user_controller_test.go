package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gokg/gomvc/controllers"
	"gokg/gomvc/repositories"
)

type TestSetup struct {
	Router         *gin.Engine
	UserController *controllers.UserController
}

func Setup() TestSetup {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	acl := repositories.NewAclAbstract()
	userController := &controllers.UserController{}
	userController.Init(acl)

	return TestSetup{
		Router:         router,
		UserController: userController,
	}
}

func TestInsertUser(t *testing.T) {
	setup := Setup()

	setup.Router.POST("/insert", setup.UserController.InsertUser)

	testCases := []struct {
		name         string
		userJson     string
		expectedCode int
	}{
		{"valid user", `{"name": "Alice", "email": "alice@example.com", "roleId": 1}`, http.StatusOK},
		{"invalid user", `{"name": "Alice", "email": "alice@example.com", "roleId": 1}`, http.StatusInternalServerError},
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

func TestUpdateUser(t *testing.T) {
	setup := Setup()

	setup.Router.POST("/update", setup.UserController.UpdateUser)

	testCases := []struct {
		name         string
		userJson     string
		expectedCode int
	}{
		{"valid user", `{"id": 1, "name": "Bob", "email": "bob@example.com", "roleId": 1}`, http.StatusOK},
		{"invalid user", `{"id": 2, "name": "Bob"}`, http.StatusInternalServerError},
		{"invalid user", `{"id": 1, "name": "Jane Doe"}`, http.StatusInternalServerError},
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

func TestDeleteUser(t *testing.T) {
	setup := Setup()

	setup.Router.DELETE("/delete/:id", setup.UserController.DeleteUser)

	testCases := []struct {
		name         string
		userId       string
		expectedCode int
	}{
		{"valid user", "1", http.StatusOK},
		{"invalid user", "abc", http.StatusBadRequest},
		{"nonexistent user", "9999", http.StatusInternalServerError},
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
