package users

import (
	"errors"
	"fmt"
	"gokg/gomvc/repositories"
	"net/http"
	"strings"
	"time"

	"gokg/gomvc/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

type SignedResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// ILoginHandler defines the interface for user operations
type ILoginHandler interface {
	Init(_acl *repositories.AclBase)
	SetRand(_num int)
	//GetAll(c *gin.Context)
	LoginUser(c *gin.Context)
	PrivateACLCheckUserWrapper(pageName string, read, write bool) gin.HandlerFunc
}

type CustomLoginHandler struct {
	num    int
	acl    *repositories.AclBase
	secret string
}

func (h *CustomLoginHandler) Init(_acl *repositories.AclBase) {
	h.acl = _acl
	h.secret = "mk<ofPDZ34T5kJNII6'zV;DqSg_2V3HvW5`<-Kel*>sv/Nz)QIBL<%kq=#u'`eY"
}

func (h *CustomLoginHandler) SetRand(_num int) {
	h.num = _num
}

/*func (h *CustomLoginHandler) GetAll(c *gin.Context) {
	users := h.acl.UsersWithRoles()
	for _, user := range users {
		user.CurrNum = h.num
	}
	// implement your logic here
	c.JSON(200, users)
}*/

func (h *CustomLoginHandler) LoginUser(c *gin.Context) {
	/*type loginUser struct {
		Username string `json:"username,omitempty"`
	}*/

	loginParams := models.LoginUser{}
	c.ShouldBindJSON(&loginParams)

	/*users := h.acl.UsersWithRoles()
	// Search for the user in the Users array
	var user models.User
	for _, u := range users {
		if u.Username == loginParams.Username {
			user = u
			break
		}
	}*/

	user := h.acl.GetUserByUsernamePassword(loginParams)

	// If the user was found, generate a JWT and include the user's RoleID in the claims
	if user.ID != 0 {
		expirationTime := time.Now().Add(1 * time.Hour) // Token expires after 1 hours
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user":   user.Username,
			"userId": user.ID,
			"nbf":    time.Now().Unix(),     // Set 'nbf' to now
			"exp":    expirationTime.Unix(), // Add the 'exp' claim
		})

		tokenStr, err := token.SignedString([]byte(h.secret))
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

func (h *CustomLoginHandler) VerifyToken(c *gin.Context) {
	// Get the token from the Authorization header
	tokenString := c.GetHeader("Authorization")

	// Replace "Bearer " prefix (if it exists)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.secret), nil // replace with your secret key
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.JSON(http.StatusOK, gin.H{"user": claims["user"]})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	}
}

func (h *CustomLoginHandler) PrivateACLCheckUserWrapper(pageName string, read, write bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		privateACLCheckUser(c, h, pageName, read, write)
	}
}

func privateACLCheckUser(c *gin.Context, h *CustomLoginHandler, pageName string, read, write bool) {
	jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := parseToken(jwtToken, h.secret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
			Message: "bad jwt token",
		})
		return
	}

	claims, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
			Message: "unable to parse claims",
		})
		return
	}

	claimedUserID, OK := claims["userId"].(float64)
	if !OK {
		c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
			Message: "no userId property in claims",
		})
		return
	}

	claimedUserIDInt := int(claimedUserID)
	// Get the user from the Users slice
	var user models.User
	users := h.acl.UsersAll()
	for _, u := range users {
		if u.ID == claimedUserIDInt {
			user = u
			break
		}
	}

	// If the user was not found, return an error
	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, UnsignedResponse{
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
			if write && !rp.Write {
				// If write access is required but the user doesn't have it, deny access
				hasPolicy = false
				break
			}
			if read && !rp.Read {
				// If read access is required but the user doesn't have it, deny access
				hasPolicy = false
				break
			}
			// If none of the above conditions were met, the user has the required access
			hasPolicy = true
			break
		}
	}

	if !hasPolicy {
		c.AbortWithStatusJSON(http.StatusForbidden, UnsignedResponse{
			Message: "user does not have required policy",
		})
		return
	}
	c.Set("user", user)
	c.Next()
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given on extract")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func parseToken(jwtToken string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New("bad jwt token")
	}

	return token, nil
}
