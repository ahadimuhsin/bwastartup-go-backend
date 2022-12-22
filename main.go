package main

import (
	// "fmt"
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"path/filepath"

	// "bwastartup/transaction"
	"bwastartup/user"

	webHandler "bwastartup/web/handler"

	// "fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

	//panggil NewRepository dari repo campaign
	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)

	//module transaction
	transactionRepository := transaction.NewRepository(db)

	//paymentService
	paymentService := payment.NewService()

	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)
	// campaigns, _ := campaignService.GetCampaigns(2);
	// fmt.Println(len(campaigns));
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	//bagian web
	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService)

	router := gin.Default()
	router.Use(cors.Default())
	// router.LoadHTMLGlob("web/templates/**/*")
	router.HTMLRender = loadTemplates("./web/templates")

	router.Static("/images", "./images")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")
	api := router.Group("/api/v1")
	// route authentication
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)
	api.POST("/email-checker", userHandler.CheckEmailAvailability)
	api.POST("/upload-avatar", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	//route campaign
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.POST("/campaign", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.GET("/campaign/:id", campaignHandler.GetCampaign)
	api.PUT("/campaign/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	//route transactions
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transaction", authMiddleware(authService, userService), transactionHandler.CreateTransaction)

	//router web
	router.GET("/users", userWebHandler.Index)
	router.POST("/users", userWebHandler.Create)
	router.GET("/users/new", userWebHandler.New)
	router.GET("/users/edit/:id", userWebHandler.Edit)
	router.POST("/users/update/:id", userWebHandler.Update)
	router.GET("/users/avatar/:id", userWebHandler.NewAvatar)
	router.POST("/users/avatar/:id", userWebHandler.UploadAvatar)

	// route campaign
	router.GET("/campaigns", campaignWebHandler.Index)
	router.Run()

}

// *middleware auth
// butuh variabel autService dan userService agar bisa menggunakan kedua service tersebut
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		//ambil nilai header Authorization: Bearer token
		authHeader := c.GetHeader("Authorization")

		//cek apakah di dalam string authHeader, ada kata Bearer
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		//ambil token dengan memisahkan authHeader berdasarkan spasi
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		//validasi token
		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		//saat sukses, ambil user_id
		//ambil user dari db berdasarkan user_id lewat service
		userId := int(claim["user_id"].(float64))

		user, err := userService.GetUserById(userId)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("currentUser", user)
		//set context isinya user
	}
}

// load Templates
func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
