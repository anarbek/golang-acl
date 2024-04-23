package users

import (
	"gokg/gomvc/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type InitialLoginHandler struct {
	num int
	acl *repositories.AclAbstract
}

func (h *InitialLoginHandler) Init(_acl *repositories.AclAbstract) {
	h.acl = _acl
}

func (h *InitialLoginHandler) SetRand(_num int) {
	h.num = _num
}

/*func (h *InitialLoginHandler) GetAll(c *gin.Context) {
	users := h.acl.UsersWithRoles()
	for _, user := range users {
		user.CurrNum = h.num
	}
	// implement your logic here
	c.JSON(200, users)
}*/

func (h *InitialLoginHandler) LoginUser(c *gin.Context) {
	type login struct {
		Username string `json:"username,omitempty"`
	}

	loginParams := login{}
	c.ShouldBindJSON(&loginParams)

	if loginParams.Username == "mike" || loginParams.Username == "rama" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": loginParams.Username,
			"nbf":  time.Date(2018, 01, 01, 12, 0, 0, 0, time.UTC).Unix(),
		})

		tokenStr, err := token.SignedString([]byte("supersaucysecret"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, UnsignedResponse{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, SignedResponse{
			Token:   tokenStr,
			Message: "logged in",
		})
		return
	}

	c.JSON(http.StatusBadRequest, UnsignedResponse{
		Message: "bad username",
	})
}

func (h *InitialLoginHandler) PrivateACLCheckUserWrapper(pageName string, read, write bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: err.Error(),
			})
			return
		}

		token, err := parseToken(jwtToken, "supersaucysecret")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: "bad jwt token",
			})
			return
		}

		claims, OK := token.Claims.(jwt.MapClaims)
		if !OK {
			c.AbortWithStatusJSON(http.StatusInternalServerError, UnsignedResponse{
				Message: "unable to parse claims",
			})
			return
		}

		claimedUID, OK := claims["user"].(string)
		if !OK {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: "no user property in claims",
			})
			return
		}

		uid := c.Param("uid")
		if claimedUID != uid {
			c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
				Message: "token uid does not match resource uid",
			})
			return
		}

		c.Next()
	}
}
