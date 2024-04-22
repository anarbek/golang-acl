package controllers

import (
	"gokg/gomvc/models"
	"gokg/gomvc/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RolesController struct {
	acl *repositories.RoleBase
}

func (u *RolesController) Init(_acl *repositories.RoleBase) {
	_acl.Init()
	u.acl = _acl
}

// @Summary Get all roles
// @Description Get all roles from the database
// @ID get-all-roles
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} models.Role
// @Router /roles [get]
// @Tags Roles
func (u *RolesController) GetAll(c *gin.Context) {
	loggedInUser, ok := GetLoggedInUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert Role"})
		return
	}
	users := u.acl.Roles(loggedInUser)
	// implement your logic here
	c.JSON(200, users)
}

// @Summary Get a role by ID
// @Description Get a specific role from the database by ID
// @ID get-role-by-id
// @Produce  json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} models.Role
// @Router /roles/{id} [get]
// @Tags Roles
func (u *RolesController) GetRole(c *gin.Context) {
	loggedInUser, ok := GetLoggedInUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert Role"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := u.acl.GetRole(loggedInUser, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, role)
}

// @Summary Get permissions for logged in user
// @Description Get permissions for the logged in user from the database
// @ID get-permissions-for-loggedin-user
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} string
// @Router /roles/permissionsforuser [get]
// @Tags Roles
func (u *RolesController) GetPermissionsForLoggedInUser(c *gin.Context) {
	loggedInUser, ok := GetLoggedInUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert Role"})
		return
	}
	roles := u.acl.GetPermissionsForLoggedinUser(loggedInUser)
	// implement your logic here
	c.JSON(200, roles)
}

// @Summary Insert a role
// @Description Insert a new role into the database
// @ID insert-role
// @Accept  json
// @Produce  json
// @Param role body models.Role true "Role to add"
// @Security BearerAuth
// @Success 200 {object} models.Role
// @Router /roles/insert [post]
// @Tags Roles
func (u *RolesController) InsertRole(c *gin.Context) {
	loggedInUser, ok := GetLoggedInUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert Role"})
		return
	}
	// Parse the role from the request
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the user using the acl
	if err := u.acl.InsertRole(&role, loggedInUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// @Summary Update a role
// @Description Update an existing role in the database
// @ID update-role
// @Accept  json
// @Produce  json
// @Param role body models.Role true "Role to update"
// @Security BearerAuth
// @Success 200 {object} models.Role
// @Router /roles/update [post]
// @Tags Roles
func (u *RolesController) UpdateRole(c *gin.Context) {
	loggedInUser, ok := GetLoggedInUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert Role"})
		return
	}
	// Parse the user from the request
	var user models.Role
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user using the acl
	if err := u.acl.UpdateRole(&user, loggedInUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete a role
// @Description delete a role from the database
// @ID delete-role
// @Produce  json
// @Param id path int true "Role ID"
// @Success 200 "Role deleted"
// @Router /roles/delete/{id} [delete]
// @Security BearerAuth
// @Tags Roles
func (u *RolesController) DeleteRole(c *gin.Context) {
	loggedInUser, ok := GetLoggedInUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert Role"})
		return
	}
	// Get the user ID from the URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Role ID"})
		return
	}

	// Delete the user using the acl
	if err := u.acl.DeleteRole(id, loggedInUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Role deleted"})
}
