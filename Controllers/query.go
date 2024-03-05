package controllers

import (
	models "Assignment_4/models"
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
)

func ShowBanks(db *pg.DB) ([]*models.Bank, error) {

	var banks []*models.Bank

	if err := db.Model(&banks).
		Relation("Branch").
		Order("id ASC").
		Select(); err != nil {
		return banks, err
	}
	return banks, nil
}

func AddNewBank(db *pg.DB, newBank *models.Bank) (uint, error) {

	if _, err := db.Model(newBank).Insert(); err != nil {
		return 0, err
	}
	return newBank.ID, nil
}

func GetBankById(db *pg.DB, bankId uint) (models.Bank, error) {
	var bank models.Bank

	if selectErr := db.Model(&bank).
		Relation("Branch").
		Where("bank.Id = ?", bankId).
		Select(); selectErr != nil {
		return bank, selectErr
	}

	return bank, nil
}

func UpdateBank(db *pg.DB, bank *models.Bank, bankID uint) (uint, error) {

	if _, updateErr := db.Model(bank).Where("id = ?", bankID).UpdateNotZero(); updateErr != nil {
		return bankID, updateErr
	}

	return bankID, nil
}

func DeleteBank(tx *pg.Tx, bankId uint) error {

	if _, err := tx.Model((*models.Bank)(nil)).Where("id = ?", bankId).Delete(); err != nil {
		return err
	}
	return nil
}

func GetBranches(db *pg.DB) ([]*models.Branch, error) {

	var branches []*models.Branch
	if err := db.Model(&branches).Relation("Bank").Select(); err != nil {
		return branches, err
	}
	return branches, nil
}

func GetBranchById(db *pg.DB, branchId uint) (*models.Branch, error) {
	branch := &models.Branch{}
	if err := db.Model(branch).Relation("Bank").Where("branch.id=?", branchId).Select(); err != nil {
		return branch, err
	}
	return branch, nil
}

func UpdateBranchData(tx *pg.Tx, branchId uint, branch *models.Branch) (uint, error) {

	res, err := tx.Model(branch).Where("id = ?", branchId).UpdateNotZero(branch)

	if err != nil {
		tx.Rollback()
		return branchId, err
	}
	if res.RowsAffected() == 0 {
		tx.Rollback()
		return branchId, errors.New("no record updated")
	}
	return branchId, nil
}

func DelBranchById(tx *pg.Tx, branchId uint) error {
	res, err := tx.Model((*models.Branch)(nil)).Where("id = ?", branchId).Delete()

	if err != nil {
		tx.Rollback()
		return err
	}
	if res.RowsAffected() == 0 {
		tx.Rollback()
		return errors.New("no record deleted")
	}
	return nil
}

func GetCustDetails(db *pg.DB, custId uint) (models.Customer, error) {

	var customer models.Customer
	if err := db.Model(&customer).Where("cust_id = ?", custId).Relation("Account").Select(); err != nil {
		return models.Customer{}, err
	}
	return customer, nil

}

func updateCustDetails(tx *pg.Tx, customer models.Customer) error {

	res, updateErr := tx.Model(&customer).WherePK().UpdateNotZero()

	if updateErr != nil {
		tx.Rollback()
		return updateErr
	}

	if res.RowsAffected() == 0 {
		tx.Rollback()
		return errors.New("no record updated")
	}
	return nil

}
func ViewTransDetails(db *pg.DB, transId uint) (models.Transaction, error) {

	var fetchedTrans models.Transaction

	if err := db.Model(&fetchedTrans).Relation("Account").Where("id=?", transId).Select(); err != nil {
		return models.Transaction{}, err
	}
	return fetchedTrans, nil
}

