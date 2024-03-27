package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"example/hello/repositories"
	"example/hello/users"
	"math/rand"
)

var privateThings = map[string]map[int64]string{
	"mike": {
		0: "MIKE: private string",
		1: "MIKE: secret thing",
		2: "MIKE: sneaky secret",
	},
	"rama": {
		0: "RAMA: private string",
		1: "RAMA: secret thing",
		2: "RAMA: sneaky secret",
	},
}

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

type SignedResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func index(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "indexMain"})
}

func private(c *gin.Context) {
	uidStr := c.Param("uid")
	pidInt, _ := strconv.ParseInt(c.Param("pid"), 10, 64)

	secret, ok := privateThings[uidStr][pidInt]

	if ok {
		c.JSON(200, gin.H{"msg": secret})
		return
	}

	c.JSON(200, gin.H{"msg": "unknown pid"})
}

func login(c *gin.Context) {
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

/*
func loginUser(c *gin.Context) {
	type loginUser struct {
		Username string `json:"username,omitempty"`
	}

	loginParams := loginUser{}
	c.ShouldBindJSON(&loginParams)

	users := repositories.UsersWithRoles()
	// Search for the user in the Users array
	var user models.User
	for _, u := range users {
		if u.Name == loginParams.Username {
			user = u
			break
		}
	}

	// If the user was found, generate a JWT and include the user's RoleID in the claims
	if user.ID != 0 {
		expirationTime := time.Now().Add(1 * time.Hour) // Token expires after 1 hours
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user":   user.Name,
			"userId": user.ID,
			"nbf":    time.Now().Unix(),     // Set 'nbf' to now
			"exp":    expirationTime.Unix(), // Add the 'exp' claim
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
}*/

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func parseToken(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte("supersaucysecret"), nil
	})

	if err != nil {
		return nil, errors.New("bad jwt token")
	}

	return token, nil
}

func jwtTokenCheck(c *gin.Context) {
	jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := parseToken(jwtToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: "bad jwt token",
		})
		return
	}

	_, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		c.AbortWithStatusJSON(http.StatusInternalServerError, UnsignedResponse{
			Message: "unable to parse claims",
		})
		return
	}
	c.Next()
}

func privateACLCheck(c *gin.Context) {
	jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := parseToken(jwtToken)
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

/*
func privateACLCheckUserWrapper(pageName string, read, write bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		privateACLCheckUser(c, pageName, read, write)
	}
}

func privateACLCheckUser(c *gin.Context, pageName string, read, write bool) {
	jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := parseToken(jwtToken)
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

	claimedUserID, OK := claims["userId"].(float64)
	if !OK {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: "no userId property in claims",
		})
		return
	}

	claimedUserIDInt := int(claimedUserID)
	// Get the user from the Users slice
	var user models.User
	for _, u := range repositories.Users {
		if u.ID == claimedUserIDInt {
			user = u
			break
		}
	}

	// If the user was not found, return an error
	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: "user not found",
		})
		return
	}

	// Get the user's role from the Roles slice
	var role models.Role = repositories.GetRole(user.RoleID)

	// Check if the role has the required policy
	hasPolicy := false
	for _, rp := range role.RolePolicies {
		if rp.Policy.Name == pageName {
			if write && rp.Write {
				hasPolicy = true
				break
			}
			if read && rp.Read {
				hasPolicy = true
				break
			}
		}
	}

	if !hasPolicy {
		c.AbortWithStatusJSON(http.StatusForbidden, UnsignedResponse{
			Message: "user does not have required policy",
		})
		return
	}

	c.Next()
}*/

func main() {
	handler := &users.MyUserHandler{}
	acl := repositories.AclAbstract{}
	handler.Init(&acl)
	number = rand.Intn(100) // Generate a random number between 0 and 99
	router := gin.New()
	router.GET("/", index)
	router.GET("/rnd", func(c *gin.Context) {
		c.String(http.StatusOK, "(2)Random number: %d", number)
	})

	router.POST("/login", login)
	router.POST("/loginUser", handler.LoginUser)
	/*router.Use(handler.PrivateACLCheckUserWrapper("UserManagement", true, false)).GET("/rndAuth", func(c *gin.Context) {
		c.String(http.StatusOK, "(2)Random number: %d", number)
	})*/

	privateRouter := router.Group("/private")
	privateRouter.Use(jwtTokenCheck)
	privateRouter.Use(privateACLCheck).GET("/:uid/:pid", private)

	handler.SetRand(number)
	v1 := router.Group("/api/v1")
	{
		usersRoutes := v1.Group("/users")
		{
			usersRoutes.Use(handler.PrivateACLCheckUserWrapper("UserManagement", true, false)).GET("/", handler.GetAll)
		}
		subjectRoutes := v1.Group("/subjects")
		{
			subjectRoutes.Use(handler.PrivateACLCheckUserWrapper("SubjectManagement", true, false)).GET("/", handler.GetAll)
		}
	}

	router.Run(":8081")
}

var number int

func init() {
	rand.Seed(time.Now().UnixNano())
}
