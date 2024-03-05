package Routes

import (
	controllers "Assignment_4/controllers"

	"github.com/gin-gonic/gin"
)

func AccountRoutes(connectingRoutes *gin.Engine) {

	connectingRoutes.POST("/accounts/view", controllers.ShowAccount())
	connectingRoutes.POST("/accounts/close/:acc_id", controllers.CloseAccount())
	connectingRoutes.GET("/accounts/:acc_id", controllers.GetAccountDetails())
	connectingRoutes.POST("/accounts/deposit", controllers.DepositFunds())
	connectingRoutes.POST("/accounts/withdraw", controllers.WithdrawFunds())
	connectingRoutes.GET("/accounts/:acc_id/transactionhistory", controllers.GetTransactionDetails())

}
