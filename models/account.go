package models

type Account struct {
	Acc_ID           uint `pg:",pk"`
	BranchID     uint `pg:",fk:branches,on_delete:SET NULL"`
	AccNo        string
	Balance      float64
	AccType      string
	Branch       *Branch        `pg:"rel:has-one"`
	Transactions []*Transaction `pg:"rel:has-many,on_delete:CASCADE"`
	Customer     []*Customer    `pg:"many2many:customer_to_accounts"`
}
