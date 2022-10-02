package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"library_app/auth"
	"library_app/config"
	"library_app/models"
	"net/http"
	"strconv"
	"time"
)

type BorrowRequest struct {
	DeadlineDate string `json:"deadline_date"`
	BookID       uint   `json:"book_id"`
}

type BorrowReport struct {
	BookID   int `json:"book_id"`
	Borrowed int
	Users    int
}

type BorrowReportResponse struct {
	BookName string `json:"book_name"`
	Borrowed int    `json:"borrowed_times"`
	Users    int    `json:"unique_users"`
}

type BorrowResponse struct {
	ID           uint   `json:"id"`
	CreatedAt    string `json:"created_at"`
	DeadlineDate string `json:"deadline_date"`
	ReturnDate   string `json:"return_date"`
	BookID       uint   `json:"book_id"`
	UserID       uint   `json:"user_id"`
}

func GetBorrow(c *gin.Context) {
	returned := c.Query("returned")
	borrows := []models.Borrow{}
	if returned != "" {
		returnBool, err := strconv.ParseBool(returned)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error parse returned param",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		if returnBool {
			config.DB.Preload(clause.Associations).Find(&borrows, "return_date is not null")
		} else {
			config.DB.Preload(clause.Associations).Find(&borrows, "return_date is null")
		}
	} else {
		config.DB.Preload(clause.Associations).Find(&borrows)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success Get Borrow List",
		"data":    borrows,
	})
}

func GetBorrowReport(c *gin.Context) {
	var report []BorrowReport
	data := config.DB.Model(&models.Borrow{}).Select("book_id, COUNT(book_id) as borrowed, COUNT(DISTINCT user_id) as users").Group("book_id").Find(&report)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"messages": "User Data Not Found",
			"error":    data.Error.Error(),
		})

		c.Abort()
		return
	}

	borrowedBookResp := []BorrowReportResponse{}
	for _, val := range report {
		var book models.Book
		config.DB.First(&book, "id = ?", val.BookID)

		resp := BorrowReportResponse{
			BookName: book.Name,
			Borrowed: val.Borrowed,
			Users:    val.Users,
		}
		borrowedBookResp = append(borrowedBookResp, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    borrowedBookResp,
		"message": "Get Borrow Report Success",
	})
}

func InsertBorrow(c *gin.Context) {
	var request BorrowRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad Request",
		})

		c.Abort()
		return
	}

	var count int64
	var borrowData models.Borrow
	config.DB.Model(&borrowData).Where("book_id = ? and return_date is null", request.BookID).Count(&count)

	var book models.Book
	config.DB.First(&book, "id = ?", request.BookID)
	if count >= int64(book.Owned) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Book can't be borrowed. Out of stock",
		})
		c.Abort()
		return
	}

	var user models.User
	username := auth.UserName
	userData := config.DB.Where("username = ?", username).First(&user)
	if userData.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"messages": "User Data Not Found",
			"error":    userData.Error.Error(),
		})

		c.Abort()
		return
	}

	deadlineDate, err := time.Parse("2006-01-02 15:04:05", request.DeadlineDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error parsing",
		})
		c.Abort()
		return
	}

	borrow := models.Borrow{
		DeadlineDate: deadlineDate,
		ReturnDate:   sql.NullTime{},
		BookID:       request.BookID,
		UserID:       user.ID,
	}

	config.DB.Create(&borrow)

	respBook := BorrowResponse{
		ID:           borrow.ID,
		CreatedAt:    borrow.Model.CreatedAt.Format(time.RFC3339),
		DeadlineDate: borrow.DeadlineDate.Format(time.RFC3339),
		BookID:       borrow.BookID,
		UserID:       borrow.UserID,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    respBook,
		"message": "Insert Borrow Success",
	})
}

func ReturnBorrowByID(c *gin.Context) {
	id := c.Param("id")
	var borrow models.Borrow
	data := config.DB.First(&borrow, "id = ? and return_date is null", id)
	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"data":    "Data Not Found",
		})
		c.Abort()
		return
	}

	config.DB.Model(&borrow).Where("id = ?", id).Update("return_date", sql.NullTime{Time: time.Now(), Valid: true})

	respBorrow := BorrowResponse{
		ID:           borrow.ID,
		CreatedAt:    borrow.Model.CreatedAt.Format(time.RFC3339),
		DeadlineDate: borrow.DeadlineDate.Format(time.RFC3339),
		ReturnDate:   borrow.ReturnDate.Time.Format(time.RFC3339),
		BookID:       borrow.BookID,
		UserID:       borrow.UserID,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    respBorrow,
		"message": "Return Borrow Success",
	})
}

func GetBorrowExceedDeadline(c *gin.Context) {
	borrows := []models.Borrow{}
	data := config.DB.Where("? > deadline_date and return_date is null", time.Now()).Find(&borrows)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"data":    "Data Not Found",
		})
		c.Abort()
		return
	}
	borrowResp := []BorrowResponse{}

	for _, val := range borrows {
		borrow := BorrowResponse{
			ID:           val.ID,
			CreatedAt:    val.CreatedAt.Format(time.RFC3339),
			DeadlineDate: val.DeadlineDate.Format(time.RFC3339),
			BookID:       val.BookID,
			UserID:       val.UserID,
		}
		borrowResp = append(borrowResp, borrow)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "These are the persons and books that exceed deadline",
		"data":    borrowResp,
	})
}
