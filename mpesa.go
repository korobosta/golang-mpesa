package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"time"
)


var TEST_MPESA_TOKEN_URL string ="https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
var TEST_MPESA_STK_PUSH_URL ="https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
var TEST_MPESA_TRANSACTION_QUERY_URL = "https://sandbox.safaricom.co.ke/mpesa/transactionstatus/v1/query"

var LIVE_MPESA_TOKEN_URL ="https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
var LIVE_MPESA_STK_PUSH_URL ="https://api.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
var LIVE_MPESA_TRANSACTION_QUERY_URL = "https://api.safaricom.co.ke/mpesa/transactionstatus/v1/query"

func GenerateMpesaToken(config Config) TokenFeedback {

	var errors string = ""
	var accessToken string = ""
	var success bool = false
	m := make(map[string]string)
	feedback := TokenFeedback{}

	var url string = config.MpesaTokenUrl
	if(url == ""){
		if config.Env == 1 {
			url = LIVE_MPESA_TOKEN_URL
		} else {
			url = TEST_MPESA_TOKEN_URL
		}
	}
	
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		errors = err.Error()
	}else{

		req.Header.Add("Content-Type", "application/json")

		req.SetBasicAuth(config.MpesaConsumerKey, config.MpesaConsumerSecret)

		res, err := client.Do(req)
		if err != nil{
			errors = err.Error()
		}else{
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)

			if err != nil{
				errors = err.Error()
			}else{
				err = json.Unmarshal(body, &m)
				if err != nil{
					errors = err.Error()
				}else{

					if token, ok := m["access_token"]; ok {
						accessToken = token
						success = true
					}
					if errorMessage, ok := m["errorMessage"]; ok {
						errors = errorMessage
					}
				}
			}
		}
	}
	feedback.AccessToken=accessToken
	feedback.Error =errors
	feedback.Success = success
	feedback.MpesaResponse = m

	return feedback
}

func StkPush(config Config) StkPushFeedback{

	var errors string = ""
	var accessToken string = ""
	var success bool = false
	m := make(map[string]string)

	tokenFeedback := TokenFeedback{}
	stkPushFeedback := StkPushFeedback{}

	tokenFeedback = GenerateMpesaToken(config)

	if tokenFeedback.Success == false{
		errors = tokenFeedback.Error
	}else{
		accessToken =tokenFeedback.AccessToken
		var url string =""
		if config.MpesaStkPushUrl==""{
			if config.Env  == 1{
				url = LIVE_MPESA_STK_PUSH_URL
			}else{
				url = TEST_MPESA_STK_PUSH_URL
			}
		}
		
		method := "POST"

		timestamp := time.Now().Format("20060102150405")

		password := base64.StdEncoding.EncodeToString([]byte(config.MpesaShortCode+config.MpesaPassKey+timestamp))
	
		payload := map[string]any{
			"BusinessShortCode": config.MpesaShortCode,
			"Password": password ,
			"Timestamp": timestamp,
			"TransactionType": "CustomerPayBillOnline",
			"Amount": config.Amount,
			"PartyA": config.PhoneNumber,
			"PartyB": config.MpesaShortCode,
			"PhoneNumber": config.PhoneNumber,
			"CallBackURL": config.MpesaCallbackUrl,
			"AccountReference": config.AccountNumber,
			"TransactionDesc": "Payment for "+config.AccountNumber,
		}
	
		jsonPayload, err := json.Marshal(payload)

		if err != nil {
			errors =err.Error()
		}else{
			client := &http.Client {
				}
			req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
			if err != nil {
				errors = err.Error()
			}else{
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", "Bearer "+accessToken)
				res, err := client.Do(req)
				if err != nil {
					errors = err.Error()
				}else{
					defer res.Body.Close()
		            body, err := ioutil.ReadAll(res.Body)
					if err != nil {
						errors = err.Error()
					}else{
						err = json.Unmarshal(body, &m)
						if err != nil {
							errors = err.Error()
						}else{
							if responseCode, ok := m["ResponseCode"]; ok {
								if responseCode == "0"{
									success = true
									stkData := StkPushData{};
									if(config.StkPushData != stkData){
										stkData.Status =config.StkPushData.DefaultStatus
										stkData.CheckoutRequestID =m["CheckoutRequestID"]
										stkData.CustomerMessage = m["CustomerMessage"]
										stkData.MerchantRequestID = m["MerchantRequestID"]
										stkData.ResponseCode = m["ResponseCode"]
										stkData.ResponseDescription= m["ResponseDescription"]
										stkData.PhoneNumber =config.PhoneNumber
										stkData.AccountNumber =config.AccountNumber
										stkData.Amount = fmt.Sprintf("%f", config.Amount)
										SaveStkPushData(stkData,config.StkPushData)
									}
								}
							}else{
								if errorMessage, ok := m["errorMessage"]; ok {
									errors = errorMessage
								}else{
									errors ="Unknown error occured"
								}
							}
						}
					}
				}

			}
		}
	}

	stkPushFeedback.Error = errors
	stkPushFeedback.MpesaResponse =m
	stkPushFeedback.Success = success
	
	return stkPushFeedback
}

