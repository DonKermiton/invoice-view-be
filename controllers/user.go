package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"invoice-view-be/models"
	"invoice-view-be/utils/token"
	"net/http"
	"os"
	"strconv"
)

type RegisterInput struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CompanyName string `json:"companyName" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company := models.Company{
		Name: input.CompanyName,
	}

	user := models.User{
		Email:    input.Email,
		Password: input.Password,
	}

	dbUser, err := models.GetUserByEmail(user.Email)

	if dbUser.Email != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "Account with provided Email already exists"})
		return
	}
	fmt.Println("teraz \n", &user.Password)
	_, _, err = company.CreateCompanyWithUser(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	setAndSendAuthCookies(c, user)
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: input.Password,
	}

	setAndSendAuthCookies(c, user)
}

func setAndSendAuthCookies(c *gin.Context, user models.User) {
	tokensPair, user, err := LoginCheck(user.Email, user.Password)
	user.Password = "_"

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect"})
		return
	}

	tokenLifeSpan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
	}

	refreshTokenLifeSpan, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAY_LIFESPAN"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Try again later"})
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    tokensPair.AccessToken,
		Path:     "localhost:8080",
		HttpOnly: true,
		Secure:   false,                // set to true in production
		MaxAge:   tokenLifeSpan * 3600, // 1 hour
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokensPair.RefreshToken,
		Path:     "localhost:8080",
		HttpOnly: true,
		Secure:   false, // set to true in production
		MaxAge:   refreshTokenLifeSpan * 3600 * 24,
	})

	c.JSON(http.StatusOK, gin.H{"data": user})
}

type RefreshTokenInput struct {
	RefreshToken string `json:"token" binding:"required"`
}

func RefreshToken(c *gin.Context) {
	var input RefreshTokenInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func CurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := models.GetUserById(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": user})
}

func LoginCheck(email string, password string) (token.LoginTokenGroup, models.User, error) {
	user, err := models.GetUserByEmail(email)

	if err != nil {
		return token.LoginTokenGroup{}, models.User{}, err
	}

	err = models.VerifyPassword(password, user.Password)

	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return token.LoginTokenGroup{}, models.User{}, err
	}

	tokensPair, err := token.GenerateTokensPair(user.ID)
	if err != nil {
		return token.LoginTokenGroup{}, models.User{}, errors.New("something went wrong. Try again later")
	}

	return tokensPair, user, err
}
