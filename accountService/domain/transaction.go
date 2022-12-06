package domain

import "time"

type Transaction struct {
	ID                int       `json:"id" sql:"type:integer;not null"`
	AccountSenderId   int       `json:"account_sender_id" sql:"type:integer;not null"`
	AccountReceiverId int       `json:"account_receiver_id" sql:"type:integer;not null"`
	Amount            int       `json:"amount" sql:"type:integer;default 0;not null"`
	IsDeleted         bool      `json:"is_deleted" sql:"type:boolean;default false;not null"`
	CreationDate      time.Time `json:"creation_date" sql:"type:timestamp(6);default CURRENT_TIMESTAMP(6);not null"`
	ModificationDate  time.Time `json:"modification_date" sql:"type:timestamp(6);default CURRENT_TIMESTAMP(6);not null"`
}
