package domain

import "time"

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

func (dest *User) CompareAndSet(src *User) {
	if src.UserName != "" && dest.UserName != src.UserName {
		dest.UserName = src.UserName
	}
	if src.Surname != "" && dest.Surname != src.Surname {
		dest.Surname = src.Surname
	}
	if src.UserLogin != "" && dest.UserLogin != src.UserLogin {
		dest.UserLogin = src.UserLogin
	}
	if src.UserEmail != "" && dest.UserEmail != src.UserEmail {
		dest.UserEmail = src.UserEmail
	}
	if src.UserPassword != "" && dest.UserPassword != src.UserPassword {
		dest.UserPassword = src.UserPassword
	}
}
