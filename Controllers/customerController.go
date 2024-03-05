package controllers

import (
	database "Assignment_4/database"
	models "Assignment_4/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCustomerDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		custIdString := c.Param("cust_id") //retrieve the customer id from url
		custId, err := strconv.ParseUint(custIdString, 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to retrieve the customer id",
			})
			return
		}
		//retrieve the customer details
		db := database.ReturnDBIns()

		customer, err := GetCustDetails(db, uint(custId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve the customer details",
			})
			return
		}

		//return response
		c.JSON(http.StatusOK, gin.H{
			"message": "Details retrieved",
			"details": customer,
		})

	}
}

func UpdateCustomeDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
      //retrieve the details from request payload
		var cust models.Customer
		err := c.ShouldBindJSON(&cust)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to retrieve details to update",
			})
			return
		}

		db := database.ReturnDBIns()   //update the details

		tx, txErr := db.Begin()    //Start transaction
		if txErr != nil {
			tx.Close()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to begin transaction",
			})
			return
		}
		 
		//update customer details
		Up_err := updateCustDetails(tx, cust)  
		
		if Up_err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update details",
			})
			return
		}
		tx.Commit()

		c.JSON(http.StatusOK, gin.H{
			"message": "Details updated successfully",
		})

	}
}
