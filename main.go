package main

import (
	"github.com/gin-gonic/gin"
	"library_app/config"
	"library_app/middlewares"
	"library_app/routes"
)

func main() {
	config.InitDB()

	router := gin.Default()

	authentication := router.Group("/authentication")
	{
		authentication.POST("get-token", routes.GenerateToken)
		authentication.POST("register", routes.RegisterUser)
	}

	user := router.Group("/user").Use(middlewares.Auth())
	{
		user.GET("/", routes.GetUsers)
		user.GET("/:id", routes.GetUserByID)
		user.PUT("/:id", routes.UpdateUser)
		user.DELETE("/:id", routes.DeleteUser)
	}

	role := router.Group("/role").Use(middlewares.Auth())
	{
		role.GET("/", routes.GetRoles)
		role.GET("/:id", routes.GetRoleByID)
		role.POST("/", routes.InsertRole)
	}

	book := router.Group("/book").Use(middlewares.Auth())
	{
		book.GET("/", routes.GetBooks)
		book.GET("/:id", routes.GetBookByID)
		book.POST("/", routes.InsertBook)
		book.PUT("/:id", routes.UpdateBook)
		book.DELETE("/:id", routes.DeleteBook)
	}

	borrow := router.Group("/borrow").Use(middlewares.Auth())
	{
		borrow.GET("/", routes.GetBorrow)
		borrow.GET("/report", routes.GetBorrowReport)
		borrow.POST("/", routes.InsertBorrow)
		borrow.POST("/return/:id", routes.ReturnBorrowByID)
		borrow.GET("/exceed-deadline", routes.GetBorrowExceedDeadline)
	}

	router.Run()
}
