package models

import "time"

type Customer struct {
	ID       uint `pg:",pk"`
	Name     string
	Age      int
	DOB      time.Time
	BranchID uint       `json:"branchId" pg:",on_delete:CASCADE"`
	PAN_ID   string     `pg:",notnull" json:"pan"`
	Phone    string     `pg:",notnull" json:"phone"`
	Address  string     `pg:",notnull" json:"address"`
	Branch   *Branch    `pg:"rel:has-one"`
	Account  []*Account `pg:"many2many:customer_to_accounts,on_delete:CASCADE"`
}
