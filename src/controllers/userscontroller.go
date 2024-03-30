package controllers

import (
	"gokg/gomvc/models"
	"gokg/gomvc/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	acl *repositories.AclBase
}

func (u *UserController) Init(_acl *repositories.AclBase) {
	u.acl = _acl
}

// @Summary Get all users
// @Description get all users with their roles
// @ID get-all-users
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get]
// @Security BearerAuth
func (u *UserController) GetAll(c *gin.Context) {
	/*loggedInUser, ok := GetLoggedInUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user"})
		return
	}
	users := u.acl.Users(loggedInUser)*/
	users := u.acl.UsersAll()
	/*for _, user := range users {
		user.CurrNum = num
	}*/
	// implement your logic here
	c.JSON(200, users)
}

func GetLoggedInUser(c *gin.Context) (value *models.User, ok bool) {
	loggedInUserInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user in context"})
		ok = false
	}
	// Perform a type assertion to convert loggedInUser to *models.User
	loggedInUser, ok := loggedInUserInterface.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user"})
		ok = false
	}
	value = &loggedInUser
	return value, ok
}

// @Summary Insert a new user
// @Description insert a new user into the database
// @ID insert-user
// @Accept  json
// @Produce  json
// @Param user body models.User true "user to insert"
// @Success 200 {object} models.User
// @Router /users/insert [post]
// @Security BearerAuth
func (u *UserController) InsertUser(c *gin.Context) {
	loggedInUser, ok := GetLoggedInUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user"})
		return
	}
	// Parse the user from the request
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the user using the acl
	if err := u.acl.InsertUser(&user, loggedInUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update an existing user
// @Description update an existing user in the database
// @ID update-user
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body models.User true "user to update"
// @Success 200 {object} models.User
// @Router /users/update [post]
// @Security BearerAuth
func (u *UserController) UpdateUser(c *gin.Context) {
	loggedInUser, ok := GetLoggedInUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user"})
		return
	}
	// Parse the user from the request
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user using the acl
	if err := u.acl.UpdateUser(&user, loggedInUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete a user
// @Description delete a user from the database
// @ID delete-user
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 "User deleted"
// @Router /users/delete/{id} [delete]
// @Security BearerAuth
func (u *UserController) DeleteUser(c *gin.Context) {
	loggedInUser, ok := GetLoggedInUser(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user"})
		return
	}
	// Get the user ID from the URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Delete the user using the acl
	if err := u.acl.DeleteUser(id, loggedInUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "User deleted"})
}

/*var num = 0

func SetRand(_num int) {
	num = _num
}

func GetAll(c *gin.Context) {
	users := repositories.UsersWithRoles()
	for _, user := range users {
		user.CurrNum = num
	}
	// implement your logic here
	c.JSON(200, users)
}

func loginUser(c *gin.Context) {

}
*/
