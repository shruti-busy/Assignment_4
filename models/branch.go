package models

type Branch struct {
	ID        uint `pg:",pk"`
	Name      string
	Address   string
	IFSCCode  string
	BankID    uint        `pg:",fk:bank_id,on_delete:SET NULL"`
	Bank      Bank        `pg:"rel:has-one"`
	Customers []*Customer `pg:"rel:has-many,on_delete:CASCADE"`
}