func SearchTransaction(db *pg.DB, startDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction

	if err := db.Model(&transactions).Where("tr_date >= ? AND tr_date < ?", startDate, startDate.AddDate(0, 0, 1)).Select(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func RetrieveAccDetails(db *pg.DB, accId uint) (models.Account, error) {

	var account models.Account

	err := db.Model(&account).Relation("Branch.Bank").Relation("Customers").Where("acc_id = ?", accId).Select()
	if err != nil {
		return models.Account{}, err
	}
	return account, nil
}

func GetAccTransDetails(db *pg.DB, accId uint) (models.Account, error) {
	var account models.Account
	if err := db.Model(&account).Relation("Transaction").Where("acc_id=?", accId).Select(); err != nil {
		return models.Account{}, err
	}
	return account, nil
}

func SaveCustomers(tx *pg.Tx, inputCredentials showCustAccdetail) ([]uint, error) {

	var custIds []uint
	for _, customer := range inputCredentials.Customers {
		//check if the customer details are already present in customers table, if yes then only retrieve
		existingCustomer := &models.Customer{}
		err := tx.Model(existingCustomer).Where("pan = ?", customer.PAN).Where("branch_id = ?", inputCredentials.BranchID).Select()

		if err == nil {
			custIds = append(custIds, existingCustomer.ID)
			//customer already exists
		} else if err == pg.ErrNoRows {

			cust := &models.Customer{
				Name:     customer.Name,
				PAN_ID:   customer.PAN,
				DOB:      customer.DOB,
				Address:  customer.Address,
				BranchID: inputCredentials.BranchID,
			}

			_, err := tx.Model(cust).Insert()

			if err != nil {
				tx.Rollback()
				return []uint{}, err
			}
			//this slice will be used for mapping accounts with customer
			custIds = append(custIds, cust.ID)
		} else {
			tx.Rollback()
			return []uint{}, errors.New("failed to check existing customer")
		}
	}
	return custIds, nil
}

func SaveAccount(tx *pg.Tx, account *models.Account) (uint, error) {

	if _, accErr := tx.Model(account).Insert(); accErr != nil {
		tx.Rollback()
		return 0, errors.New("failed to insert record in accounts table")
	}
	return account.Acc_ID, nil
}

func SaveCustomerAccount(tx *pg.Tx, custIds []uint, accId uint) error {

	for _, custId := range custIds {

		_, err := tx.Model(&models.CustomerToAccount{
			CustomerID: accId,
			AccountID:  custId,
		}).Insert()

		if err != nil {
			tx.Rollback()
			return errors.New("failed to insert record in accounts table")
		}
	}
	return nil
}

func GetCustomerAcc(tx *pg.Tx, accId uint) ([]models.CustomerToAccount, error) {

	var accCust []models.CustomerToAccount
	if err := tx.Model(&accCust).Where("acc_id=?", accId).Select(); err != nil {

		tx.Rollback()
		return []models.CustomerToAccount{}, err
	}
	return accCust, nil
}

func DeleteCustomer(tx *pg.Tx, accCust []models.CustomerToAccount) error {

	for _, customer := range accCust {

		count, err := tx.Model((*models.CustomerToAccount)(nil)).Where("cust_id = ?", customer.CustomerID).Count()
		if err != nil {
			tx.Rollback()
			return errors.New("failed to retrieve count of customer")
		}
		if count == 1 {
			//delete record
			if _, err := tx.Model((*models.Customer)(nil)).Where("cust_id=?", customer.CustomerID).Delete(); err != nil {
				tx.Rollback()
				return errors.New("failed to delete customer data from customers table")
			}

		}
	}

	return nil
}

func DeleteAccount(tx *pg.Tx, accId uint) error {
	if _, err := tx.Model((*models.Account)(nil)).Where("acc_id=?", accId).Delete(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func InsertTransaction(tx *pg.Tx, transaction models.Transaction) error {

	if _, err := tx.Model(&transaction).Insert(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func UpdateAccount(tx *pg.Tx, account models.Account) error {

	if _, err := tx.Model(&account).Where("acc_id = ?", account.Acc_ID).Update(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
