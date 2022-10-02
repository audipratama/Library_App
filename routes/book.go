package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"library_app/config"
	"library_app/models"
	"net/http"
)

type BookRequest struct {
	Name   string `json:"name"`
	Writer string `json:"writer"`
	Year   int    `json:"year"`
	Owned  int    `json:"owned"`
}

type BooksResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Writer string `json:"writer"`
	Year   int    `json:"year"`
	Owned  int    `json:"owned"`
}

func GetBooks(c *gin.Context) {
	books := []models.Book{}

	config.DB.Preload(clause.Associations).Find(&books)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success Get Book List",
		"data":    books,
	})
}

func GetBookByID(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	data := config.DB.Preload(clause.Associations).First(&book, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"data":    "Data Not Found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success Get Book By ID",
		"data":    book,
	})
}

func InsertBook(c *gin.Context) {
	var request BookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad Request",
		})

		c.Abort()
		return
	}

	book := models.Book{
		Name:   request.Name,
		Writer: request.Writer,
		Year:   request.Year,
		Owned:  request.Owned,
	}

	config.DB.Create(&book)

	respBook := BooksResponse{
		ID:     book.ID,
		Name:   book.Name,
		Year:   book.Year,
		Writer: book.Writer,
		Owned:  book.Owned,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    respBook,
		"message": "Insert Success",
	})
}

func UpdateBook(c *gin.Context) {
	var request BookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Bad Request",
		})

		c.Abort()
		return
	}

	id := c.Param("id")
	var book models.Book
	data := config.DB.First(&book, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"data":    "Data Not Found",
		})
		return
	}

	var count int64
	var borrowData models.Borrow
	config.DB.Model(&borrowData).Where("book_id = ? and return_date is null", book.ID).Count(&count)
	if request.Owned < int(count) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Can't change book stock lower than active borrow",
		})
		return
	}

	config.DB.Model(&book).Where("id = ?", id).Update("Name", request.Name)
	config.DB.Model(&book).Where("id = ?", id).Update("Writer", request.Writer)
	config.DB.Model(&book).Where("id = ?", id).Update("Year", request.Year)
	config.DB.Model(&book).Where("id = ?", id).Update("Owned", request.Owned)

	respDept := BooksResponse{
		ID:     book.ID,
		Name:   book.Name,
		Year:   book.Year,
		Writer: book.Writer,
		Owned:  book.Owned,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update Success",
		"data":    respDept,
	})
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	data := config.DB.First(&book, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Data Not Found",
			"data":    "Data Not Found",
		})
		return
	}

	config.DB.Delete(&book, id)
	c.JSON(http.StatusNoContent, gin.H{
		"message": "Delete Success",
	})
}
