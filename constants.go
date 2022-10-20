package spacecow_common

import (
	"time"
)

const MinHoursBetweenScans = 3
const DefaultGopiSleep = 60
const SkimScanDefaultInterval = 60 * 24 * 1       // regular - ads
const LowFatScanDefaultInterval = 60 * 24 * 1     // regular - no ads
const WholeScanDefaultInterval = 60 * 12 * 1      // premium
const HyperpasturizedScanDefaultInterval = 60 * 1 // gomez level
type SubscriptionLevel int

const TopicName = "queue"
const SubName = "queue-sub"
const FbTopicName = "fbupdate"   // "fbupdate"
const FbSubName = "fbupdate-sub" // "fbupdate"

const (
	Trial = iota
	Skim
	LowFat
	Whole
	Hyperpasturized
)

type LanguageCode int

const (
	En = iota
)

type EventTypes int

const (
	EventCheckBalances = iota
	EventLoadNewTransactions
	SetupNewInstitution
	EventNoop
	EventSkimFound
	EventAlert
	EventUpdateMap
	EventRefreshUI
	EventIntrospect
	EventPushAccounts
	EventRelinkAccount
	EventRelinkCompleted
	EventExportTransactions
)

// XactionType is the type of xaction in the db
type XactionType int

const (
	XactionPayment = iota
	XactionCharge
	XactionCredit
	XactionInterestCharge
	XactionLateFee
)

// ShipTypes is the ship level in the game
type ShipTypes int

// classes on richness levels
const (
	ScoutShip    = iota // less than 5k
	AberdeenShip        // 5k-50k
	HerefordShip        // 50-100k
	LonghornShip        // 100k-1m
	WagyuShip           // > 1 million
)

// TransactionMap fixes the broken Plaid transaction mapping
type TransactionMap struct {
	ID                  string      `bson:"_id" json:"id"`
	PhysicalLocation    string      `json:"physical_location" bson:"physical_location"`
	TransactionType     XactionType `bson:"transaction_type" json:"transaction_type"`
	Description         string      `json:"description" bson:"description"`
	DetailedDescription []string    `bson:"detailed_description" json:"detailed_description"`
}

// Q is a mapping for the internal work queue
type Q struct {
	ID        string     `json:"id" bson:"_id"` // id
	Done      bool       `json:"done" bson:"done"`
	Added     time.Time  `json:"added" bson:"added"`
	Processed time.Time  `json:"processed" bson:"processed"`
	UID       string     `json:"uid" bson:"uid"`
	AT        string     `json:"AT" bson:"at"`
	Event     EventTypes `json:"event" bson:"event"`
	Extra     string     `json:"extra" bson:"extra"`
}

type PlaidError struct {
	DisplayMessage  string `json:"display_message"`
	ErrorCode       string `json:"error_code"`
	ErrorMessage    string `json:"error_message"`
	ErrorType       string `json:"error_type"`
	RequestID       string `json:"request_id"`
	SuggestedAction string `json:"suggested_action"`
}

