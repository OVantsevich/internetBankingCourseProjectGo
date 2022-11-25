package domain

import (
	"time"
)

type Account struct {
	ID               int       `json:"id" sql:"type:integer;not null"`
	UserId           int       `json:"user_id" sql:"type:integer;not null"`
	AccountName      string    `json:"account_name" sql:"type:varchar(40);default 'account'::character varying;not null"`
	Amount           int       `json:"amount" sql:"type:integer;default 0;not null"`
	CreationDate     time.Time `json:"creation_date" sql:"type:timestamp(6);default CURRENT_TIMESTAMP(6);not null"`
	ModificationDate time.Time `json:"modification_date" sql:"type:timestamp(6);default CURRENT_TIMESTAMP(6);not null"`
	IsDeleted        bool      `json:"is_deleted" sql:"type:boolean;default false;not null"`
}
