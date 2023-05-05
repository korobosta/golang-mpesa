# What is golang-mpesa

**golang-mpesa** is an golang mpesa package that helps developers to call mpesa apis without much hussle. The package helps in generating tokens, making stk push, performing transaction query requests and even saving the transactions in the database.

## How to setup

Get the package with following command :

```bash
go get github.com/korobosta/golang-mpesa

```

## Functions

The package provides the following functions:

- **GenerateMpesaToken** : For generating mpesa access tokens

- **StkPush** : For making stk push to the customer phone

- **GetStkPushResponse** : For getting the stk push callback data from mpesa and saving it in the database.

- **DecodeStkCallbackResponse**  A function for decoding mpesa stk response data.

- **SaveMpesaPaymentConfirmation** : This functions takes the mpesa payment confirmation response and saves it in the database

- **TransactionQuery** : This function makes the transaction query request to MPESA

- **GetTransactionQueryResponse** : This is a callback function for taking transaction query response and saving the payment in the database.

- **EncryptWithPublicKey** : This function encrypts the password with MPESA public key to form the security credentail


Configuration:

For easy configuration of the package, we have created a model called ***Config*** so that all static variables required can be initialized at once. The configuration can be done as shown below.

```Go

import (
	"github.com/korobosta/golang-mpesa"
)

//Configuration function for all static mpesa variables
func MpesaConfig() mpesa.Config {
	var config = mpesa.Config{}

	//Mandatory Fiels
	config.MpesaConsumerKey = "your-consumer-key"
	config.MpesaConsumerSecret = "your-consumer-secret"
	config.MpesaShortCode = "your-mpesa-short code" //e.g 174379
	config.Initiator = "Your initiator username" // e.g testapi
	config.InitiatorPassword = "Your Initation Password" // e.g Safaricom1234!!
	config.TransQueryOriginatorConversationID ="your transaction query originator conversation id" //  Your can use our ***mpesa.RandomString(10)*** function to ganarate a random string 
	config.TransQueryRemarks = "Transaction Query Request Remarks"
	config.IdentifierType = "4" // Transaction Query Identity type can be 1, 2, 3, 4(Business Short Code)
	config.TransQueryCommandID = "TransactionStatusQuery"
	config.TransQueryResultURL = "Your Result URL for handling transaction query" 
	config.TransQueryQueueTimeOutURL = "Your QueueTimeOut transaction query callback url"
	config.MpesaCallbackUrl = "Your url for handling stk callback"
	config.MpesaPassKey = "your mpesa pass key"
	config.Env = 0 // 0- Test, 1- Live

	//Optional

	//I you would like the package to decode and save the transaction query response from safaricom to your database table, configure the model as shown below
	transQueryTable := mpesa.TransQueryTable{}

	transQueryTableColumns := mpesa.TransQueryTableColumns{}

	//Replace the details below with the way your mpesa transaction query table columns
	transQueryTableColumns.AccountReference = "account_reference"  // 
	transQueryTableColumns.ConversationID = "ConversationID"
	transQueryTableColumns.OriginatorConversationID = "OriginatorConversationID"
	transQueryTableColumns.ResponseCode = "ResponseCode"
	transQueryTableColumns.ResponseDescription = "ResponseDescription"
	transQueryTableColumns.Status = "status"
	transQueryTableColumns.TransactionReference = "mpesa_code_entered"
	transQueryTableColumns.AccountReference = "account_reference"

	transQueryTable.Columns = transQueryTableColumns
	transQueryTable.TableName = "transaction_query_requests" //Name of your table
	transQueryTable.DbConnection = Connect()
	transQueryTable.DefaultStatus = "0"  // Default status for the status column
	transQueryTable.SuccessMpesaStatus = "1" // The status we will save when mpesa brings back success response

	config.TransQueryTable = transQueryTable


	//If you would like the package to save the stk data response in your database table, configure as shown below

	stk := mpesa.StkPushData{}
	//Replace the details below with the column names of your table
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
	stk.DbConnection = MyDBConnect() // Your Function that returns your database connection
	stk.TableName ="mpesa_stk_payments" // Your table name
	stk.DefaultStatus = "0" //Status that will be saved at the beginning
	stk.SuccessMpesaStatus = "1" // Status that will be saved after a successfull mpesa respnse


	config.StkPushData = stk
	

	config.TransQueryTable = transQueryTable

	return config
}

// A function for configuring your MPESA payments table
func MpesaPaymentTable() mpesa.PaymentTable {
	payment := mpesa.Payment{}

	//Replace with you mpesa payments table column names
	payment.BillRefNumber = "account_number"
	payment.BusinessShortCode = "business_short_code"
	payment.FirstName = "first_name"
	payment.InvoiceNumber = "invoice_number"
	payment.LastName = "last_name"
	payment.MiddleName = "middle_name"
	payment.TransTime = "transaction_time"
	payment.TransID = "transaction_id"
	payment.TransAmount = "amount"
	payment.ThirdPartyTransID = "third_party_trans_id"
	payment.MSISDN = "phone"
	payment.OrgAccountBalance = "organisation_account_balance"
	payment.TransactionType = "transaction_type"

	paymentTable := mpesa.PaymentTable{}
	paymentTable.TableName = "mpesa_payments"  //Your mpesa payments table name
	paymentTable.Columns = payment
	paymentTable.DbConnection = MyDBConnect() // Your function for return DB connection instance

	return paymentTable
}

func main() {
	config := MpesaConfig()

	//Generating Tokens

	// All the functions that call MPESA api will generate the tokens on there own but incase you want to generate the token, call the functon below

	tokenFeedback := mpesa.TokenFeedback{}
	tokenFeedback = mpesa.GenerateMpesaToken(config)
	if tokenFeedback.Success == false{
		errors = tokenFeedback.Error
	}else{
		accessToken =tokenFeedback.AccessToken
	}

	//STK Push

	mpesaResponse := mpesa.StkPushFeedback{}
	config.PhoneNumber = "254712345678" // Customer phone number
	config.Amount = 1 //Amount for cutomer to pay
	config.AccountNumber = "12345" //Account number which the customer is paying e.g order number
	config.MpesaShortCode = "174379" //Short code receiving the payment , You can use config.MpesaShortCode if you configured it above
	mpesaResponse = mpesa.StkPush(config) // Call the package function to do the stk push

	if(mpesaResponse.Success == true){
		//STK was sent successfully
	}
	else{
		// log.Println(mpesaResponse.Error)
	}

	//mpesaResponse.MpesaResponse return the raw response from MPESA

	// If you would like the package to handle the stk transansaction response from mpesa and save it in the database, make a route as shown below and call the package function.

	// For this to happen your config should have configured config.StkPushData above

	router := mux.NewRouter()
	router.Handle("/stk/callback", mpesa.GetStkPushResponse(config)).Methods("POST","OPTIONS")


	// Transaction Query 

	transactionQueryResponse := mpesa.TransactionQueryFeedback{}

	config.TransactionReference = ""  //Transaction Code your are checking"
	config.AccountNumber = "" //Account number related to the transaction

	transactionQueryResponse = mpesa.TransactionQuery(config)

	if(transactionQueryResponse.Success == true){
		//Query sent successfully
	}
	else{
		// log.Println(transactionQueryResponse.Error)
	}

	//transactionQueryResponse.MpesaResponse return the raw response from MPESA

	// If you would like the package to handle the callback transansaction response from mpesa and save it in the database, make a route as shown below and call the package function.

	// For this to happen your config should have configured config.TransQueryTable above

	router := mux.NewRouter()
	router.Handle("/payment/transaction-query", mpesa.GetTransactionQueryResponse(paymentTable,config.TransQueryTable,config)).Methods("POST","OPTIONS")


	// MPESA Confirmation
	// If you would like the package  to process for you the data from mpesa confirmation response, configure your confiramtion url as shown below

	paymentTable := MpesaPaymentTable()
	router.Handle("/payment/confirmation", mpesa.SaveMpesaPaymentConfirmation(paymentTable,config)).Methods("POST","OPTIONS")
	
	

	
}

```

## Todo list
- Connection to other databases that are not sql
- More updates to give more controll to the database
- Any suggestion you make through ***kevinkorobosta@gmail.com***
