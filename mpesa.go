package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)


var TEST_MPESA_TOKEN_URL string ="https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
var TEST_MPESA_STK_PUSH_URL ="https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"

var LIVE_MPESA_TOKEN_URL ="https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
var LIVE_MPESA_STK_PUSH_URL ="https://api.safaricom.co.ke/mpesa/stkpush/v1/processrequest"

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
										stkData.Status ="0"
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

		
