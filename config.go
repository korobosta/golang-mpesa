package main

import "database/sql"

type Config struct {
	Id                  int
	MpesaShortCode      string
	MpesaCallbackUrl    string
	MpesaConsumerSecret string
	MpesaConsumerKey    string
	MpesaPassKey        string
	Env                 int
	MpesaTokenUrl       string
	MpesaStkPushUrl     string
	Amount              float64
	AccountNumber       string
	PhoneNumber         string
	StkPushData         StkPushData
}

type StkPushData struct {
	MerchantRequestID   string
	CheckoutRequestID   string
	ResponseCode        string
	ResponseDescription string
	CustomerMessage     string
	Status              string
	DbConnection        *sql.DB
	TableName string
	PhoneNumber string
	AccountNumber string
	Amount string
}

type TokenFeedback struct {
	AccessToken   string
	Error         string
	Success       bool
	MpesaResponse map[string]string
}

type StkPushFeedback struct {
	Error         string
	Success       bool
	MpesaResponse map[string]string
}