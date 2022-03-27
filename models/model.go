package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type UserWallet struct {
	gorm.Model
	Username string  ` json:"username"`
	Wallet   float64 `json:"wallet"`
}

type CreditAmount struct {
	Username string  `json:"username"`
	Amount   float64 `json:"amount"`
}

type DebitAmount struct {
	Amount float64 `json:"amount"`
}
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserWalletReturn struct {
	Username string  ` json:"username"`
	Wallet   float64 `json:"wallet"`
}
