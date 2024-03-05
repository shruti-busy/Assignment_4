package database

import (
	"Assignment_4/models"
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func Connect() *pg.DB {
	opts := &pg.Options{
		User:     "postgres",
		Password: "admin",
		Addr:     "localhost:5432",
		Database: "postgres",
	}

	db := pg.Connect(opts)
	if db == nil {
		fmt.Print("Connection failed")
		os.Exit(100)
	} else {

		fmt.Println("Successfully connected to the database")
	}

	err := CreateBankTables(db)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func CreateBankTables(db *pg.DB) error {
	models := []interface{}{
		(*models.Bank)(nil),
		(*models.Branch)(nil),
		(*models.Customer)(nil),
		(*models.Account)(nil),
		(*models.CustomerToAccount)(nil),
		(*models.Transaction)(nil),
	}

	for _, model := range models {
		createTableError := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if createTableError != nil {
			log.Printf("Error while creating the table,Reason:%v", createTableError)
			return createTableError
		}

	}

	return nil

}

var db *pg.DB

func ReturnDBIns() *pg.DB {
	return db
}
