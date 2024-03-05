package controllers

import (
	database "Assignment_4/database"
	models "Assignment_4/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type showCustAccdetail struct {
	Customers []struct {
		Name    string    `json:"name"`
		DOB     time.Time `json:"dob"`
		Phone   int       `json:"phone"`
		Address string    `json:"address"`
		PAN     string    `json:"pan"`
	} `json:"customers"`
	BranchID uint    `json:"branch_id"`
	Balance  float64 `json:"balance"`
	AccType  string  `json:"acc_type"`
}

type WithdrawOrDeposit struct {
	AccID  uint    `json:"acc_id"`
	Amount float64 `json:"amount"`
	Mode   string  `json:"mode"`
}

func ShowAccount() gin.HandlerFunc {
	return func(c *gin.Context) {

		var inputCredentials showCustAccdetail
		err := c.ShouldBindJSON(&inputCredentials)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Failed to parse request payload",
			})
			return
		}
		db := database.ReturnDBIns()

		//start a transaction
		tx, err := db.Begin()
		if err != nil {
			tx.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		var custIds []uint

		//insert record in customer table
		// custIds []uint will be used for mapping
		custIds, custErr := SaveCustomers(tx, inputCredentials)
		if custErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": custErr.Error(),
			})
			return
		}

		//an instance of Account struct
		acc := &models.Account{
			Balance:  inputCredentials.Balance,
			AccType:  inputCredentials.AccType,
			BranchID: inputCredentials.BranchID,
		}
		//insert record in accounts table
		accId, accErr := SaveAccount(tx, acc)
		if accErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": accErr.Error(),
			})
			return
		}

		//establish relation between accounts and customers table
		if mapErr := SaveCustomerAccount(tx, custIds, accId); mapErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   mapErr.Error(),
				"message": "Failed to insert record in accounts table",
			})
			return
		}

		//return the response
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}
		//return success response
		c.JSON(http.StatusCreated, gin.H{
			"message": "Account opened successfully",
		})
	}
}

func CloseAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve acc_id
		accIdStr := c.Param("acc_id")
		accID, err := strconv.ParseUint(accIdStr, 10, strconv.IntSize)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		db := database.ReturnDBIns()

		tx, txErr := db.Begin()
		if txErr != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to start transaction",
			})
			return
		}

		//retrieve corresponding customer details
		// accCust is a []models.CustomerToAccount
		accCust, selErr := GetCustomerAcc(tx, uint(accID))

		if selErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve account data from mapping table",
			})
			return
		}

		//iterating over this slice to delete corresponding customer details
		if delErr := DeleteCustomer(tx, accCust); delErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": delErr.Error(),
			})
			return
		}

		if delErr := DeleteAccount(tx, uint(accID)); delErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to delete account details from account table",
				"error":   delErr.Error(),
			})
			return
		}

		//commit the transaction
		tx.Commit()
		//return response of success
		c.JSON(http.StatusOK, gin.H{
			"message": "Account closed successfully",
		})

	}
}

func GetAccountDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve acc_id from url
		accIdStr := c.Param("acc_id")
		accId, err := strconv.ParseUint(accIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Retrieve the account details
		db := database.ReturnDBIns()

		if account, selectErr := GetAccTransDetails(db, uint(accId)); selectErr == nil { /////

			// Return the account details
			c.JSON(http.StatusOK, gin.H{
				"message": "Account details retrieved successfully",
				"account": account,
			})
		} else {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve account details",
			})
		}
	}
}

func WithdrawFunds() gin.HandlerFunc {
	return func(c *gin.Context) {
		//to store request payload
		var withdrawData WithdrawOrDeposit

		if err := c.ShouldBindJSON(&withdrawData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//get db instance
		db := database.ReturnDBIns()

		//retrieve the account record from accounts table
		account, err := GetAccTransDetails(db, withdrawData.AccID) ///
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		//open a transaction
		tx, txErr := db.Begin()

		if txErr != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to start transaction",
			})
			return
		}

		account.Balance -= withdrawData.Amount

		//record this in transaction table
		// Create a new transaction record for deposit
		transaction := models.Transaction{
			Mode:              withdrawData.Mode,
			ReceiverAccNumber: account.AccNo,
			Amount:            withdrawData.Amount,
			AccountID:         account.Acc_ID,
		}
		//add record in transaction table
		if insertErr := InsertTransaction(tx, transaction); insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
			return
		}

		// Update the account record with the new balance
		if err := UpdateAccount(tx, account); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance"})
			return
		}
		// commit the transaction
		tx.Commit()

		//return appropriate response
		c.JSON(http.StatusOK, gin.H{
			"message":     "Amount withdrawn successfully",
			"new Balance": account.Balance,
		})

	}
}

func DepositFunds() gin.HandlerFunc {
	return func(c *gin.Context) {

		//to store request payload
		var depositData WithdrawOrDeposit

		if err := c.ShouldBindJSON(&depositData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//get db instance
		db := database.ReturnDBIns()

		//retrieve the account record from accounts table
		account, err := GetAccTransDetails(db, depositData.AccID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		//open a transaction
		tx, txErr := db.Begin()

		if txErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to start transaction",
			})
			return
		}

		// Add the deposit amount to the current balance
		account.Balance += depositData.Amount

		//record this in transaction table
		// Create a new transaction record for the deposit
		transaction := models.Transaction{
			Mode:              depositData.Mode,
			ReceiverAccNumber: account.AccNo,
			Amount:            depositData.Amount,
			AccountID:         account.Acc_ID,
		}
		//add record in transaction table
		if insertErr := InsertTransaction(tx, transaction); insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
			return
		}

		// Update the account record with the new balance
		if err := UpdateAccount(tx, account); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance"})
			return
		}

		// commit the transaction
		tx.Commit()

		//return appropriate response
		c.JSON(http.StatusOK, gin.H{
			"message":     "Amount deposited successfully",
			"new Balance": account.Balance,
		})

	}
}

func GetTransactionDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve acc_id from request url
		accIdStr := c.Param("acc_id")
		accId, err := strconv.ParseUint(accIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to retrieve account id from request url",
			})
			return
		}

		//retrieve all the transactions related to that account
		db := database.ReturnDBIns()

		if account, selErr := GetAccTransDetails(db, uint(accId)); selErr == nil {
			//return the appropriate response
			c.JSON(http.StatusOK, gin.H{
				"message":             "Transaction details retrieved successfully",
				"transaction details": account,
			})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to retrieve transaction details related to thi account",
			})

		}

	}
}