func SaveStkPushData(data StkPushData,db StkPushData){
	save, err := db.DbConnection.Prepare("INSERT INTO "+db.TableName+"("+db.MerchantRequestID+","+db.CheckoutRequestID+","+db.ResponseCode+","+db.ResponseDescription+","+db.CustomerMessage+","+db.Status+","+db.PhoneNumber+","+db.AccountNumber+","+db.Amount+") VALUES(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Println(err)
	}
	save.Exec(data.MerchantRequestID, data.CheckoutRequestID, data.ResponseCode, data.ResponseDescription,data.CustomerMessage,data.Status,data.PhoneNumber,data.AccountNumber,data.Amount)
}

func GetStkPushResponse(config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var errors string = ""
	var success bool = false

	m := make(map[string]any)

	stkCallbackFeedback := StkCallbackFeedback{}

	formatedStkCallback :=FormatedStkCallback{}

	
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		errors =err.Error()
	}

	if(errors == ""){
		err = json.Unmarshal(b, &m)
		if err != nil {
			errors = err.Error()
		}else{
			formatedStkCallback = DecodeStkCallbackResponse(b)
			//stkData := FormatedStkCallback{};
			if(formatedStkCallback.MpesaReceiptNumber != nil){
				success = true
			}else{
				errors = "Some fields are missing"
			}
		}
		
    }

	if (formatedStkCallback.MpesaReceiptNumber != nil){
		updateOrder, err := config.StkPushData.DbConnection.Prepare("UPDATE "+config.StkPushData.TableName+" SET status = ?,"+config.StkPushData.ReferenceNumber+"=? where "+config.StkPushData.MerchantRequestID+"=?")
		if err != nil {
			log.Println(err)
		}
		updateOrder.Exec(config.StkPushData.SuccessMpesaStatus, formatedStkCallback.MpesaReceiptNumber,formatedStkCallback.MerchantRequestID)
	}
	stkCallbackFeedback.Error =errors
	stkCallbackFeedback.MpesaResponse = m
	stkCallbackFeedback.Success = success

	out, err := json.Marshal(formatedStkCallback)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
})

	//return stkCallbackFeedback

}

func DecodeStkCallbackResponse (b []byte) FormatedStkCallback{

	formatedStkCallback := FormatedStkCallback{}

	stkCallbackResponse := StkCallbackResponse{}

	err := json.Unmarshal(b, &stkCallbackResponse)
	if err != nil {
		log.Println(err)
	}else{
		formatedStkCallback.Amount =stkCallbackResponse.Body.StkCallback.CallbackMetadata.Item[0].Value
		formatedStkCallback.CheckoutRequestID =stkCallbackResponse.Body.StkCallback.CheckoutRequestID
		formatedStkCallback.MerchantRequestID = stkCallbackResponse.Body.StkCallback.MerchantRequestID
		formatedStkCallback.MpesaReceiptNumber = stkCallbackResponse.Body.StkCallback.CallbackMetadata.Item[1].Value
		formatedStkCallback.PhoneNumber  = stkCallbackResponse.Body.StkCallback.CallbackMetadata.Item[3].Value
		formatedStkCallback.ResultCode = stkCallbackResponse.Body.StkCallback.ResultCode
		formatedStkCallback.TransactionDate = stkCallbackResponse.Body.StkCallback.CallbackMetadata.Item[2].Value
		formatedStkCallback.ResultDesc = stkCallbackResponse.Body.StkCallback.ResultDesc
	}

	return formatedStkCallback
    
}

func SaveMpesaPaymentConfirmation(table PaymentTable,config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var errors string = ""

	payment := Payment{}
	
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		errors =err.Error()
	}

	if(errors == ""){
		err = json.Unmarshal(b, &payment)
		if err != nil {
			errors = err.Error()
		}else{
			SaveMpesaPayment(table,payment,config);
		}
		
    }

	out, err := json.Marshal(payment)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
})

	//return stkCallbackFeedback

}

