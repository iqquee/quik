package controllers

import (
	"fmt"
	"math"
	"net/http"
	"quik/models"

	"quik/config"

	"github.com/gin-gonic/gin"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
)

var (
	paramsIdErr = "the wallet id needs to be passed in the request parameter"
	headerErr   = "the username is needs to be passed in the request header"
)

func GetWalletBalance(c *gin.Context) {
	var user models.User
	var wallet models.UserWallet

	//logs the incoming request
	log.Info(c.Request)

	userName := c.Request.Header.Get("username")
	if userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": headerErr,
		})
		return
	}

	//query the redis database to find a balance
	//if there was no balance found then get it from the mysql database
	//and thereafter save the balance into redis
	val := config.Cached.Get(userName)
	log.Println(val)
	if val == nil {
		userResult, err := models.GetUserByUsername(&user, userName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("the user with username %s does not exist", userName),
			})
			log.Println(err.Error())
			return
		}

		params := c.Param("wallet_id")
		if params == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": paramsIdErr,
			})
			return
		}
		walletResult, err := models.GetUserWalletById(&wallet, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("the wallet with id of %s does not exist", params),
			})
			log.Println(err.Error())
			return
		}
		balance := decimal.NewFromFloat(walletResult.Wallet)

		// check if user is the owner of the wallet by comparing the username
		if userResult.Username == walletResult.Username {
			// save the user balance into redis
			config.Cached.Set(walletResult.Username, walletResult.Wallet, 0).Err()
			c.JSON(http.StatusOK, gin.H{
				"balance": balance,
			})
		} else {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "sorry, you are not the owner of the wallet you are trying to access",
			})
			return
		}
	} else {
		data, _ := val.Int()

		c.JSON(http.StatusOK, gin.H{
			"balance": data,
		})
	}

}

func CreditWallet(c *gin.Context) {
	var foundUser models.User
	var wallet models.UserWallet
	var creditDetails models.CreditAmount

	if err := c.BindJSON(&creditDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//logs the incoming request
	log.Info(c.Request)

	//check the users input if it contains a negative sign before even
	//querying the database to retrieve data to minimize resources and response time
	if err := math.Signbit(creditDetails.Amount); err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you cannot credit a users wallet with a negative value",
		})
		return
	}
	//check if the amout to credit with is not 0
	if creditDetails.Amount <= 0 || creditDetails.Amount <= 0.0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "your credit amount cannot be 0 or less than 0",
		})
		return
	}
	userResult, err := models.GetUserByUsername(&foundUser, creditDetails.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the user with username %s does not exist", creditDetails.Username),
		})
		log.Println(err.Error())
		return
	}

	params := c.Param("wallet_id")
	if params == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": paramsIdErr,
		})
		return
	}
	walletResult, err := models.GetUserWalletById(&wallet, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("the wallet with id of %s does not exist", params),
		})
		log.Println(err.Error())
		return
	}

	//we need to be sure that you are crediting the right wallet
	if userResult.Username == walletResult.Username {

		previousBalance := decimal.NewFromFloat(walletResult.Wallet)
		credit := decimal.NewFromFloat(creditDetails.Amount)
		newBalance := previousBalance.Add(credit).String()
		cacheErr := config.Cached.Set(walletResult.Username, newBalance, 0).Err()
		if cacheErr != nil {
			log.Println(cacheErr)
			return
		}

		//check if the fund to be sent contains a minus sign before authorizing the credit
		err := models.UpdateWalletFund(&wallet, creditDetails.Username, newBalance) //the value in database has to be a float e.g 0.0
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "an error occured while updating your wallet",
			})
			log.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"balance": newBalance,
			"message": fmt.Sprintf("%s your wallet has been credited with %v", userResult.First_Name, creditDetails.Amount),
		})
	} else {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("the user with username %s is not the owner of th wallet with an id of %d", userResult.Username, walletResult.ID),
		})
		return
	}
}

func DebitWallet(c *gin.Context) {
	var foundUser models.User
	var wallet models.UserWallet
	var debitDetails models.DebitAmount

	if err := c.BindJSON(&debitDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//logs the incoming request
	log.Info(c.Request)

	//check the users input if it contains a negative sign before even
	//querying the database to retrieve data to minimize resources and response time
	if err := math.Signbit(debitDetails.Amount); err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you cannot debit your wallet with a negative value",
		})
		return
	}
	//check if the amout to debit with is not 0
	if debitDetails.Amount <= 0 || debitDetails.Amount <= 0.0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "your debit amount cannot be 0 or less than 0",
		})
		return
	}

	username := c.Request.Header.Get("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": headerErr,
		})
		return
	}
	userResult, err := models.GetUserByUsername(&foundUser, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("the user with username %s does not exist", username),
		})
		log.Println(err.Error())
		return
	}

	params := c.Param("wallet_id")
	if params == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": paramsIdErr,
		})
		return
	}
	walletResult, err := models.GetUserWalletById(&wallet, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("the wallet with id of %s does not exist", params),
		})
		log.Println(err.Error())
		return
	}

	//check if the user the is owner of the wallet they are trying to debit from
	// as a user must not be allowed to debit from another persons wallet
	if userResult.Username == walletResult.Username {
		if walletResult.Wallet == 0 || walletResult.Wallet == 0.0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("%s your wallet is %v and so you cannot debit from it. please top up your wallet first.", userResult.Username, walletResult.Wallet),
			})
			return
		} else if walletResult.Wallet < debitDetails.Amount {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("%s you are trying to debit %v from your wallet which is %v. please top up your wallet first or reduce the amount you would like to debit.", userResult.Username, debitDetails.Amount, walletResult.Wallet),
			})
			return
		} else {
			previousBalance := decimal.NewFromFloat(walletResult.Wallet)
			debit := decimal.NewFromFloat(debitDetails.Amount)
			newBalance := previousBalance.Sub(debit).String()
			cacheErr := config.Cached.Set(walletResult.Username, newBalance, 0).Err()
			if cacheErr != nil {
				log.Println(cacheErr)
				return
			}

			//check if the fund to be sent contains a minus sign before authorizing the credit
			err := models.UpdateWalletFund(&wallet, username, newBalance) //the value in database has to be a float e.g 0.0
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "an error occured while updating your wallet",
				})
				log.Println(err.Error())
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"balance": newBalance,
				"message": fmt.Sprintf("%s your wallet has been debited with %v", userResult.First_Name, debitDetails.Amount),
			})

		}
	} else {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("the user with username %s is not the owner of th wallet with an id of %d", userResult.Username, walletResult.ID),
		})
		return
	}
}
