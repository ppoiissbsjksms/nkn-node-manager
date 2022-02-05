package main

import (
	"github.com/gin-gonic/gin"
	"nkn-node-manager/controllers"
	"nkn-node-manager/models"
	"nkn-node-manager/utils"
)

func main() {

	r := gin.Default()

	models.ConnectDatabase() // new
	defer models.DB.Close()

	go utils.CheckOffline()

	r.GET("/wallets", controllers.FindWallets) // only when you want to export all your wallets
	r.GET("/wallet/active", controllers.FindActiveWallets)
	r.POST("/wallet", controllers.UploadWallet)
	r.POST("/walletform", controllers.UploadWalletFile)
	r.GET("/wallets/:address", controllers.FindWallet)
	r.GET("/wallet/idle", controllers.FindIdleWallet)
	r.GET("/remove/:address", controllers.DeleteWallet)
	//r.GET("/generateid/:address", controllers.GenerateID)

	r.Run("0.0.0.0:30050")
}
