package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	
	config := Config{}

	config.MpesaConsumerKey = "kAPfIJxLk0Ieas12QLsANwXPHmbvNCwn"
	config.MpesaConsumerSecret = "23KqoLyCB3OHVvAI"
	config.MpesaCallbackUrl = "https://e515-197-156-137-142.ngrok-free.app/stk/callback"
	config.MpesaShortCode ="174379"
	config.PhoneNumber ="254713887070"
	config.Env =0
	config.AccountNumber = "12345"
	config.Amount =2
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
	stk.ReferenceNumber = "mpesa_reference"
	stk.DbConnection = Connect()
	stk.TableName ="mpesa_stk_payments"

	config.StkPushData = stk

	payment := Payment{}

	payment.BillRefNumber ="account_number"
	payment.BusinessShortCode ="business_short_code"
	payment.FirstName ="first_name"
	payment.InvoiceNumber ="invoice_number"
	payment.LastName ="last_name"
	payment.MiddleName ="middle_name"
	payment.TransTime ="transaction_time"
	payment.TransID ="transaction_id"
	payment.TransAmount ="amount"
	payment.ThirdPartyTransID ="third_party_trans_id"
	payment.MSISDN ="phone"
	payment.OrgAccountBalance ="organisation_account_balance"
	payment.TransactionType ="transaction_type"

	paymentTable := PaymentTable{}
	paymentTable.TableName = "mpesa_payments"
	paymentTable.Columns = payment
	paymentTable.DbConnection = Connect()

	router.Handle("/stk/callback", GetStkPushResponse(config)).Methods("POST","OPTIONS")
	router.Handle("/payment/confirmation", SaveMpesaPaymentConfirmation(paymentTable)).Methods("POST","OPTIONS")
	router.Handle("/payment/transaction-query", GetTransactionQueryResponse(paymentTable)).Methods("POST","OPTIONS")


	//feedback := GenerateMpesaToken(config)

	//feedback := StkPush(config)

	feedback :=TokenFeedback{}

	out, err := json.Marshal(feedback)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))

	http.Handle("/", router)
	fmt.Println("Connected to port 10000")
	log.Fatal(http.ListenAndServe(":10000", router))
}