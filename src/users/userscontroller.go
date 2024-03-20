package users

import (
	"example/hello/repositories"

	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	users := repositories.UsersWithRoles()
	// implement your logic here
	c.JSON(200, users)
}
