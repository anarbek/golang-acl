package users

import (
	"example/hello/repositories"

	"github.com/gin-gonic/gin"
)

var num = 0

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
