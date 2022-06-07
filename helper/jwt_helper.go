package helper

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type JwtHelper struct{}

type ClaimsData struct {
	UserId    string    `json:"userId"`
	Username  string    `json:"username"`
	UserRole  string    `json:"userRole"`
	Issued_at time.Time `json:"issued_at"`
	jwt.StandardClaims
}

func (jh *JwtHelper) GenerateToken(data map[string]interface{}) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    data["userId"],
		"username":  data["username"],
		"userRole":  data["userRole"],
		"issued_at": time.Now(),
	})

	tokenStr, err := token.SignedString([]byte("asdfasdfasdf"))

	return tokenStr, err
}

func (jh *JwtHelper) VerifyToken(c *fiber.Ctx) error {

	tokenStr := c.GetReqHeaders()["Authorization"]
	resp := new(ResponseHelper)

	if tokenStr == "" {
		return resp.Unauthorized(c, "Access Unauthorized")
	}
	tokenStr = strings.Split(tokenStr, " ")[1]

	// fix this, still error
	token, err := jwt.ParseWithClaims(tokenStr, &ClaimsData{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("asdfasdfasdf"), nil
	})

	if err == nil {
		data, _ := token.Claims.(*ClaimsData)
		c.Locals("user", data)
		return c.Next()
	} else if valErr, ok := err.(*jwt.ValidationError); ok {
		if valErr.Errors&jwt.ValidationErrorMalformed != 0 {
			return resp.ServerError(c, "Token unreadable")
		} else if valErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return resp.ServerError(c, "Token expired or not active")
		} else {
			return resp.ServerError(c, "Error: "+err.Error())
		}
	} else {
		return resp.ServerError(c, "Error: "+err.Error())
	}

}
