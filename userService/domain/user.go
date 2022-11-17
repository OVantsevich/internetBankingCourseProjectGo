package domain

import "time"

type User struct {
	ID               string    `json:"id" sql:"type:integer;not null"`
	UserName         string    `json:"user_name" sql:"type:varchar(20);not null"`
	Surname          string    `json:"surname" sql:"type:varchar(50);not null"`
	IsDeleted        bool      `json:"is_deleted" sql:"type:boolean;not null"`
	CreationDate     time.Time `json:"creation_date" sql:"type:timestamp(6);default CURRENT_TIMESTAMP(6);not null"`
	ModificationDate time.Time `json:"modification_date" sql:"type:timestamp(6);default CURRENT_TIMESTAMP(6);not null"`
	UserLogin        string    `json:"user_login" sql:"type:varchar(100);not null"`
	UserPassword     string    `json:"user_password" sql:"type:varchar(200);not null"`
}

func (dest *User) CompareAndSet(src *User) {

	if dest.UserName != src.UserName {
		dest.UserName = src.UserName
	}
	if dest.Surname != src.Surname {
		dest.Surname = src.Surname
	}
	if dest.IsDeleted != src.IsDeleted {
		dest.IsDeleted = src.IsDeleted
	}
	if dest.UserLogin != src.UserLogin {
		dest.UserLogin = src.UserLogin
	}
	if dest.UserPassword != src.UserPassword {
		dest.UserPassword = src.UserPassword
	}
}