// CowTransaction is our xaction storage type with extensions from plaid
type CowTransaction struct {
	// The ID of a posted transaction's associated pending transaction, where applicable.
	PendingTransactionID string `json:"pending_transaction_id" bson:"pendingTransactionId"`
	// The ID of the category to which this transaction belongs. For a full list of categories, see [`/categories/get`](https://plaid.com/docs/api/products/transactions/#categoriesget).  If the `transactions` object was returned by an endpoint such as `/asset_report/get/` or `/asset_report/pdf/get`, this field will only appear in an Asset Report with Insights.
	CategoryID string `json:"category_id" bson:"categoryId"`
	// A hierarchical array of the categories to which this transaction belongs. For a full list of categories, see [`/categories/get`](https://plaid.com/docs/api/products/transactions/#categoriesget).  If the `transactions` object was returned by an endpoint such as `/asset_report/get/` or `/asset_report/pdf/get`, this field will only appear in an Asset Report with Insights.
	Category    []string `json:"category" bson:"category"`
	Address     string   `json:"address" bson:"address"`
	GeoHash     string   `bson:"geoHash" bson:"geoHash"`
	PaymentMeta string   `json:"payment_meta"` // only for ACH xfers, containers payment id and others
	// The merchant name or transaction description.  If the `transactions` object was returned by an endpoint such as `/transactions/get`, this field will always appear. If the `transactions` object was returned by an endpoint such as `/asset_report/get/` or `/asset_report/pdf/get`, this field will only appear in an Asset Report with Insights.
	Name string `json:"name" bson:"name"`
	// The string returned by the financial institution to describe the transaction. For transactions returned by `/transactions/get`, this field is in beta and will be omitted unless the client is both enrolled in the closed beta program and has set `options.include_original_description` to `true`.
	OriginalDescription string `json:"original_description" bson:"originalDescription"`
	// The ID of the account in which this transaction occurred.
	AccountID string `json:"account_id" bson:"accountId"`
	// The settled value of the transaction, denominated in the account's currency, as stated in `iso_currency_code` or `unofficial_currency_code`. Positive values when money moves out of the account; negative values when money moves in. For example, debit card purchases are positive; credit card payments, direct deposits, and refunds are negative.
	Amount float64 `json:"amount" bson:"amount"`
	// The ISO-4217 currency code of the transaction. Always `null` if `unofficial_currency_code` is non-null.
	IsoCurrencyCode string `json:"iso_currency_code" bson:"isoCurrencyCode"`
	// The unofficial currency code associated with the transaction. Always `null` if `iso_currency_code` is non-`null`. Unofficial currency codes are used for currencies that do not have official ISO currency codes, such as cryptocurrencies and the currencies of certain countries.  See the [currency code schema](https://plaid.com/docs/api/accounts#currency-code-schema) for a full listing of supported `iso_currency_code`s.
	UnofficialCurrencyCode string `json:"unofficial_currency_code" bson:"unofficialCurrencyCode"`
	// For pending transactions, the date that the transaction occurred; for posted transactions, the date that the transaction posted. Both dates are returned in an [ISO 8601](https://wikipedia.org/wiki/ISO_8601) format ( `YYYY-MM-DD` ).
	Date string `json:"date" bson:"date" bson:"date"`
	// When `true`, identifies the transaction as pending or unsettled. Pending transaction details (name, type, amount, category ID) may change before they are settled.
	Pending bool `json:"pending" bson:"pending"`
	// TransactionID The unique ID of the transaction.
	TransactionID string `json:"transaction_id" bson:"_id"`
	// The merchant name, as extracted by Plaid from the `name` field.
	MerchantName string `json:"merchant_name" bson:"merchantName"`
	// The check number of the transaction. This field is only populated for check transactions.
	CheckNumber string `json:"check_number" bson:"checkNumber"`
	// The channel used to make a payment. `online:` transactions that took place online.  `in store:` transactions that were made at a physical location.  `other:` transactions that relate to banking, e.g. fees or deposits.  This field replaces the `transaction_type` field.
	PaymentChannel string `json:"payment_channel" bson:"paymentChannel"`
	// The date that the transaction was authorized. Dates are returned in an [ISO 8601](https://wikipedia.org/wiki/ISO_8601) format ( `YYYY-MM-DD` ). The `authorized_date` field uses machine learning to determine a transaction date for transactions where the `date_transacted` is not available. If the `date_transacted` field is present and not `null`, the `authorized_date` field will have the same value as the `date_transacted` field.
	AuthorizedDate string `json:"authorized_date" bson:"authorizedDate"`
	// Date and time when a transaction was authorized in [ISO 8601]
	AuthorizedDatetime time.Time `json:"authorized_datetime" bson:"authorizedDatetime"`
	// Date and time when a transaction was posted in [ISO 8601](https://wikipedia.org/wiki/ISO_8601)
	Datetime                time.Time `json:"datetime" bson:"datetime"`                                 // gmt posted
	PersonalFinanceCategory string    `json:"personal_finance_category" bson:"personalFinanceCategory"` // most likely useless
	UID                     string    `bson:"uid" json:"uid"`
	IID                     string    `bson:"IID" json:"IID"`
}
