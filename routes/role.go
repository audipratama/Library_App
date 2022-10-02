package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"library_app/config"
	"library_app/models"
	"net/http"
)

type RoleResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Code  string `json:"code"`
	Users []UserResponse
}

type RoleRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func GetRoleByID(c *gin.Context) {
	id := c.Param("id")

	var role models.Role
	data := config.DB.Preload(clause.Associations).Find(&role, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"data":    "Data Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success Get Role By ID",
		"data":    role,
	})
}

func GetRoles(c *gin.Context) {
	roles := []models.Role{}
	config.DB.Preload(clause.Associations).Find(&roles)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success Get Role List",
		"data":    roles,
	})
}

func InsertRole(c *gin.Context) {
	var request RoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad Request",
		})

		c.Abort()
		return
	}

	role := models.Role{
		Name: request.Name,
		Code: request.Code,
	}
	config.DB.Create(&role)

	respRole := RoleResponse{
		ID:   role.ID,
		Name: role.Name,
		Code: role.Code,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    respRole,
		"message": "Insert Success",
	})
}
