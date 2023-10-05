package token

import (
	"github.com/golang-jwt/jwt"
	"os"
	"strconv"
	"time"
)

type RefreshToken string

func GenerateRefreshToken(tokenUuid string, userId uint) (RefreshToken, error) {
	refreshTokenLifespan, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAY_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims[ClaimUserId] = userId
	claims[ClaimUUID] = tokenUuid
	claims["exp"] = time.Now().Add((time.Hour * 24) * time.Duration(refreshTokenLifespan)).Unix()
	claims[ClaimType] = Refresh

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := refreshToken.SignedString([]byte(refreshTokenKey()))
	return RefreshToken(signedToken), err
}

func (t RefreshToken) ValidateToken() (*jwt.Token, error) {
	token := Token(t)
	return token.validateToken(refreshTokenKey())
}

func refreshTokenKey() string {
	return os.Getenv("REFRESH_API_SECRET")
}
