package middleware

import (
	"fmt"
	"github.com/labstack/echo"
	microContext "golang.org/x/net/context"
	"net/http"
	proto "protos"
	"strings"
)

type AuthMiddleware interface {
	Auth(next echo.HandlerFunc) echo.HandlerFunc
}

type authMiddleware struct {
	authMicroservice proto.AuthClient
}

func NewAuth(authMicroservice proto.AuthClient) AuthMiddleware {
	return &authMiddleware{authMicroservice}
}

func (a *authMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		httpErr := echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		rawToken := c.Request().Header.Get("Authorization")
		token := strings.TrimPrefix(rawToken, "Bearer ")

		if len(token) < 44 {
			return httpErr
		}

		response, err := a.authMicroservice.ProlongSession(microContext.TODO(), &proto.AuthTokenData{Token: &token})

		id := response.GetId()

		if err != nil || id == 0 {
			fmt.Println(err)
			return httpErr
		}

		// temporary
		c.Set("userId", id)
		c.Set("token", token)

		return next(c)
	}
}
