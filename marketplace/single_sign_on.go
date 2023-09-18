package marketplace

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	OrganizationName string `json:"organization_name"`
	QuicknodeID      string `json:"quicknode-id"`
	Plan      			 string `json:"plan"`
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

func OpenDashboard(url string) (int, string, error) {
	client := &http.Client{}

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return 0, "", err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bodyStr := fmt.Sprintf("%s", body)

	if res.StatusCode != http.StatusOK {
		return res.StatusCode, bodyStr, fmt.Errorf("HTTP Request failed with status code: %d", res.StatusCode)
	} else {
		return res.StatusCode, bodyStr, nil
	}
}
