package helper

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type JwtHelper struct{}

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

	token := c.GetReqHeaders()["Authorization"]
	resp := ResponseHelper{}

	if token == "" {
		return resp.Unauthorized(c, "Access unauthorized")
	}
	token = strings.Split(token, " ")[1]

	// fix this, still error
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("asdfasdfasdf"), nil
	})

	if err == nil {
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
