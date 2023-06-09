package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil{
		log.Fatal(err.Error())
	}

	fmt.Println("Connection to database is good")

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	
	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)
	api.GET("/campaigns",campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id",campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService),campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService),campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService),campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.Run()


	// gambaran struktur flow:
		// request
		// handler
		// Service
		// repository
		// database

	// gambaran flow register:
		// input
		// handler mapping input dari user -> struct input
		// Service mapping dari struct input ke struct User
		// repository save struct user ke db
		// database
	
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func (c *gin.Context){
		authHeader := c.GetHeader("Authorization")
	
		if !strings.Contains(authHeader, "Bearer"){
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}
	
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
	
		token, err := authService.ValidateToken(tokenString)

		if err != nil{
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil{
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized,response)
			return
		}
		c.Set("currentUser",user)
	}
}
// ambil nilai header authorization: bearer tokentokentoken
// dari header authorization, kita ambil nilai tokkennya saja
// kita validasi token
// kita ambil user_id
// ambil user dari db berdasarkan user_id lewat service
// kalau user ada set context isinya user