func GetTransactionQueryResponse(table PaymentTable,transQueryTable TransQueryTable, config Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var errors string = ""

	transactionQueryCallbackFeedback := TransactionStatusResponse{}

	var originatorConversationID string = ""

	payment :=Payment{}

	
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		errors =err.Error()
	}

	if(errors == ""){
		err = json.Unmarshal(b, &transactionQueryCallbackFeedback)
		if err != nil {
			errors = err.Error()
		}else{
			log.Println(transactionQueryCallbackFeedback)
			payment,originatorConversationID = DecodeTransactionQueryCallbackResponse(b)
			
			emptyPayment := Payment{};
			if(payment != emptyPayment){
				var  accountNumber string = UpdateTransQueryTable(transQueryTable,originatorConversationID)
				if accountNumber != ""{
					payment.BillRefNumber = accountNumber
					SaveMpesaPayment(table,payment,config);
				}
				
			}
		}
		
    }
	
	// out, err := json.Marshal(payment)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(out))
})

	//return stkCallbackFeedback

}

func DecodeTransactionQueryCallbackResponse (b []byte) (Payment,string){

	payment := Payment{}

	var originatorConversationID string = ""

	transactionQueryCallbackResponse := TransactionStatusResponse{}

	err := json.Unmarshal(b, &transactionQueryCallbackResponse)
	if err != nil {
		log.Println(err)
	}else{
		resultType := transactionQueryCallbackResponse.Result.ResultType
		if(resultType ==0){
			paybill := transactionQueryCallbackResponse.Result.ResultParameters.ResultParameter[0].Value.(string)
			log.Println(paybill)
			arrPaybill := strings.Split(paybill, "-")
			businessShortCode := arrPaybill[0]

			originatorConversationID = transactionQueryCallbackResponse.Result.OriginatorConversationID

			customer := transactionQueryCallbackResponse.Result.ResultParameters.ResultParameter[1].Value.(string)
			arrCustomer := strings.Split(customer, "-")
			phoneNumber := arrCustomer[0]
			name := arrCustomer[1]
			name = strings.TrimSpace(name)
			nameErr := strings.Split(name, " ")
			firstName := nameErr[0]
			lastName := nameErr[len(nameErr)-1]

			payment.FirstName = firstName
			payment.LastName =lastName
			payment.BusinessShortCode =businessShortCode
			payment.MSISDN =phoneNumber
			payment.TransID = transactionQueryCallbackResponse.Result.ResultParameters.ResultParameter[12].Value.(string)
			payment.TransAmount = transactionQueryCallbackResponse.Result.ResultParameters.ResultParameter[10].Value.(string)
			payment.TransTime =transactionQueryCallbackResponse.Result.ResultParameters.ResultParameter[9].Value.(string)
		}
	}

	return payment,originatorConversationID
}

func TransactionQuery(config Config) TransactionQueryFeedback{

	var errors string = ""
	var accessToken string = ""
	var success bool = false
	m := make(map[string]string)

	tokenFeedback := TokenFeedback{}
	transactionQueryFeedback := TransactionQueryFeedback{}

	tokenFeedback = GenerateMpesaToken(config)

	if tokenFeedback.Success == false{
		errors = tokenFeedback.Error
	}else{
		accessToken =tokenFeedback.AccessToken
		var url string =""
		if config.MpesaStkPushUrl==""{
			if config.Env  == 1{
				url = LIVE_MPESA_TRANSACTION_QUERY_URL
			}else{
				url = TEST_MPESA_TRANSACTION_QUERY_URL
			}
		}
		
		method := "POST"
	
		payload := map[string]any{
			"Initiator":config.Initiator,
			"SecurityCredential": EncryptWithPublicKey(config.InitiatorPassword,config.Env),
			"CommandID": config.TransQueryCommandID,
			"TransactionID": config.TransactionReference,
			"OriginatorConversationID":config.TransQueryOriginatorConversationID,
			"PartyA":config.MpesaShortCode,
			"IdentifierType":config.IdentifierType,
			"ResultURL":config.TransQueryResultURL,
			"QueueTimeOutURL":config.TransQueryQueueTimeOutURL,
			"Remarks":config.TransQueryRemarks,
			"Occasion":config.TransQueryOccassion,
		}
	
		jsonPayload, err := json.Marshal(payload)

		if err != nil {
			errors =err.Error()
		}else{
			client := &http.Client {
				}
			req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
			if err != nil {
				errors = err.Error()
			}else{
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Authorization", "Bearer "+accessToken)
				res, err := client.Do(req)
				if err != nil {
					errors = err.Error()
				}else{
					defer res.Body.Close()
		            body, err := ioutil.ReadAll(res.Body)
					if err != nil {
						errors = err.Error()
					}else{
						err = json.Unmarshal(body, &m)
						if err != nil {
							errors = err.Error()
						}else{
							if responseCode, ok := m["ResponseCode"]; ok {
								if responseCode == "0"{
									success = true
									transQueryTable := TransQueryTable{};
									if(config.TransQueryTable != transQueryTable){
										transQueryData := TransQueryTableColumns{}
										transQueryData.Status =config.TransQueryTable.DefaultStatus
										transQueryData.ResponseDescription =m["ResponseDescription"]
										transQueryData.ConversationID = m["ConversationID"]
										transQueryData.OriginatorConversationID = m["OriginatorConversationID"]
										transQueryData.TransactionReference = config.TransactionReference
										transQueryData.AccountReference= config.AccountNumber
										SaveTransQueryData(transQueryData,config.TransQueryTable)
									}
								}else{
									if errorMessage, ok := m["errorMessage"]; ok {
										errors = errorMessage
									}else{
										errors ="Unknown error occured"
									}
								}
							}
						}
					}
				}

			}
		}
	}

	transactionQueryFeedback.Error = errors
	transactionQueryFeedback.MpesaResponse = m
	transactionQueryFeedback.Success = success
	
	return transactionQueryFeedback
}

