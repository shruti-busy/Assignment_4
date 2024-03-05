package Routes

import (
	controllers "Assignment_4/controllers"

	"github.com/gin-gonic/gin"
)

func BranchRoutes(connectingRoutes *gin.Engine) {

	connectingRoutes.GET("/branches/", controllers.ShowAllBranches())
	connectingRoutes.GET("/branches/:branch_id", controllers.GetBranchDetails())
	connectingRoutes.POST("/branches", controllers.CreateBranch())
	connectingRoutes.PUT("/branches/:branch_id", controllers.UpdateBranch())
	connectingRoutes.DELETE("/branches/:branch_id", controllers.DeleteBranch())
}
