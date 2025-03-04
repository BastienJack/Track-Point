package db

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(40);not null;index:idx_username,unique;" json:"user_name,omitempty"`
	Password string `gorm:"type:varchar(256);not null;" json:"password,omitempty"`
}

func (User) TableName() string {
	return "Users"
}

func GetUserByName(ctx context.Context, username string) (*User, error) {
	res := new(User)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).
		Select("id, username, password").Where("username = ?", username).
		First(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func GetUserById(ctx context.Context, id int64) (*User, error) {
	res := new(User)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).
		First(&res, id).Error; err != nil {
		return res, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

func CreateUser(ctx context.Context, user *User) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(db *gorm.DB) error {
		if err := db.Create(user).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
