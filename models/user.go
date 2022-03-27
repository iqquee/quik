package models

import (
	"quik/config"
)

//this function created new user into the database
func CreateUser(user *User) (err error) {
	//chcek if the username is already taken

	if err := config.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByUsername(user *User, username string) (*User, error) {
	if err := config.DB.Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserWalletById(user *UserWallet, id string) (*UserWallet, error) {
	if err := config.DB.Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
