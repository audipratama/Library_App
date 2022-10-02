package routes

import (
	"github.com/gin-gonic/gin"
	"library_app/auth"
	"library_app/config"
	"library_app/models"
	"net/http"
)

type TokenRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad Request",
		})

		c.Abort()
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Internal Server Error",
		})
	}

	insertUser := config.DB.Create(&user)
	if insertUser.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   insertUser.Error.Error(),
			"message": "Internal Server Error",
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"UserID":   user.ID,
		"Email":    user.Email,
		"username": user.Username,
		"RoleID": user.RoleID,
	})
}

func GenerateToken(c *gin.Context) {
	var request TokenRequest
	var user models.User

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad Request",
		})

		c.Abort()
		return
	}

	checkEmail := config.DB.Where("email = ?", request.Email).First(&user)
	if checkEmail.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"messages":"Data Not Found",
			"error": checkEmail.Error.Error(),
		})

		c.Abort()
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"messages":"Invalid Credential",
			"error": credentialError.Error(),
		})
	}

	tokenString, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"messages":"Invalid Credential",
			"error": err.Error(),
		})

		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":tokenString,
	})
}
