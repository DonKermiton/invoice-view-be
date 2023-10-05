package token

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"os"
	"strconv"
)

type AccessToken string
type AuthClaims jwt.MapClaims

func GenerateToken(tokenUuid string, userId uint) (string, error) {
	tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := generateClaims(ClaimsData{
		UserId:    userId,
		ClaimType: Access,
		ClaimUuid: tokenUuid,
		TokenExp:  tokenLifespan,
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(accessTokenKey()))
}

func accessTokenKey() string {
	return os.Getenv("API_SECRET")
}

func IsTokenValid(c *gin.Context) error {
	validatedToken, err := ExtractTokenAndValidate(c)

	if err != nil {
		return err
	}

	expired := checkTokenExpired(validatedToken)

	if expired {
		return errors.New("token expired")
	}

	return nil
}

func ExtractToken(c *gin.Context) (Token, error) {
	bearerToken, err := c.Cookie("access_token")

	if err != nil {
		fmt.Println(err)
	}

	return Token(bearerToken), nil
}

func ExtractTokenAndValidate(c *gin.Context) (*jwt.Token, error) {
	t, err := ExtractToken(c)

	if err != nil {
		return nil, err
	}

	token, err := t.validateToken(accessTokenKey())

	if err != nil {
		return nil, err
	}

	return token, nil
}

func HandleToken(err error) uint32 {

	if err != nil {
		validationError, _ := err.(*jwt.ValidationError)

		return validationError.Errors
	}

	return 0
}

func (t *AuthClaims) hasAccessTokenValidType() error {
	if (*t)[ClaimType] == Access {
		return nil
	}

	return errors.New("access token has invalid claim type")
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTYyNTc4ODAsInRva2VuVXVpZCI6IjJiYjZlNWEyLTU0OGQtNGI2ZC05ZjA5LTIyY2YxOGE3ZmNhOCIsInR5cGUiOiJhY2Nlc3NfdG9rZW4iLCJ1c2VySWQiOjF9.zbDEe_Vo_gSlo1RhTGrXmkFy5naIxrug5PSItp1nXnw
