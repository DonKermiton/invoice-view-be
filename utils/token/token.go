package token

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt"
	uuidLib "github.com/google/uuid"
	"invoice-view-be/models"
	"strconv"
	"time"
)

const (
	Access      = "access_token"
	Refresh     = "refresh_token"
	ClaimUserId = "userId"
	ClaimUUID   = "tokenUuid"
	ClaimType   = "type"
)

type LoginTokenGroup struct {
	AccessToken  string
	RefreshToken string
}

type Token string

// ClaimsData share claims data between tokens
// tokenExp in hours
type ClaimsData struct {
	UserId    uint
	ClaimType string
	ClaimUuid string
	TokenExp  int
}

func GenerateTokensPair(userId uint) (LoginTokenGroup, error) {
	tokenUuid := uuidLib.New().String()

	accessToken, err := GenerateToken(tokenUuid, userId)

	if err != nil {
		return LoginTokenGroup{}, err
	}

	refreshToken, err := GenerateRefreshToken(tokenUuid, userId)

	if err != nil {
		return LoginTokenGroup{}, err
	}

	return LoginTokenGroup{accessToken, string(refreshToken)}, nil
}

func generateClaims(data ClaimsData) jwt.MapClaims {
	return jwt.MapClaims{
		ClaimUserId: data.UserId,
		ClaimUUID:   data.ClaimUuid,
		ClaimType:   data.ClaimType,
		"exp":       time.Now().Add(time.Hour * time.Duration(data.TokenExp)).Unix(),
	}
}

func ExtractTokenId(c *gin.Context) (uint, error) {
	token, err := ExtractTokenAndValidate(c)

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.f", claims[ClaimUserId]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}

func (t *Token) validateToken(privateKey string) (*jwt.Token, error) {
	return jwt.Parse(string(*t), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(privateKey), nil
	})
}

func (t *Token) CheckTokensAreInPair(userId uint, r *RefreshToken) (bool, error) {
	// todo:: check is Token valid
	validatedAccessToken, err := t.validateToken(accessTokenKey())

	if err != nil {
		errCode := HandleToken(err)

		if errCode != 0 || errCode != jwt.ValidationErrorExpired {
			return false, err
		}
	}

	accessTokenClaims := getClaims(validatedAccessToken)

	accessTokenUserId, err := parseClaimToInt(accessTokenClaims[ClaimUserId])
	accessTokenExp := accessTokenClaims["exp"]
	accessTokenUuid := accessTokenClaims[ClaimUUID]

	if err := accessTokenClaims.hasAccessTokenValidType(); err != nil {
		return false, err
	}

	user, err := models.GetUserById(userId)

	if err != nil {
		return false, errors.New("something went wrong with auth. Try logging again")
	}

	if user.IsUserDeleted() {
		return false, errors.New("user Deleted")
	}

	if user.HasBeenDataModified(accessTokenExp.(time.Time)) {
		return false, errors.New("user data updated after validatedAccessToken was generated. Try Login again")
	}

	// check refresh token
	refreshToken, err := r.ValidateToken()

	if !refreshToken.Valid {
		return false, errors.New("refresh token invalid")
	}

	if err != nil {
		return false, err
	}

	refreshTokenClaims := getClaims(refreshToken)

	refreshTokenUuid := refreshTokenClaims[ClaimUUID]

	if accessTokenUuid != refreshTokenUuid {
		return false, errors.New("tokens are not in pair")
	}

	refreshTokenUserId, err := parseClaimToInt(refreshTokenClaims[ClaimUserId])
	if err != nil {
		return false, err
	}

	refreshTokenType := refreshTokenClaims[ClaimType]

	if refreshTokenUserId != accessTokenUserId {
		return false, errors.New("different user id in access and refresh")
	}

	if refreshTokenType == Refresh {
		return false, errors.New("refresh token is in invalid type")
	}

	return true, nil
}

// returns true if token expire
func checkTokenExpired(token *jwt.Token) bool {
	return !token.Claims.(jwt.MapClaims).VerifyExpiresAt(time.Now().Unix(), false)
}

func getClaims(token *jwt.Token) AuthClaims {
	return AuthClaims(token.Claims.(jwt.MapClaims))
}
func parseClaimToInt(claim interface{}) (uint, error) {
	parsedClaim, err := strconv.ParseUint(fmt.Sprintf("%.f", claim), 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(parsedClaim), nil
}
