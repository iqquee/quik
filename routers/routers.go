package routers

import (
	"quik/controllers"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("api/v1/user")
	{
		auth.POST("/signup", controllers.CreateUser)
		auth.POST("/signin", controllers.SignIn)
		auth.GET("/users", controllers.GetUsers)
		auth.DELETE("/:id", controllers.DeleteUser)
	}

	wallet := r.Group("api/v1/wallets")
	{
		wallet.GET("/:wallet_id/balance", controllers.GetWalletBalance)
		wallet.POST("/:wallet_id/credit", controllers.CreditWallet)
		wallet.POST("/:wallet_id/debit", controllers.DebitWallet)
		wallet.GET("/wallets", controllers.GetWallets)
	}
	return r
}
