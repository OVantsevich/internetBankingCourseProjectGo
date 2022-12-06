package domain

type User struct {
	ID        int    `json:"id" sql:"type:integer;not null"`
	UserLogin string `json:"user_login" sql:"type:varchar(50);not null"`
	UserName  string `json:"user_name" sql:"type:varchar(50);not null"`
	Surname   string `json:"surname" sql:"type:varchar(50);not null"`
	IsDeleted bool   `json:"is_deleted" sql:"type:boolean;default false;not null"`
}

func (old *User) UpdateUser(new *User) {
	if new.UserName != "" && old.UserName != new.UserName {
		old.UserName = new.UserName
	}
	if new.Surname != "" && old.Surname != new.Surname {
		old.Surname = new.Surname
	}
}
