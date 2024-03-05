package models

type Bank struct {
	ID       uint `pg:",pk"`
	BankName string
	Branches []*Branch `pg:"rel:has-many,on_delete:CASCADE"`
}
