package controllers

import (
	database "Assignment_4/database"
	models "Assignment_4/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBranchDetails() gin.HandlerFunc {
	return func(c *gin.Context) {

		branchIdString := c.Param("branch_id")
		if branchIdString == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please enter branch ID",
			})
			return
		}

		branchId, err := strconv.ParseUint(branchIdString, 10, strconv.IntSize)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		db := database.ReturnDBIns()

		if branch, selectErr := GetBranchById(db, uint(branchId)); selectErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": selectErr.Error(),
			})

		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":"Branch details retrieved",
				"branch details": branch,
			})
		}

	}
}

func ShowAllBranches() gin.HandlerFunc {
	return func(c *gin.Context) {

		var branches []*models.Branch
		db := database.ReturnDBIns()

		branches, selectErr := GetBranches(db)

		if selectErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": selectErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "Branch details retrieved",
			"branch details": branches,
		})
	}
}

func CreateBranch() gin.HandlerFunc {
	return func(c *gin.Context) {

		var branch models.Branch

		if err := c.ShouldBindJSON(&branch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		//check if the table already exits or not in the table
		db := database.ReturnDBIns()
		count, err := db.Model((*models.Bank)(nil)).Where("id=?", branch.BankID).Count()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while checking bank existence",
			})
			return
		}
		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "bank does not exist",
			})
			return
		}

		//begin a transaction
		tx, txErr := db.Begin()

		// Make sure to close it if something goes wrong.

		if txErr != nil {
			tx.Close()
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error in establishing transaction",
				"error":   txErr.Error(),
			})
			return
		}

		defer tx.Rollback() // Rollback transaction incase of any error

		//add new record in branch table.
		_, createErr := tx.Model(&branch).Insert()

		if createErr != nil {
			tx.Rollback()

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while creating the branch record",
				"error":   createErr.Error(),
			})

			return
		}

		//coomit the transaction and return response
		tx.Commit()

		c.JSON(http.StatusOK, gin.H{
			"message":         "Branch created successfully",
			"inserted record": branch,
		})

	}
}

func UpdateBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		branchIdStr := c.Param("branch_id")   //get id from url
		branchID, err := strconv.ParseUint(branchIdStr, 10, strconv.IntSize)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		//retrieve the data from json which is to be updated
		branch := &models.Branch{}

		if err := c.ShouldBindJSON(branch); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		db := database.ReturnDBIns()

		//begin transaction
		tx, txErr := db.Begin()
		if txErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to start transaction",
			})
		}

		updatedBranchId, err := UpdateBranchData(tx, uint(branchID), branch)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//commit the transaction
		tx.Commit()
		c.JSON(http.StatusCreated, gin.H{
			"message":   "Branch has been updated successfully",
			"Branch ID": updatedBranchId,
		})

	}
}

func DeleteBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//retrieve the branchId from url
		branchIdString := c.Param("branch_id")
		branchId, parseErr := strconv.ParseUint(branchIdString, 10, strconv.IntSize)

		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": parseErr.Error(),
			})
			return
		}

		db := database.ReturnDBIns()
		
  //start the transaction 
		tx, txErr := db.Begin()
		if txErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": txErr.Error(),
			})
			return
		}

		// Delete the bank record
		deleteErr := DelBranchById(tx, uint(branchId))

		if deleteErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": deleteErr.Error()})
			return
		}

		//commit all changes and return response
		tx.Commit()

		c.JSON(http.StatusOK, gin.H{
			"message": "Branch deleted successfully",
		})
	}
}
