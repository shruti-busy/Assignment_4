package models

import "time"

type Transaction struct {
	ID                uint `pg:",pk"`
	AccountID         uint `pg:",fk:accounts,on_delete:SET NULL"`
	Mode              string
	ReceiverAccNumber string
	Timestamp         time.Time `pg:",default:now()"`
	Amount            float64
	Account           *Account `pg:"rel:has-one"`
}
