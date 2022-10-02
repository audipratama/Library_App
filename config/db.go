package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"library_app/models"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/library_app?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&models.User{}, &models.Book{}, &models.Role{},&models.Borrow{})

	//Uncomment for DB Migration
	//DB.Create(
	//	&models.Book{
	//		Name:   "Ihya Ulumuddin",
	//		Writer: "Abu Hamid Al-Ghazali",
	//		Year:   1105,
	//		Owned:  2,
	//	},
	//)
	//
	//DB.Create(
	//	&models.Book{
	//		Name:   "Kitab Al-Hikam",
	//		Writer: "Ibnu Atha’illah As-Sakandari",
	//		Year:   1296,
	//		Owned:  1,
	//	},
	//)
	//
	//DB.Create(
	//	&models.Role{
	//		Name:  "Guest",
	//		Code:  "G",
	//	},
	//)
	//
	//DB.Create(
	//	&models.Role{
	//		Name:  "Admin",
	//		Code:  "A",
	//	},
	//)
}
