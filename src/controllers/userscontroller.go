package controllers

import (
	"example/hello/models"
	"example/hello/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	acl *repositories.AclAbstract
}

func (u *UserController) Init(_acl *repositories.AclAbstract) {
	u.acl = _acl
}

func (u *UserController) GetAll(c *gin.Context) {
	users := u.acl.UsersWithRoles()
	/*for _, user := range users {
		user.CurrNum = num
	}*/
	// implement your logic here
	c.JSON(200, users)
}

// Insert a new user
func (u *UserController) InsertUser(c *gin.Context) {
	// Parse the user from the request
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the user using the acl
	if err := u.acl.InsertUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update an existing user
func (u *UserController) UpdateUser(c *gin.Context) {
	// Parse the user from the request
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user using the acl
	if err := u.acl.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete a user
func (u *UserController) DeleteUser(c *gin.Context) {
	// Get the user ID from the URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Delete the user using the acl
	if err := u.acl.DeleteUser(id); err != nil {
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
