package marketplace

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	OrganizationName string `json:"organization_name"`
	QuicknodeID      string `json:"quicknode-id"`
}

type JWTClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

func GetJWT(secretKey string, user User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["organization_name"] = user.OrganizationName
	claims["quicknode_id"] = user.QuicknodeID

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
