package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "pck_payments"
    dbHost := "127.0.0.1:3306"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)

	if err != nil {
		panic(err.Error())
	}
	
	return db
}

func main() {
	config := Config{}

	config.MpesaConsumerKey = "kAPfIJxLk0Ieas12QLsANwXPHmbvNCwn"
	config.MpesaConsumerSecret = "23KqoLyCB3OHVvAI"
	config.MpesaCallbackUrl = "https://mydomain.com:7000"
	config.MpesaShortCode ="174379"
	config.PhoneNumber ="254713887070"
	config.Env =0
	config.AccountNumber = "12345"
	config.Amount =1
	config.MpesaPassKey ="bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919"

	stk := StkPushData{}
	stk.CheckoutRequestID = "CheckoutRequestID"
	stk.CustomerMessage = "CustomerMessage"
	stk.MerchantRequestID ="MerchantRequestID"
	stk.ResponseCode = "ResponseCode"
	stk.ResponseDescription = "ResponseDescription"
	stk.Status ="status"
	stk.PhoneNumber ="phone_number"
	stk.AccountNumber ="account_number"
	stk.Amount = "amount"
	stk.DbConnection = Connect()
	stk.TableName ="mpesa_stk_payments"

	config.StkPushData = stk


	//feedback := GenerateMpesaToken(config)

	feedback := StkPush(config)

	out, err := json.Marshal(feedback)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}