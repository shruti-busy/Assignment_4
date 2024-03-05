package Routes

import (
	controllers "Assignment_4/controllers"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(connectingRoutes *gin.Engine) {

	connectingRoutes.GET("/customers/:cust_id", controllers.GetCustomerDetails())

	connectingRoutes.PUT("/customers/update", controllers.UpdateCustomeDetails())
}
