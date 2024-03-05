package Routes

import (
	controllers "Assignment_4/controllers"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(connectingRoutes *gin.Engine) {

	connectingRoutes.GET("/transaction/:id", controllers.ShowTransactionDetails())
	connectingRoutes.POST("/transaction/transfer", controllers.TransferMoney())
	connectingRoutes.GET("/transaction/search", controllers.SearchTransactions())

}
