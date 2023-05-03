package mpesa

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
	Initiator           string
	InitiatorPassword   string
	TransactionReference string
	IdentifierType      string
	TransQueryResultURL string
	TransQueryRemarks   string
	TransQueryOccassion string
	TransQueryCommandID string
	TransQueryQueueTimeOutURL     string
	TransQueryOriginatorConversationID  string
	StkPushData         StkPushData
	TransQueryTable TransQueryTable
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
	ReferenceNumber string
	DefaultStatus string
	SuccessMpesaStatus string
}

type TransQueryTableColumns struct {
	OriginatorConversationID   string
	ConversationID   string
	ResponseCode        string
	ResponseDescription string
	TransactionReference string
	Status string
	AccountReference string
}

type TransQueryTable struct{
	TableName string
	Columns TransQueryTableColumns
	DbConnection        *sql.DB
	DefaultStatus string
	SuccessMpesaStatus string

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

type TransactionQueryFeedback struct {
	Error         string
	Success       bool
	MpesaResponse map[string]string
}

type StkCallbackResponse struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string `form:"MerchantRequestID" json:"MerchantRequestID"`
			CheckoutRequestID string `form:"CheckoutRequestID" json:"CheckoutRequestID"`
			ResultCode int `form:"ResultCode" json:"ResultCode"`
			ResultDesc string `form:"ResultDesc" json:"ResultDesc"`
			CallbackMetadata struct {
				Item [] struct {
					Name string `form:"Name" json:"Name"`
					Value any  `form:"Value" json:"Value"`
				}
			} `form:"CallbackMetadata" json:"CallbackMetadata"`
		}
	} 
}

type TransactionStatusResponse struct {
	Result struct {
		TransactionID string `form:"TransactionID" json:"TransactionID"`
		ConversationID string `form:"ConversationID" json:"ConversationID"`
		OriginatorConversationID string `form:"OriginatorConversationID" json:"OriginatorConversationID"`
		ResultType int `form:"ResultType" json:"ResultType"`
		ResultDesc string `form:"ResultDesc" json:"ResultDesc"`
		ResultParameters struct {
			ResultParameter [] struct {
				Key string `form:"Key" json:"Key"`
				Value any  `form:"Value" json:"Value"`
			} `form:"ResultParameter" json:"ResultParameter"`
		}  `form:"ResultParameters" json:"ResultParameters"`
		
	} 
}

type Payment struct {
	TransactionType string `form:"TransactionType" json:"TransactionType"`
    TransID string `form:"TransID" json:"TransID"`
    TransTime string `form:"TransTime" json:"TransTime"`
    TransAmount string `form:"TransAmount" json:"TransAmount"`
    BusinessShortCode string `form:"BusinessShortCode" json:"BusinessShortCode"`
    BillRefNumber string `form:"BillRefNumber" json:"BillRefNumber"`
    InvoiceNumber string `form:"InvoiceNumber" json:"InvoiceNumber"`
    OrgAccountBalance string `form:"OrgAccountBalance" json:"OrgAccountBalance"`
    ThirdPartyTransID string `form:"ThirdPartyTransID" json:"ThirdPartyTransID"`
    MSISDN string `form:"MSISDN" json:"MSISDN"`
    FirstName string `form:"FirstName" json:"FirstName"`
    MiddleName string `form:"MiddleName" json:"MiddleName"`
    LastName string `form:"LastName" json:"LastName"`
}

type PaymentTable struct{
	TableName string
	Columns Payment
	DbConnection        *sql.DB

}

type FormatedStkCallback struct {
	MpesaReceiptNumber any `form:"MpesaReceiptNumber" json:"MpesaReceiptNumber"`
    Amount any `form:"Amount" json:"Amount"`
    TransactionDate any `form:"TransactionDate" json:"TransactionDate"`
    PhoneNumber any `form:"PhoneNumber" json:"PhoneNumber"`
    MerchantRequestID string `form:"MerchantRequestID" json:"MerchantRequestID"`
    CheckoutRequestID string `form:"CheckoutRequestID" json:"CheckoutRequestID"`
    ResultCode int `form:"ResultCode" json:"ResultCode"`
    ResultDesc string `form:"ResultDesc" json:"ResultDesc"`
}

type StkCallbackFeedback struct {
	Error         string
	Success       bool
	MpesaResponse map[string]any
}
