package models

import (
	"fmt"
	"quik/config"
)

func CreateWallet(wallet *UserWallet) (err error) {
	if err := config.DB.Create(wallet).Error; err != nil {
		return err
	}
	return nil
}

func UpdateWalletFund(wallet *UserWallet, username string, amount string) (err error) {
	if err := config.DB.Exec(fmt.Sprintf("UPDATE user_wallets SET wallet = %v WHERE username = '%s'", amount, username)).Error; err != nil {
		return err
	}
	return nil
}

func GetAllWallets(wallet *[]UserWallet) (err error) {
	if err := config.DB.Find(wallet).Error; err != nil {
		return err
	}
	return nil
}
