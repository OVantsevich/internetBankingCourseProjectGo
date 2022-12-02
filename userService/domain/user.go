package domain

import (
	"time"
)

type User struct {
	ID               int       `json:"id" sql:"type:integer;not null"`
	UserLogin        string    `json:"user_login" sql:"type:varchar(50);not null"`
	UserEmail        string    `json:"user_email" sql:"type:varchar(50);not null"`
	UserPassword     string    `json:"user_password" sql:"type:varchar(50);not null"`
	UserName         string    `json:"user_name" sql:"type:varchar(50);not null"`
	Surname          string    `json:"surname" sql:"type:varchar(50);not null"`
	IsDeleted        bool      `json:"is_deleted" sql:"type:boolean;default false;not null"`
	CreationDate     time.Time `json:"creation_date" sql:"type:timestamp(6);default CURRENT_TIMESTAMP(6);not null"`
	ModificationDate time.Time `json:"modification_date" sql:"type:timestamp(6);default CURRENT_TIMESTAMP(6);not null"`
}

func (old *User) UpdateUser(new *User) {
	if new.UserName != "" && old.UserName != new.UserName {
		old.UserName = new.UserName
	}
	if new.Surname != "" && old.Surname != new.Surname {
		old.Surname = new.Surname
	}
	if new.UserEmail != "" && old.UserEmail != new.UserEmail {
		old.UserEmail = new.UserEmail
	}
	if new.UserPassword != "" && old.UserPassword != new.UserPassword {
		old.UserPassword = new.UserPassword
	}
}
