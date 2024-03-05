package models

type CustomerToAccount struct {
	ID	uint
	CustomerID uint      `pg:"on_delete:CASCADE"`
	Customer   *Customer `pg:"rel:has-one"`
	AccountID  uint      `pg:"on_delete:CASCADE"`
	Account    *Account  `pg:"rel:has-one"`
}