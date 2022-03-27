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
	}

	wallet := r.Group("api/v1/wallets")
	{
		wallet.GET("/:wallet_id/balance", controllers.GetWalletBalance)
		wallet.POST("/:wallet_id/credit", controllers.CreditWallet)
		wallet.POST("/:wallet_id/debit", controllers.DebitWallet)
	}
	return r
}
