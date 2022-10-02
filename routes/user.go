package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"library_app/config"
	"library_app/models"
	"net/http"
)

type UserResponse struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     models.Role `json:"role"`
}

func GetUsers(c *gin.Context) {
	users := []models.User{}
	config.DB.Preload(clause.Associations).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success Get User List",
		"data":    users,
	})
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	data := config.DB.Preload(clause.Associations).Find(&user, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"data":    "Data Not Found",
		})
		c.Abort()
		return
	}

	respUser := UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success Get User By ID",
		"data":    respUser,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var userRequest models.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad Request",
		})

		c.Abort()
		return
	}

	if err := userRequest.HashPassword(userRequest.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Internal Server Error",
		})
	}

	var user models.User
	data := config.DB.First(&user, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"data":    "Data Not Found",
		})
		return
	}

	config.DB.Model(&user).Where("id = ?", id).Update("Name", userRequest.Name)
	config.DB.Model(&user).Where("id = ?", id).Update("Username", userRequest.Username)
	config.DB.Model(&user).Where("id = ?", id).Update("Email", userRequest.Email)
	config.DB.Model(&user).Where("id = ?", id).Update("Role_ID", userRequest.RoleID)
	config.DB.Model(&user).Where("id = ?", id).Update("Password", userRequest.Password)

	respUser := UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update Success",
		"data":    respUser,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	data := config.DB.First(&user, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"data":    "Data Not Found",
		})
		return
	}

	config.DB.Delete(&user, id)
	c.JSON(http.StatusNoContent, gin.H{
		"message": "Delete Success",
	})
}