func SaveTransQueryData(data TransQueryTableColumns,db TransQueryTable){
	save, err := db.DbConnection.Prepare("INSERT INTO "+db.TableName+"("+db.Columns.ConversationID+","+db.Columns.OriginatorConversationID+","+db.Columns.ResponseCode+","+db.Columns.AccountReference+","+db.Columns.ResponseDescription+","+db.Columns.Status+","+db.Columns.TransactionReference+") VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		log.Println(err)
	}
	save.Exec(data.ConversationID, data.OriginatorConversationID, data.ResponseCode, data.AccountReference,data.ResponseDescription,data.Status,data.TransactionReference)
}


func UpdateTransQueryTable(transQueryTable TransQueryTable, originatorConversationID string) string{
	var accountNumber string= GetTransactionQueryAccountNumber(transQueryTable, originatorConversationID)
	if(accountNumber == ""){
		updateOrder, err := transQueryTable.DbConnection.Prepare("UPDATE "+transQueryTable.TableName+" SET status = ?,"+transQueryTable.SuccessMpesaStatus+"=? where "+transQueryTable.Columns.OriginatorConversationID+"=?")
		if err != nil {
			log.Println(err)
		}else{
			_,err = updateOrder.Exec(transQueryTable.SuccessMpesaStatus,originatorConversationID)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return accountNumber
}

func GetTransactionQueryAccountNumber(table TransQueryTable, originatorConversationID string) string{
	var accountNumber string = ""
	err := table.DbConnection.QueryRow("select "+table.Columns.AccountReference+" from "+table.TableName+" where "+table.Columns.OriginatorConversationID+" = ?", originatorConversationID).Scan(&accountNumber)
	if err != nil {
		log.Println(err)
	}
	return accountNumber
}

func SaveMpesaPayment(table PaymentTable, payment Payment,config Config){
	columns := table.Columns.TransactionType+","+table.Columns.TransID+","+table.Columns.TransTime+","+table.Columns.TransAmount+","+table.Columns.BusinessShortCode+","+table.Columns.BillRefNumber+","+table.Columns.InvoiceNumber+","+table.Columns.OrgAccountBalance+","+table.Columns.ThirdPartyTransID+","+table.Columns.MSISDN+","+table.Columns.FirstName+","+table.Columns.MiddleName+","+table.Columns.LastName
	log.Println(columns)
	save, err := table.DbConnection.Prepare("INSERT INTO "+table.TableName+"("+columns+") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Println(err)
	}
	_, err = save.Exec(payment.TransactionType,payment.TransID,payment.TransTime,payment.TransAmount,payment.BusinessShortCode,payment.BillRefNumber,payment.InvoiceNumber,payment.OrgAccountBalance,payment.ThirdPartyTransID,payment.MSISDN,payment.FirstName,payment.MiddleName,payment.LastName)
	if err == nil{
		if config.AfterPaymentFunction != nil{
			config.AfterPaymentFunction(payment)
		}
	}
}
		
