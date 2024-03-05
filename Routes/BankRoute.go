package Routes

import (
	controllers "Assignment_4/controllers"

	"github.com/gin-gonic/gin"
)

func BankRoutes(connectingRoutes *gin.Engine) {

	connectingRoutes.GET("/bank", controllers.GetAllBanks())
	connectingRoutes.GET("/bank/:bank_id", controllers.GetBankDetails())
	connectingRoutes.PUT("/bank/:bank_id", controllers.UpdateBankDetails())
	connectingRoutes.POST("/bank", controllers.CreateBank())
	connectingRoutes.DELETE("/bank/:bank_id", controllers.DeleteBankfunc())
}
