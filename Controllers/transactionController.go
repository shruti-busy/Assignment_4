package controllers

import (
	database "Assignment_4/database"
	models "Assignment_4/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ShowTransactionDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		transIdString := c.Param("id") //get the transaction id from request url
		transactionId, err := strconv.ParseUint(transIdString, 10, strconv.IntSize)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to retrieve transaction id from the url",
			})
			return
		}

		db := database.ReturnDBIns()

		if fetchedTransaction, selErr := ViewTransDetails(db, uint(transactionId)); selErr == nil {
			c.JSON(http.StatusOK, gin.H{
				"message":             "Details retrieved successfully",
				"transaction details": fetchedTransaction,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve  transaction details",
			})
		}

	}
}

func TransferMoney() gin.HandlerFunc {
	return func(c *gin.Context) {

		transferDetails := models.Transaction{}
		//retrieve details from request payload
		if err := c.ShouldBindJSON(&transferDetails); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		db := database.ReturnDBIns()

		//get sender and reciever accounts
		sender, err := GetAccTransDetails(db, transferDetails.AccountID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve sender account details",
				"error":   err.Error(),
			})
			return
		}

		var receiver models.Account
		if err := db.Model(&receiver).Where("acc_no=?", transferDetails.ReceiverAccNumber).Select(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve reciever's account details",
			})
			return
		}
		//begin transaction
		tx, txErr := db.Begin()

		if txErr != nil {
			tx.Close()
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to start transaction",
			})
			return
		}

		//insert this in transactions table
		if err := InsertTransaction(tx, transferDetails); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to insert record in transactions table",
				"error":   err.Error(),
			})
			return
		}

		sender.Balance -= transferDetails.Amount
		receiver.Balance += transferDetails.Amount
		//update the record in accounts table for both sender and reciever
		if err := UpdateAccount(tx, sender); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sender's account balance"})
			return
		}

		if err := UpdateAccount(tx, receiver); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receiver's account balance"})
			return
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}
		//return appropriate response
		c.JSON(http.StatusOK, gin.H{
			"message": "Transfer completed successfully",
		})

	}
}

func SearchTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		//query param
		startDateString := c.Query("start_date")

		startDate, err := time.Parse(time.RFC3339, startDateString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to retrieve the date",
			})
			return
		}

		db := database.ReturnDBIns()
		//retrieve the transactions
		if transactions, err := SearchTransaction(db, startDate); err == nil {
			//return the response
			c.JSON(http.StatusOK, gin.H{
				"message":      "Transactions retrieved",
				"Transactions": transactions,
			})
		} else {

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve the transactions",
			})

		}
	}
}
