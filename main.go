package main

import (
	// "fmt"
	"bwastartup/auth"
	"bwastartup/handler"
	"bwastartup/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//buat koneksi ke database
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	/*proses input
	input -> handler mapping input ke struct
	-> service mapping ke struct user (model)
	->= repository save struct ke db
	*/
	//

	//panggil NewRepository dari repo user
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	authService := auth.NewService()
	fmt.Println(authService.GenerateToken(1))

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)
	api.POST("/email-checker", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UploadAvatar)
	router.Run()

}

