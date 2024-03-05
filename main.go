package main

import (
	Routes "Assignment_4/Routes"
	database "Assignment_4/database"
	models "Assignment_4/models"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func init() {
	orm.RegisterTable((*models.CustomerToAccount)(nil))
}

var pgdb *pg.DB

func main() {

	pgdb = database.Connect()
	createErr := database.CreateBankTables(pgdb)
	if createErr != nil {
		panic(createErr)
	}

	router := gin.Default()

	Routes.AccountRoutes(router)
	Routes.BankRoutes(router)
	Routes.BranchRoutes(router)
	Routes.CustomerRoutes(router)
	Routes.TransactionRoutes(router)
	router.Run()
}
