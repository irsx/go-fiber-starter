package middlewares

import (
	"errors"
	"go-fiber-starter/utils"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	Issuer string `json:"issuer"`
	jwt.StandardClaims
}

type SkipperRoutesData struct {
	Method  string
	UrlPath string
}

func JwtMiddleware(ctx *fiber.Ctx) error {
	ctx.Set("X-XSS-Protection", "1; mode=block")
	ctx.Set("Strict-Transport-Security", "max-age=5184000")
	ctx.Set("X-DNS-Prefetch-Control", "off")

	// skip whitelist routes
	for _, whiteList := range whiteListRoutes() {
		if ctx.Method() == whiteList.Method && ctx.Path() == whiteList.UrlPath {
			return ctx.Next()
		}
	}

	// check header token
	authorizationToken := getAuthorizationToken(ctx)
	if authorizationToken == "" {
		err := errors.New("missing Bearer token")
		return utils.JsonErrorUnauthorized(ctx, err)
	}

	// verify token
	jwtToken, err := jwt.ParseWithClaims(authorizationToken, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return utils.JsonErrorUnauthorized(ctx, err)
	}

	claimsData := jwtToken.Claims.(*JwtCustomClaims)
	utils.Logger.Info("✅ SET USER AUTH")
	ctx.Locals("user_auth", claimsData.Issuer)

	return ctx.Next()
}

func getAuthorizationToken(ctx *fiber.Ctx) string {
	authorizationToken := string(ctx.Request().Header.Peek("Authorization"))
	return strings.Replace(authorizationToken, "Bearer ", "", 1)
}

func whiteListRoutes() []SkipperRoutesData {
	return []SkipperRoutesData{
		{"POST", "/api/auth/login"},
		{"POST", "/api/auth/register"},
		{"GET", "/api/import/products/stream"},
	}
}
