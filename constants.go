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
	EventClearAllWarnings
	EventUpdateSubscriptions
	EventUpdateChart
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

type PossibleSubscriptions struct {
	// The merchant name or transaction description.  If the `transactions` object was returned by an endpoint such as `/transactions/get`, this field will always appear. If the `transactions` object was returned by an endpoint such as `/asset_report/get/` or `/asset_report/pdf/get`, this field will only appear in an Asset Report with Insights.
	Name string `json:"name" bson:"name"`
	// The string returned by the financial institution to describe the transaction. For transactions returned by `/transactions/get`, this field is in beta and will be omitted unless the client is both enrolled in the closed beta program and has set `options.include_original_description` to `true`.
	OriginalDescription string  `json:"original_description" bson:"originalDescription"`
	Amount              float64 `json:"amount" bson:"amount"`
	// The merchant name, as extracted by Plaid from the `name` field.
	MerchantName        string `json:"merchant_name" bson:"merchantName"`
	UID                 string `bson:"uid" json:"uid"`
	IsPhysicalLocation  bool   `bson:"isPhysicalLocation" json:"is_physical_location"`
	FlatType            string `json:"flatType" bson:"flatType"`
	DetailedDescription string `bson:"detailedDescription" json:"detailed_description"`
}

type Categories struct {
	UID      string  `bson:"uid" json:"uid"`
	FlatType string  `json:"flatType" bson:"flatType"`
	Total    float64 `bson:"total" json:"total"`
}

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
	PhysicalLocation    bool        `json:"physical_location" bson:"physical_location"`
	TransactionType     XactionType `bson:"transaction_type" json:"transaction_type"`
	Description         string      `json:"description" bson:"description"`
	DetailedDescription string      `bson:"detailed_description" json:"detailed_description"`
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
	// these are introspected - plaid api is usually broken
	IsPhysicalLocation  bool        `bson:"isPhysicalLocation" json:"is_physical_location"`
	SCType              XactionType `json:"sc_type" bson:"SCType"`
	DetailedDescription string      `bson:"detailedDescription" json:"detailed_description"`
}

// DetailedClassify maps everything over
func DetailedClassify(transaction CowTransaction) TransactionMap {
	switch transaction.CategoryID {

	case "10000000":

		return TransactionMap{
			ID:                  "10000000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "bank fees",
			DetailedDescription: "bank fees",
		}
	case "10001000":

		return TransactionMap{
			ID:                  "10001000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "bank fees",
			DetailedDescription: "bank fees=>overdraft",
		}
	case "10002000":

		return TransactionMap{
			ID:                  "10002000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "bank fees",
			DetailedDescription: "bank fees=>atm",
		}
	case "10003000":

		return TransactionMap{
			ID:                  "10003000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "bank fees",
			DetailedDescription: "bank fees=>late payment",
		}
	case "10004000":

		return TransactionMap{
			ID:                  "10004000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "bank fees",
			DetailedDescription: "bank fees=>fraud dispute",
		}
	case "10005000":

		return TransactionMap{
			ID:                  "10005000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "bank fees",
			DetailedDescription: "bank fees=>foreign transaction",
		}
	case "10006000":

		return TransactionMap{
			ID:                  "10006000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "bank fees",
			DetailedDescription: "bank fees=>wire transfer",
		}
	case "10007000":

		return TransactionMap{
			ID:                  "10007000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "bank fees",
			DetailedDescription: "bank fees=>insufficient funds",
		}
	case "10008000":

		return TransactionMap{
			ID:                  "10008000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "bank fees",
			DetailedDescription: "bank fees=>cash advance",
		}
	case "10009000":

		return TransactionMap{
			ID:                  "10009000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "bank fees",
			DetailedDescription: "bank fees=>excess activity",
		}
	case "11000000":

		return TransactionMap{
			ID:                  "11000000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "cash advance",
			DetailedDescription: "cash advance",
		}
	case "12000000":

		return TransactionMap{
			ID:                  "12000000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community",
		}
	case "12001000":

		return TransactionMap{
			ID:                  "12001000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>animal shelter",
		}
	case "12002000":

		return TransactionMap{
			ID:                  "12002000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>assisted living services",
		}
	case "12002001":

		return TransactionMap{
			ID:                  "12002001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>assisted living services=>facilities and nursing homes",
		}
	case "12002002":

		return TransactionMap{
			ID:                  "12002002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>assisted living services=>caretakers",
		}
	case "12003000":

		return TransactionMap{
			ID:                  "12003000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>cemetery",
		}
	case "12004000":

		return TransactionMap{
			ID:                  "12004000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>courts",
		}
	case "12005000":

		return TransactionMap{
			ID:                  "12005000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>day care and preschools",
		}
	case "12006000":

		return TransactionMap{
			ID:                  "12006000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>disabled persons services",
		}
	case "12007000":

		return TransactionMap{
			ID:                  "12007000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>drug and alcohol services",
		}
	case "12008000":

		return TransactionMap{
			ID:                  "12008000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education",
		}
	case "12008001":

		return TransactionMap{
			ID:                  "12008001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>vocational schools",
		}
	case "12008002":

		return TransactionMap{
			ID:                  "12008002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>tutoring and educational services",
		}
	case "12008003":

		return TransactionMap{
			ID:                  "12008003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>primary and secondary schools",
		}
	case "12008004":

		return TransactionMap{
			ID:                  "12008004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>fraternities and sororities",
		}
	case "12008005":

		return TransactionMap{
			ID:                  "12008005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>driving schools",
		}
	case "12008006":

		return TransactionMap{
			ID:                  "12008006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>dance schools",
		}
	case "12008007":

		return TransactionMap{
			ID:                  "12008007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>culinary lessons and schools",
		}
	case "12008008":

		return TransactionMap{
			ID:                  "12008008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>computer training",
		}
	case "12008009":

		return TransactionMap{
			ID:                  "12008009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>colleges and universities",
		}
	case "12008010":

		return TransactionMap{
			ID:                  "12008010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>art school",
		}
	case "12008011":

		return TransactionMap{
			ID:                  "12008011",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>education=>adult education",
		}
	case "12009000":

		return TransactionMap{
			ID:                  "12009000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>government departments and agencies",
		}
	case "12010000":

		return TransactionMap{
			ID:                  "12010000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>government lobbyists",
		}
	case "12011000":

		return TransactionMap{
			ID:                  "12011000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>housing assistance and shelters",
		}
	case "12012000":

		return TransactionMap{
			ID:                  "12012000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>law enforcement",
		}
	case "12012001":

		return TransactionMap{
			ID:                  "12012001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>law enforcement=>police stations",
		}
	case "12012002":

		return TransactionMap{
			ID:                  "12012002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>law enforcement=>fire stations",
		}
	case "12012003":

		return TransactionMap{
			ID:                  "12012003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>law enforcement=>correctional institutions",
		}
	case "12013000":

		return TransactionMap{
			ID:                  "12013000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>libraries",
		}
	case "12014000":

		return TransactionMap{
			ID:                  "12014000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>military",
		}
	case "12015000":

		return TransactionMap{
			ID:                  "12015000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>organizations and associations",
		}
	case "12015001":

		return TransactionMap{
			ID:                  "12015001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>organizations and associations=>youth organizations",
		}
	case "12015002":

		return TransactionMap{
			ID:                  "12015002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>organizations and associations=>environmental",
		}
	case "12015003":

		return TransactionMap{
			ID:                  "12015003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>organizations and associations=>charities and non-profits",
		}
	case "12016000":

		return TransactionMap{
			ID:                  "12016000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>post offices",
		}
	case "12017000":

		return TransactionMap{
			ID:                  "12017000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>public and social services",
		}
	case "12018000":

		return TransactionMap{
			ID:                  "12018000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>religious",
		}
	case "12018001":

		return TransactionMap{
			ID:                  "12018001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>religious=>temple",
		}
	case "12018002":

		return TransactionMap{
			ID:                  "12018002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>religious=>synagogues",
		}
	case "12018003":

		return TransactionMap{
			ID:                  "12018003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>religious=>mosques",
		}
	case "12018004":

		return TransactionMap{
			ID:                  "12018004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>religious=>churches",
		}
	case "12019000":

		return TransactionMap{
			ID:                  "12019000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>senior citizen services",
		}
	case "12019001":

		return TransactionMap{
			ID:                  "12019001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "community",
			DetailedDescription: "community=>senior citizen services=>retirement",
		}
	case "13000000":

		return TransactionMap{
			ID:                  "13000000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink",
		}
	case "13001000":

		return TransactionMap{
			ID:                  "13001000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>bar",
		}
	case "13001001":

		return TransactionMap{
			ID:                  "13001001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>bar=>wine bar",
		}
	case "13001002":

		return TransactionMap{
			ID:                  "13001002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>bar=>sports bar",
		}
	case "13001003":

		return TransactionMap{
			ID:                  "13001003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>bar=>hotel lounge",
		}
	case "13002000":

		return TransactionMap{
			ID:                  "13002000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>breweries",
		}
	case "13003000":

		return TransactionMap{
			ID:                  "13003000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>internet cafes",
		}
	case "13004000":

		return TransactionMap{
			ID:                  "13004000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>nightlife",
		}
	case "13004001":

		return TransactionMap{
			ID:                  "13004001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>nightlife=>strip club",
		}
	case "13004002":

		return TransactionMap{
			ID:                  "13004002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>nightlife=>night clubs",
		}
	case "13004003":

		return TransactionMap{
			ID:                  "13004003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>nightlife=>karaoke",
		}
	case "13004004":

		return TransactionMap{
			ID:                  "13004004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>nightlife=>jazz and blues cafe",
		}
	case "13004005":

		return TransactionMap{
			ID:                  "13004005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>nightlife=>hookah lounges",
		}
	case "13004006":

		return TransactionMap{
			ID:                  "13004006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>nightlife=>adult entertainment",
		}
	case "13005000":

		return TransactionMap{
			ID:                  "13005000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants",
		}
	case "13005001":

		return TransactionMap{
			ID:                  "13005001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>winery",
		}
	case "13005002":

		return TransactionMap{
			ID:                  "13005002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>vegan and vegetarian",
		}
	case "13005003":

		return TransactionMap{
			ID:                  "13005003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>turkish",
		}
	case "13005004":

		return TransactionMap{
			ID:                  "13005004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>thai",
		}
	case "13005005":

		return TransactionMap{
			ID:                  "13005005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>swiss",
		}
	case "13005006":

		return TransactionMap{
			ID:                  "13005006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>sushi",
		}
	case "13005007":

		return TransactionMap{
			ID:                  "13005007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>steakhouses",
		}
	case "13005008":

		return TransactionMap{
			ID:                  "13005008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>spanish",
		}
	case "13005009":

		return TransactionMap{
			ID:                  "13005009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>seafood",
		}
	case "13005010":

		return TransactionMap{
			ID:                  "13005010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>scandinavian",
		}
	case "13005011":

		return TransactionMap{
			ID:                  "13005011",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>portuguese",
		}
	case "13005012":

		return TransactionMap{
			ID:                  "13005012",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>pizza",
		}
	case "13005013":

		return TransactionMap{
			ID:                  "13005013",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>moroccan",
		}
	case "13005014":

		return TransactionMap{
			ID:                  "13005014",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>middle eastern",
		}
	case "13005015":

		return TransactionMap{
			ID:                  "13005015",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>mexican",
		}
	case "13005016":

		return TransactionMap{
			ID:                  "13005016",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>mediterranean",
		}
	case "13005017":

		return TransactionMap{
			ID:                  "13005017",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>latin american",
		}
	case "13005018":

		return TransactionMap{
			ID:                  "13005018",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>korean",
		}
	case "13005019":

		return TransactionMap{
			ID:                  "13005019",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>juice bar",
		}
	case "13005020":

		return TransactionMap{
			ID:                  "13005020",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>japanese",
		}
	case "13005021":

		return TransactionMap{
			ID:                  "13005021",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>italian",
		}
	case "13005022":

		return TransactionMap{
			ID:                  "13005022",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>indonesian",
		}
	case "13005023":

		return TransactionMap{
			ID:                  "13005023",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>indian",
		}
	case "13005024":

		return TransactionMap{
			ID:                  "13005024",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>ice cream",
		}
	case "13005025":

		return TransactionMap{
			ID:                  "13005025",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>greek",
		}
	case "13005026":

		return TransactionMap{
			ID:                  "13005026",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>german",
		}
	case "13005027":

		return TransactionMap{
			ID:                  "13005027",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>gastropub",
		}
	case "13005028":

		return TransactionMap{
			ID:                  "13005028",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>french",
		}
	case "13005029":

		return TransactionMap{
			ID:                  "13005029",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>food truck",
		}
	case "13005030":

		return TransactionMap{
			ID:                  "13005030",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>fish and chips",
		}
	case "13005031":

		return TransactionMap{
			ID:                  "13005031",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>filipino",
		}
	case "13005032":

		return TransactionMap{
			ID:                  "13005032",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>fast food",
		}
	case "13005033":

		return TransactionMap{
			ID:                  "13005033",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>falafel",
		}
	case "13005034":

		return TransactionMap{
			ID:                  "13005034",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>ethiopian",
		}
	case "13005035":

		return TransactionMap{
			ID:                  "13005035",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>eastern european",
		}
	case "13005036":

		return TransactionMap{
			ID:                  "13005036",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>donuts",
		}
	case "13005037":

		return TransactionMap{
			ID:                  "13005037",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>distillery",
		}
	case "13005038":

		return TransactionMap{
			ID:                  "13005038",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>diners",
		}
	case "13005039":

		return TransactionMap{
			ID:                  "13005039",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>dessert",
		}
	case "13005040":

		return TransactionMap{
			ID:                  "13005040",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>delis",
		}
	case "13005041":

		return TransactionMap{
			ID:                  "13005041",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>cupcake shop",
		}
	case "13005042":

		return TransactionMap{
			ID:                  "13005042",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>cuban",
		}
	case "13005043":

		return TransactionMap{
			ID:                  "13005043",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>coffee shop",
		}
	case "13005044":

		return TransactionMap{
			ID:                  "13005044",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>chinese",
		}
	case "13005045":

		return TransactionMap{
			ID:                  "13005045",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>caribbean",
		}
	case "13005046":

		return TransactionMap{
			ID:                  "13005046",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>cajun",
		}
	case "13005047":

		return TransactionMap{
			ID:                  "13005047",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>cafe",
		}
	case "13005048":

		return TransactionMap{
			ID:                  "13005048",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>burrito",
		}
	case "13005049":

		return TransactionMap{
			ID:                  "13005049",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>burgers",
		}
	case "13005050":

		return TransactionMap{
			ID:                  "13005050",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>breakfast spot",
		}
	case "13005051":

		return TransactionMap{
			ID:                  "13005051",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>brazilian",
		}
	case "13005052":

		return TransactionMap{
			ID:                  "13005052",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>barbecue",
		}
	case "13005053":

		return TransactionMap{
			ID:                  "13005053",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>bakery",
		}
	case "13005054":

		return TransactionMap{
			ID:                  "13005054",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>bagel shop",
		}
	case "13005055":

		return TransactionMap{
			ID:                  "13005055",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>australian",
		}
	case "13005056":

		return TransactionMap{
			ID:                  "13005056",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>asian",
		}
	case "13005057":

		return TransactionMap{
			ID:                  "13005057",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>american",
		}
	case "13005058":

		return TransactionMap{
			ID:                  "13005058",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>african",
		}
	case "13005059":

		return TransactionMap{
			ID:                  "13005059",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "food and drink",
			DetailedDescription: "food and drink=>restaurants=>afghan",
		}
	case "14000000":

		return TransactionMap{
			ID:                  "14000000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare",
		}
	case "14001000":

		return TransactionMap{
			ID:                  "14001000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services",
		}
	case "14001001":

		return TransactionMap{
			ID:                  "14001001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>psychologists",
		}
	case "14001002":

		return TransactionMap{
			ID:                  "14001002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>pregnancy and sexual health",
		}
	case "14001003":

		return TransactionMap{
			ID:                  "14001003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>podiatrists",
		}
	case "14001004":

		return TransactionMap{
			ID:                  "14001004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>physical therapy",
		}
	case "14001005":

		return TransactionMap{
			ID:                  "14001005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>optometrists",
		}
	case "14001006":

		return TransactionMap{
			ID:                  "14001006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>nutritionists",
		}
	case "14001007":

		return TransactionMap{
			ID:                  "14001007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>nurses",
		}
	case "14001008":

		return TransactionMap{
			ID:                  "14001008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>mental health",
		}
	case "14001009":

		return TransactionMap{
			ID:                  "14001009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>medical supplies and labs",
		}
	case "14001010":

		return TransactionMap{
			ID:                  "14001010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>hospitals, clinics and medical centers",
		}
	case "14001011":

		return TransactionMap{
			ID:                  "14001011",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>emergency services",
		}
	case "14001012":

		return TransactionMap{
			ID:                  "14001012",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>dentists",
		}
	case "14001013":

		return TransactionMap{
			ID:                  "14001013",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>counseling and therapy",
		}
	case "14001014":

		return TransactionMap{
			ID:                  "14001014",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>chiropractors",
		}
	case "14001015":

		return TransactionMap{
			ID:                  "14001015",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>blood banks and centers",
		}
	case "14001016":

		return TransactionMap{
			ID:                  "14001016",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>alternative medicine",
		}
	case "14001017":

		return TransactionMap{
			ID:                  "14001017",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>healthcare services=>acupuncture",
		}
	case "14002000":

		return TransactionMap{
			ID:                  "14002000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians",
		}
	case "14002001":

		return TransactionMap{
			ID:                  "14002001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>urologists",
		}
	case "14002002":

		return TransactionMap{
			ID:                  "14002002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>respiratory",
		}
	case "14002003":

		return TransactionMap{
			ID:                  "14002003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>radiologists",
		}
	case "14002004":

		return TransactionMap{
			ID:                  "14002004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>psychiatrists",
		}
	case "14002005":

		return TransactionMap{
			ID:                  "14002005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>plastic surgeons",
		}
	case "14002006":

		return TransactionMap{
			ID:                  "14002006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>pediatricians",
		}
	case "14002007":

		return TransactionMap{
			ID:                  "14002007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>pathologists",
		}
	case "14002008":

		return TransactionMap{
			ID:                  "14002008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>orthopedic surgeons",
		}
	case "14002009":

		return TransactionMap{
			ID:                  "14002009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>ophthalmologists",
		}
	case "14002010":

		return TransactionMap{
			ID:                  "14002010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>oncologists",
		}
	case "14002011":

		return TransactionMap{
			ID:                  "14002011",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>obstetricians and gynecologists",
		}
	case "14002012":

		return TransactionMap{
			ID:                  "14002012",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>neurologists",
		}
	case "14002013":

		return TransactionMap{
			ID:                  "14002013",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>internal medicine",
		}
	case "14002014":

		return TransactionMap{
			ID:                  "14002014",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>general surgery",
		}
	case "14002015":

		return TransactionMap{
			ID:                  "14002015",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>gastroenterologists",
		}
	case "14002016":

		return TransactionMap{
			ID:                  "14002016",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>family medicine",
		}
	case "14002017":

		return TransactionMap{
			ID:                  "14002017",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>ear, nose and throat",
		}
	case "14002018":

		return TransactionMap{
			ID:                  "14002018",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>dermatologists",
		}
	case "14002019":

		return TransactionMap{
			ID:                  "14002019",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>cardiologists",
		}
	case "14002020":

		return TransactionMap{
			ID:                  "14002020",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "healthcare",
			DetailedDescription: "healthcare=>physicians=>anesthesiologists",
		}
	case "15000000":

		return TransactionMap{
			ID:                  "15000000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "interest",
			DetailedDescription: "interest",
		}
	case "15001000":

		return TransactionMap{
			ID:                  "15001000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "interest",
			DetailedDescription: "interest=>interest earned",
		}
	case "15002000":

		return TransactionMap{
			ID:                  "15002000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "interest",
			DetailedDescription: "interest=>interest charged",
		}
	case "16000000":

		return TransactionMap{
			ID:                  "16000000",
			PhysicalLocation:    false,
			TransactionType:     XactionPayment,
			Description:         "payment",
			DetailedDescription: "payment",
		}
	case "16001000":

		return TransactionMap{
			ID:                  "16001000",
			PhysicalLocation:    false,
			TransactionType:     XactionPayment,
			Description:         "payment",
			DetailedDescription: "payment=>credit card",
		}
	case "16002000":

		return TransactionMap{
			ID:                  "16002000",
			PhysicalLocation:    false,
			TransactionType:     XactionPayment,
			Description:         "payment",
			DetailedDescription: "payment=>rent",
		}
	case "16003000":

		return TransactionMap{
			ID:                  "16003000",
			PhysicalLocation:    false,
			TransactionType:     XactionPayment,
			Description:         "payment",
			DetailedDescription: "payment=>loan",
		}
	case "17000000":

		return TransactionMap{
			ID:                  "17000000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation",
		}
	case "17001000":

		return TransactionMap{
			ID:                  "17001000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment",
		}
	case "17001001":

		return TransactionMap{
			ID:                  "17001001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>theatrical productions",
		}
	case "17001002":

		return TransactionMap{
			ID:                  "17001002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>symphony and opera",
		}
	case "17001003":

		return TransactionMap{
			ID:                  "17001003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>sports venues",
		}
	case "17001004":

		return TransactionMap{
			ID:                  "17001004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>social clubs",
		}
	case "17001005":

		return TransactionMap{
			ID:                  "17001005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>psychics and astrologers",
		}
	case "17001006":

		return TransactionMap{
			ID:                  "17001006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>party centers",
		}
	case "17001007":

		return TransactionMap{
			ID:                  "17001007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>music and show venues",
		}
	case "17001008":

		return TransactionMap{
			ID:                  "17001008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>museums",
		}
	case "17001009":

		return TransactionMap{
			ID:                  "17001009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>movie theatres",
		}
	case "17001010":

		return TransactionMap{
			ID:                  "17001010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>fairgrounds and rodeos",
		}
	case "17001011":

		return TransactionMap{
			ID:                  "17001011",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>entertainment",
		}
	case "17001012":

		return TransactionMap{
			ID:                  "17001012",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>dance halls and saloons",
		}
	case "17001013":

		return TransactionMap{
			ID:                  "17001013",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>circuses and carnivals",
		}
	case "17001014":

		return TransactionMap{
			ID:                  "17001014",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>casinos and gaming",
		}
	case "17001015":

		return TransactionMap{
			ID:                  "17001015",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>bowling",
		}
	case "17001016":

		return TransactionMap{
			ID:                  "17001016",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>billiards and pool",
		}
	case "17001017":

		return TransactionMap{
			ID:                  "17001017",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>art dealers and galleries",
		}
	case "17001018":

		return TransactionMap{
			ID:                  "17001018",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>arcades and amusement parks",
		}
	case "17001019":

		return TransactionMap{
			ID:                  "17001019",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>arts and entertainment=>aquarium",
		}
	case "17002000":

		return TransactionMap{
			ID:                  "17002000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>athletic fields",
		}
	case "17003000":

		return TransactionMap{
			ID:                  "17003000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>baseball",
		}
	case "17004000":

		return TransactionMap{
			ID:                  "17004000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>basketball",
		}
	case "17005000":

		return TransactionMap{
			ID:                  "17005000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>batting cages",
		}
	case "17006000":

		return TransactionMap{
			ID:                  "17006000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>boating",
		}
	case "17007000":

		return TransactionMap{
			ID:                  "17007000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>campgrounds and rv parks",
		}
	case "17008000":

		return TransactionMap{
			ID:                  "17008000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>canoes and kayaks",
		}
	case "17009000":

		return TransactionMap{
			ID:                  "17009000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>combat sports",
		}
	case "17010000":

		return TransactionMap{
			ID:                  "17010000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>cycling",
		}
	case "17011000":

		return TransactionMap{
			ID:                  "17011000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>dance",
		}
	case "17012000":

		return TransactionMap{
			ID:                  "17012000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>equestrian",
		}
	case "17013000":

		return TransactionMap{
			ID:                  "17013000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>football",
		}
	case "17014000":

		return TransactionMap{
			ID:                  "17014000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>go carts",
		}
	case "17015000":

		return TransactionMap{
			ID:                  "17015000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>golf",
		}
	case "17016000":

		return TransactionMap{
			ID:                  "17016000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>gun ranges",
		}
	case "17017000":

		return TransactionMap{
			ID:                  "17017000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>gymnastics",
		}
	case "17018000":

		return TransactionMap{
			ID:                  "17018000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>gyms and fitness centers",
		}
	case "17019000":

		return TransactionMap{
			ID:                  "17019000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>hiking",
		}
	case "17020000":

		return TransactionMap{
			ID:                  "17020000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>hockey",
		}
	case "17021000":

		return TransactionMap{
			ID:                  "17021000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>hot air balloons",
		}
	case "17022000":

		return TransactionMap{
			ID:                  "17022000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>hunting and fishing",
		}
	case "17023000":

		return TransactionMap{
			ID:                  "17023000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>landmarks",
		}
	case "17023001":

		return TransactionMap{
			ID:                  "17023001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>landmarks=>monuments and memorials",
		}
	case "17023002":

		return TransactionMap{
			ID:                  "17023002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>landmarks=>historic sites",
		}
	case "17023003":

		return TransactionMap{
			ID:                  "17023003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>landmarks=>gardens",
		}
	case "17023004":

		return TransactionMap{
			ID:                  "17023004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>landmarks=>buildings and structures",
		}
	case "17024000":

		return TransactionMap{
			ID:                  "17024000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>miniature golf",
		}
	case "17025000":

		return TransactionMap{
			ID:                  "17025000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>outdoors",
		}
	case "17025001":

		return TransactionMap{
			ID:                  "17025001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>outdoors=>rivers",
		}
	case "17025002":

		return TransactionMap{
			ID:                  "17025002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>outdoors=>mountains",
		}
	case "17025003":

		return TransactionMap{
			ID:                  "17025003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>outdoors=>lakes",
		}
	case "17025004":

		return TransactionMap{
			ID:                  "17025004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>outdoors=>forests",
		}
	case "17025005":

		return TransactionMap{
			ID:                  "17025005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>outdoors=>beaches",
		}
	case "17026000":

		return TransactionMap{
			ID:                  "17026000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>paintball",
		}
	case "17027000":

		return TransactionMap{
			ID:                  "17027000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>parks",
		}
	case "17027001":

		return TransactionMap{
			ID:                  "17027001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>parks=>playgrounds",
		}
	case "17027002":

		return TransactionMap{
			ID:                  "17027002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>parks=>picnic areas",
		}
	case "17027003":

		return TransactionMap{
			ID:                  "17027003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>parks=>natural parks",
		}
	case "17028000":

		return TransactionMap{
			ID:                  "17028000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>personal trainers",
		}
	case "17029000":

		return TransactionMap{
			ID:                  "17029000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>race tracks",
		}
	case "17030000":

		return TransactionMap{
			ID:                  "17030000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>racquet sports",
		}
	case "17031000":

		return TransactionMap{
			ID:                  "17031000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>racquetball",
		}
	case "17032000":

		return TransactionMap{
			ID:                  "17032000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>rafting",
		}
	case "17033000":

		return TransactionMap{
			ID:                  "17033000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>recreation centers",
		}
	case "17034000":

		return TransactionMap{
			ID:                  "17034000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>rock climbing",
		}
	case "17035000":

		return TransactionMap{
			ID:                  "17035000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>running",
		}
	case "17036000":

		return TransactionMap{
			ID:                  "17036000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>scuba diving",
		}
	case "17037000":

		return TransactionMap{
			ID:                  "17037000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>skating",
		}
	case "17038000":

		return TransactionMap{
			ID:                  "17038000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>skydiving",
		}
	case "17039000":

		return TransactionMap{
			ID:                  "17039000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>snow sports",
		}
	case "17040000":

		return TransactionMap{
			ID:                  "17040000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>soccer",
		}
	case "17041000":

		return TransactionMap{
			ID:                  "17041000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>sports and recreation camps",
		}
	case "17042000":

		return TransactionMap{
			ID:                  "17042000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>sports clubs",
		}
	case "17043000":

		return TransactionMap{
			ID:                  "17043000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>stadiums and arenas",
		}
	case "17044000":

		return TransactionMap{
			ID:                  "17044000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>swimming",
		}
	case "17045000":

		return TransactionMap{
			ID:                  "17045000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>tennis",
		}
	case "17046000":

		return TransactionMap{
			ID:                  "17046000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>water sports",
		}
	case "17047000":

		return TransactionMap{
			ID:                  "17047000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>yoga and pilates",
		}
	case "17048000":

		return TransactionMap{
			ID:                  "17048000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "recreation",
			DetailedDescription: "recreation=>zoo",
		}
	case "18000000":

		return TransactionMap{
			ID:                  "18000000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service",
		}
	case "18001000":

		return TransactionMap{
			ID:                  "18001000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing",
		}
	case "18001001":

		return TransactionMap{
			ID:                  "18001001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing=>writing, copywriting and technical writing",
		}
	case "18001002":

		return TransactionMap{
			ID:                  "18001002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing=>search engine marketing and optimization",
		}
	case "18001003":

		return TransactionMap{
			ID:                  "18001003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing=>public relations",
		}
	case "18001004":

		return TransactionMap{
			ID:                  "18001004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing=>promotional items",
		}
	case "18001005":

		return TransactionMap{
			ID:                  "18001005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing=>print, tv, radio and outdoor advertising",
		}
	case "18001006":

		return TransactionMap{
			ID:                  "18001006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing=>online advertising",
		}
	case "18001007":

		return TransactionMap{
			ID:                  "18001007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing=>market research and consulting",
		}
	case "18001008":

		return TransactionMap{
			ID:                  "18001008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing=>direct mail and email marketing services",
		}
	case "18001009":

		return TransactionMap{
			ID:                  "18001009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing=>creative services",
		}
	case "18001010":

		return TransactionMap{
			ID:                  "18001010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>advertising and marketing=>advertising agencies and media buyers",
		}
	case "18003000":

		return TransactionMap{
			ID:                  "18003000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>art restoration",
		}
	case "18004000":

		return TransactionMap{
			ID:                  "18004000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>audiovisual",
		}
	case "18005000":

		return TransactionMap{
			ID:                  "18005000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automation and control systems",
		}
	case "18006000":

		return TransactionMap{
			ID:                  "18006000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automotive",
		}
	case "18006001":

		return TransactionMap{
			ID:                  "18006001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automotive=>towing",
		}
	case "18006002":

		return TransactionMap{
			ID:                  "18006002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automotive=>motorcycle, moped and scooter repair",
		}
	case "18006003":

		return TransactionMap{
			ID:                  "18006003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automotive=>maintenance and repair",
		}
	case "18006004":

		return TransactionMap{
			ID:                  "18006004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automotive=>car wash and detail",
		}
	case "18006005":

		return TransactionMap{
			ID:                  "18006005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automotive=>car appraisers",
		}
	case "18006006":

		return TransactionMap{
			ID:                  "18006006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automotive=>auto transmission",
		}
	case "18006007":

		return TransactionMap{
			ID:                  "18006007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automotive=>auto tires",
		}
	case "18006008":

		return TransactionMap{
			ID:                  "18006008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automotive=>auto smog check",
		}
	case "18006009":

		return TransactionMap{
			ID:                  "18006009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>automotive=>auto oil and lube",
		}
	case "18007000":

		return TransactionMap{
			ID:                  "18007000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>business and strategy consulting",
		}
	case "18008000":

		return TransactionMap{
			ID:                  "18008000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>business services",
		}
	case "18008001":

		return TransactionMap{
			ID:                  "18008001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>business services=>printing and publishing",
		}
	case "18009000":

		return TransactionMap{
			ID:                  "18009000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>cable",
		}
	case "18010000":

		return TransactionMap{
			ID:                  "18010000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>chemicals and gasses",
		}
	case "18011000":

		return TransactionMap{
			ID:                  "18011000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>cleaning",
		}
	case "18012000":

		return TransactionMap{
			ID:                  "18012000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>computers",
		}
	case "18012001":

		return TransactionMap{
			ID:                  "18012001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>computers=>maintenance and repair",
		}
	case "18012002":

		return TransactionMap{
			ID:                  "18012002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>computers=>software development",
		}
	case "18013000":

		return TransactionMap{
			ID:                  "18013000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction",
		}
	case "18013001":

		return TransactionMap{
			ID:                  "18013001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction=>specialty",
		}
	case "18013002":

		return TransactionMap{
			ID:                  "18013002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction=>roofers",
		}
	case "18013003":

		return TransactionMap{
			ID:                  "18013003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction=>painting",
		}
	case "18013004":

		return TransactionMap{
			ID:                  "18013004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction=>masonry",
		}
	case "18013005":

		return TransactionMap{
			ID:                  "18013005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction=>infrastructure",
		}
	case "18013006":

		return TransactionMap{
			ID:                  "18013006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction=>heating, ventilating and air conditioning",
		}
	case "18013007":

		return TransactionMap{
			ID:                  "18013007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction=>electricians",
		}
	case "18013008":

		return TransactionMap{
			ID:                  "18013008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction=>contractors",
		}
	case "18013009":

		return TransactionMap{
			ID:                  "18013009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction=>carpet and flooring",
		}
	case "18013010":

		return TransactionMap{
			ID:                  "18013010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>construction=>carpenters",
		}
	case "18014000":

		return TransactionMap{
			ID:                  "18014000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>credit counseling and bankruptcy services",
		}
	case "18015000":

		return TransactionMap{
			ID:                  "18015000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>dating and escort",
		}
	case "18016000":

		return TransactionMap{
			ID:                  "18016000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>employment agencies",
		}
	case "18017000":

		return TransactionMap{
			ID:                  "18017000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>engineering",
		}
	case "18018000":

		return TransactionMap{
			ID:                  "18018000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>entertainment",
		}
	case "18018001":

		return TransactionMap{
			ID:                  "18018001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>entertainment=>media",
		}
	case "18019000":

		return TransactionMap{
			ID:                  "18019000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>events and event planning",
		}
	case "18020000":

		return TransactionMap{
			ID:                  "18020000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial",
		}
	case "18020001":

		return TransactionMap{
			ID:                  "18020001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>taxes",
		}
	case "18020002":

		return TransactionMap{
			ID:                  "18020002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>student aid and grants",
		}
	case "18020003":

		return TransactionMap{
			ID:                  "18020003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>stock brokers",
		}
	case "18020004":

		return TransactionMap{
			ID:                  "18020004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>loans and mortgages",
		}
	case "18020005":

		return TransactionMap{
			ID:                  "18020005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>holding and investment offices",
		}
	case "18020006":

		return TransactionMap{
			ID:                  "18020006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>fund raising",
		}
	case "18020007":

		return TransactionMap{
			ID:                  "18020007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>financial planning and investments",
		}
	case "18020008":

		return TransactionMap{
			ID:                  "18020008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>credit reporting",
		}
	case "18020009":

		return TransactionMap{
			ID:                  "18020009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>collections",
		}
	case "18020010":

		return TransactionMap{
			ID:                  "18020010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>check cashing",
		}
	case "18020011":

		return TransactionMap{
			ID:                  "18020011",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>business brokers and franchises",
		}
	case "18020012":

		return TransactionMap{
			ID:                  "18020012",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>banking and finance",
		}
	case "18020013":

		return TransactionMap{
			ID:                  "18020013",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>atms",
		}
	case "18020014":

		return TransactionMap{
			ID:                  "18020014",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>financial=>accounting and bookkeeping",
		}
	case "18021000":

		return TransactionMap{
			ID:                  "18021000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>food and beverage",
		}
	case "18021001":

		return TransactionMap{
			ID:                  "18021001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>food and beverage=>distribution",
		}
	case "18021002":

		return TransactionMap{
			ID:                  "18021002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>food and beverage=>catering",
		}
	case "18022000":

		return TransactionMap{
			ID:                  "18022000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>funeral services",
		}
	case "18023000":

		return TransactionMap{
			ID:                  "18023000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>geological",
		}
	case "18024000":

		return TransactionMap{
			ID:                  "18024000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement",
		}
	case "18024001":

		return TransactionMap{
			ID:                  "18024001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>upholstery",
		}
	case "18024002":

		return TransactionMap{
			ID:                  "18024002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>tree service",
		}
	case "18024003":

		return TransactionMap{
			ID:                  "18024003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>swimming pool maintenance and services",
		}
	case "18024004":

		return TransactionMap{
			ID:                  "18024004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>storage",
		}
	case "18024005":

		return TransactionMap{
			ID:                  "18024005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>roofers",
		}
	case "18024006":

		return TransactionMap{
			ID:                  "18024006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>pools and spas",
		}
	case "18024007":

		return TransactionMap{
			ID:                  "18024007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>plumbing",
		}
	case "18024008":

		return TransactionMap{
			ID:                  "18024008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>pest control",
		}
	case "18024009":

		return TransactionMap{
			ID:                  "18024009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>painting",
		}
	case "18024010":

		return TransactionMap{
			ID:                  "18024010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>movers",
		}
	case "18024011":

		return TransactionMap{
			ID:                  "18024011",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>mobile homes",
		}
	case "18024012":

		return TransactionMap{
			ID:                  "18024012",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>lighting fixtures",
		}
	case "18024013":

		return TransactionMap{
			ID:                  "18024013",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>landscaping and gardeners",
		}
	case "18024014":

		return TransactionMap{
			ID:                  "18024014",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>kitchens",
		}
	case "18024015":

		return TransactionMap{
			ID:                  "18024015",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>interior design",
		}
	case "18024016":

		return TransactionMap{
			ID:                  "18024016",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>housewares",
		}
	case "18024017":

		return TransactionMap{
			ID:                  "18024017",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>home inspection services",
		}
	case "18024018":

		return TransactionMap{
			ID:                  "18024018",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>home appliances",
		}
	case "18024019":

		return TransactionMap{
			ID:                  "18024019",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>heating, ventilation and air conditioning",
		}
	case "18024020":

		return TransactionMap{
			ID:                  "18024020",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>hardware and services",
		}
	case "18024021":

		return TransactionMap{
			ID:                  "18024021",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>fences, fireplaces and garage doors",
		}
	case "18024022":

		return TransactionMap{
			ID:                  "18024022",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>electricians",
		}
	case "18024023":

		return TransactionMap{
			ID:                  "18024023",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>doors and windows",
		}
	case "18024024":

		return TransactionMap{
			ID:                  "18024024",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>contractors",
		}
	case "18024025":

		return TransactionMap{
			ID:                  "18024025",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>carpet and flooring",
		}
	case "18024026":

		return TransactionMap{
			ID:                  "18024026",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>carpenters",
		}
	case "18024027":

		return TransactionMap{
			ID:                  "18024027",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>home improvement=>architects",
		}
	case "18025000":

		return TransactionMap{
			ID:                  "18025000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>household",
		}
	case "18026000":

		return TransactionMap{
			ID:                  "18026000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>human resources",
		}
	case "18027000":

		return TransactionMap{
			ID:                  "18027000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>immigration",
		}
	case "18028000":

		return TransactionMap{
			ID:                  "18028000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>import and export",
		}
	case "18029000":

		return TransactionMap{
			ID:                  "18029000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>industrial machinery and vehicles",
		}
	case "18030000":

		return TransactionMap{
			ID:                  "18030000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>insurance",
		}
	case "18031000":

		return TransactionMap{
			ID:                  "18031000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>internet services",
		}
	case "18032000":

		return TransactionMap{
			ID:                  "18032000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>leather",
		}
	case "18033000":

		return TransactionMap{
			ID:                  "18033000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>legal",
		}
	case "18034000":

		return TransactionMap{
			ID:                  "18034000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>logging and sawmills",
		}
	case "18035000":

		return TransactionMap{
			ID:                  "18035000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>machine shops",
		}
	case "18036000":

		return TransactionMap{
			ID:                  "18036000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>management",
		}
	case "18037000":

		return TransactionMap{
			ID:                  "18037000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing",
		}
	case "18037001":

		return TransactionMap{
			ID:                  "18037001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>apparel and fabric products",
		}
	case "18037002":

		return TransactionMap{
			ID:                  "18037002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>chemicals and gasses",
		}
	case "18037003":

		return TransactionMap{
			ID:                  "18037003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>computers and office machines",
		}
	case "18037004":

		return TransactionMap{
			ID:                  "18037004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>electrical equipment and components",
		}
	case "18037005":

		return TransactionMap{
			ID:                  "18037005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>food and beverage",
		}
	case "18037006":

		return TransactionMap{
			ID:                  "18037006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>furniture and fixtures",
		}
	case "18037007":

		return TransactionMap{
			ID:                  "18037007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>glass products",
		}
	case "18037008":

		return TransactionMap{
			ID:                  "18037008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>industrial machinery and equipment",
		}
	case "18037009":

		return TransactionMap{
			ID:                  "18037009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>leather goods",
		}
	case "18037010":

		return TransactionMap{
			ID:                  "18037010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>metal products",
		}
	case "18037011":

		return TransactionMap{
			ID:                  "18037011",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>nonmetallic mineral products",
		}
	case "18037012":

		return TransactionMap{
			ID:                  "18037012",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>paper products",
		}
	case "18037013":

		return TransactionMap{
			ID:                  "18037013",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>petroleum",
		}
	case "18037014":

		return TransactionMap{
			ID:                  "18037014",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>plastic products",
		}
	case "18037015":

		return TransactionMap{
			ID:                  "18037015",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>rubber products",
		}
	case "18037016":

		return TransactionMap{
			ID:                  "18037016",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>service instruments",
		}
	case "18037017":

		return TransactionMap{
			ID:                  "18037017",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>textiles",
		}
	case "18037018":

		return TransactionMap{
			ID:                  "18037018",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>tobacco",
		}
	case "18037019":

		return TransactionMap{
			ID:                  "18037019",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>transportation equipment",
		}
	case "18037020":

		return TransactionMap{
			ID:                  "18037020",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>manufacturing=>wood products",
		}
	case "18038000":

		return TransactionMap{
			ID:                  "18038000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>media production",
		}
	case "18039000":

		return TransactionMap{
			ID:                  "18039000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>metals",
		}
	case "18040000":

		return TransactionMap{
			ID:                  "18040000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>mining",
		}
	case "18040001":

		return TransactionMap{
			ID:                  "18040001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>mining=>coal",
		}
	case "18040002":

		return TransactionMap{
			ID:                  "18040002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>mining=>metal",
		}
	case "18040003":

		return TransactionMap{
			ID:                  "18040003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>mining=>non-metallic minerals",
		}
	case "18041000":

		return TransactionMap{
			ID:                  "18041000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>news reporting",
		}
	case "18042000":

		return TransactionMap{
			ID:                  "18042000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>oil and gas",
		}
	case "18043000":

		return TransactionMap{
			ID:                  "18043000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>packaging",
		}
	case "18044000":

		return TransactionMap{
			ID:                  "18044000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>paper",
		}
	case "18045000":

		return TransactionMap{
			ID:                  "18045000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care",
		}
	case "18045001":

		return TransactionMap{
			ID:                  "18045001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care=>tattooing",
		}
	case "18045002":

		return TransactionMap{
			ID:                  "18045002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care=>tanning salons",
		}
	case "18045003":

		return TransactionMap{
			ID:                  "18045003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care=>spas",
		}
	case "18045004":

		return TransactionMap{
			ID:                  "18045004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care=>skin care",
		}
	case "18045005":

		return TransactionMap{
			ID:                  "18045005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care=>piercing",
		}
	case "18045006":

		return TransactionMap{
			ID:                  "18045006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care=>massage clinics and therapists",
		}
	case "18045007":

		return TransactionMap{
			ID:                  "18045007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care=>manicures and pedicures",
		}
	case "18045008":

		return TransactionMap{
			ID:                  "18045008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care=>laundry and garment services",
		}
	case "18045009":

		return TransactionMap{
			ID:                  "18045009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care=>hair salons and barbers",
		}
	case "18045010":

		return TransactionMap{
			ID:                  "18045010",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>personal care=>hair removal",
		}
	case "18046000":

		return TransactionMap{
			ID:                  "18046000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>petroleum",
		}
	case "18047000":

		return TransactionMap{
			ID:                  "18047000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>photography",
		}
	case "18048000":

		return TransactionMap{
			ID:                  "18048000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>plastics",
		}
	case "18049000":

		return TransactionMap{
			ID:                  "18049000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>rail",
		}
	case "18050000":

		return TransactionMap{
			ID:                  "18050000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate",
		}
	case "18050001":

		return TransactionMap{
			ID:                  "18050001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate=>real estate development and title companies",
		}
	case "18050002":

		return TransactionMap{
			ID:                  "18050002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate=>real estate appraiser",
		}
	case "18050003":

		return TransactionMap{
			ID:                  "18050003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate=>real estate agents",
		}
	case "18050004":

		return TransactionMap{
			ID:                  "18050004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate=>property management",
		}
	case "18050005":

		return TransactionMap{
			ID:                  "18050005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate=>corporate housing",
		}
	case "18050006":

		return TransactionMap{
			ID:                  "18050006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate=>commercial real estate",
		}
	case "18050007":

		return TransactionMap{
			ID:                  "18050007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate=>building and land surveyors",
		}
	case "18050008":

		return TransactionMap{
			ID:                  "18050008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate=>boarding houses",
		}
	case "18050009":

		return TransactionMap{
			ID:                  "18050009",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate=>apartments, condos and houses",
		}
	case "18050010":

		return TransactionMap{
			ID:                  "18050010",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>real estate=>rent",
		}
	case "18051000":

		return TransactionMap{
			ID:                  "18051000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>refrigeration and ice",
		}
	case "18052000":

		return TransactionMap{
			ID:                  "18052000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>renewable energy",
		}
	case "18053000":

		return TransactionMap{
			ID:                  "18053000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>repair services",
		}
	case "18054000":

		return TransactionMap{
			ID:                  "18054000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>research",
		}
	case "18055000":

		return TransactionMap{
			ID:                  "18055000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>rubber",
		}
	case "18056000":

		return TransactionMap{
			ID:                  "18056000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>scientific",
		}
	case "18057000":

		return TransactionMap{
			ID:                  "18057000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>security and safety",
		}
	case "18058000":

		return TransactionMap{
			ID:                  "18058000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>shipping and freight",
		}
	case "18059000":

		return TransactionMap{
			ID:                  "18059000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>software development",
		}
	case "18060000":

		return TransactionMap{
			ID:                  "18060000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>storage",
		}
	case "18061000":

		return TransactionMap{
			ID:                  "18061000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>subscription",
		}
	case "18062000":

		return TransactionMap{
			ID:                  "18062000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>tailors",
		}
	case "18063000":

		return TransactionMap{
			ID:                  "18063000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>telecommunication services",
		}
	case "18064000":

		return TransactionMap{
			ID:                  "18064000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>textiles",
		}
	case "18065000":

		return TransactionMap{
			ID:                  "18065000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>tourist information and services",
		}
	case "18066000":

		return TransactionMap{
			ID:                  "18066000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>transportation",
		}
	case "18067000":

		return TransactionMap{
			ID:                  "18067000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>travel agents and tour operators",
		}
	case "18068000":

		return TransactionMap{
			ID:                  "18068000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>utilities",
		}
	case "18068001":

		return TransactionMap{
			ID:                  "18068001",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>utilities=>water",
		}
	case "18068002":

		return TransactionMap{
			ID:                  "18068002",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>utilities=>sanitary and waste management",
		}
	case "18068003":

		return TransactionMap{
			ID:                  "18068003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>utilities=>heating, ventilating, and air conditioning",
		}
	case "18068004":

		return TransactionMap{
			ID:                  "18068004",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>utilities=>gas",
		}
	case "18068005":

		return TransactionMap{
			ID:                  "18068005",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>utilities=>electric",
		}
	case "18069000":

		return TransactionMap{
			ID:                  "18069000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>veterinarians",
		}
	case "18070000":

		return TransactionMap{
			ID:                  "18070000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>water and waste management",
		}
	case "18071000":

		return TransactionMap{
			ID:                  "18071000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>web design and development",
		}
	case "18072000":

		return TransactionMap{
			ID:                  "18072000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>welding",
		}
	case "18073000":

		return TransactionMap{
			ID:                  "18073000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>agriculture and forestry",
		}
	case "18073001":

		return TransactionMap{
			ID:                  "18073001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>agriculture and forestry=>crop production",
		}
	case "18073002":

		return TransactionMap{
			ID:                  "18073002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>agriculture and forestry=>forestry",
		}
	case "18073003":

		return TransactionMap{
			ID:                  "18073003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>agriculture and forestry=>livestock and animals",
		}
	case "18073004":

		return TransactionMap{
			ID:                  "18073004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>agriculture and forestry=>services",
		}
	case "18074000":

		return TransactionMap{
			ID:                  "18074000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "service",
			DetailedDescription: "service=>art and graphic design",
		}
	case "19000000":

		return TransactionMap{
			ID:                  "19000000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops",
		}
	case "19001000":

		return TransactionMap{
			ID:                  "19001000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>adult",
		}
	case "19002000":

		return TransactionMap{
			ID:                  "19002000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>antiques",
		}
	case "19003000":

		return TransactionMap{
			ID:                  "19003000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>arts and crafts",
		}
	case "19004000":

		return TransactionMap{
			ID:                  "19004000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>auctions",
		}
	case "19005000":

		return TransactionMap{
			ID:                  "19005000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>automotive",
		}
	case "19005001":

		return TransactionMap{
			ID:                  "19005001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>automotive=>used car dealers",
		}
	case "19005002":

		return TransactionMap{
			ID:                  "19005002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>automotive=>salvage yards",
		}
	case "19005003":

		return TransactionMap{
			ID:                  "19005003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>automotive=>rvs and motor homes",
		}
	case "19005004":

		return TransactionMap{
			ID:                  "19005004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>automotive=>motorcycles, mopeds and scooters",
		}
	case "19005005":

		return TransactionMap{
			ID:                  "19005005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>automotive=>classic and antique car",
		}
	case "19005006":

		return TransactionMap{
			ID:                  "19005006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>automotive=>car parts and accessories",
		}
	case "19005007":

		return TransactionMap{
			ID:                  "19005007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>automotive=>car dealers and leasing",
		}
	case "19006000":

		return TransactionMap{
			ID:                  "19006000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>beauty products",
		}
	case "19007000":

		return TransactionMap{
			ID:                  "19007000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>bicycles",
		}
	case "19008000":

		return TransactionMap{
			ID:                  "19008000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>boat dealers",
		}
	case "19009000":

		return TransactionMap{
			ID:                  "19009000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>bookstores",
		}
	case "19010000":

		return TransactionMap{
			ID:                  "19010000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>cards and stationery",
		}
	case "19011000":

		return TransactionMap{
			ID:                  "19011000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>children",
		}
	case "19012000":

		return TransactionMap{
			ID:                  "19012000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>clothing and accessories",
		}
	case "19012001":

		return TransactionMap{
			ID:                  "19012001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>clothing and accessories=>women's store",
		}
	case "19012002":

		return TransactionMap{
			ID:                  "19012002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>clothing and accessories=>swimwear",
		}
	case "19012003":

		return TransactionMap{
			ID:                  "19012003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>clothing and accessories=>shoe store",
		}
	case "19012004":

		return TransactionMap{
			ID:                  "19012004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>clothing and accessories=>men's store",
		}
	case "19012005":

		return TransactionMap{
			ID:                  "19012005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>clothing and accessories=>lingerie store",
		}
	case "19012006":

		return TransactionMap{
			ID:                  "19012006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>clothing and accessories=>kids' store",
		}
	case "19012007":

		return TransactionMap{
			ID:                  "19012007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>clothing and accessories=>boutique",
		}
	case "19012008":

		return TransactionMap{
			ID:                  "19012008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>clothing and accessories=>accessories store",
		}
	case "19013000":

		return TransactionMap{
			ID:                  "19013000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>computers and electronics",
		}
	case "19013001":

		return TransactionMap{
			ID:                  "19013001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>computers and electronics=>video games",
		}
	case "19013002":

		return TransactionMap{
			ID:                  "19013002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>computers and electronics=>mobile phones",
		}
	case "19013003":

		return TransactionMap{
			ID:                  "19013003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>computers and electronics=>cameras",
		}
	case "19014000":

		return TransactionMap{
			ID:                  "19014000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>construction supplies",
		}
	case "19015000":

		return TransactionMap{
			ID:                  "19015000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>convenience stores",
		}
	case "19016000":

		return TransactionMap{
			ID:                  "19016000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>costumes",
		}
	case "19017000":

		return TransactionMap{
			ID:                  "19017000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>dance and music",
		}
	case "19018000":

		return TransactionMap{
			ID:                  "19018000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>department stores",
		}
	case "19019000":

		return TransactionMap{
			ID:                  "19019000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>digital purchase",
		}
	case "19020000":

		return TransactionMap{
			ID:                  "19020000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>discount stores",
		}
	case "19021000":

		return TransactionMap{
			ID:                  "19021000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>electrical equipment",
		}
	case "19022000":

		return TransactionMap{
			ID:                  "19022000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>equipment rental",
		}
	case "19023000":

		return TransactionMap{
			ID:                  "19023000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>flea markets",
		}
	case "19024000":

		return TransactionMap{
			ID:                  "19024000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>florists",
		}
	case "19025000":

		return TransactionMap{
			ID:                  "19025000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>food and beverage store",
		}
	case "19025001":

		return TransactionMap{
			ID:                  "19025001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>food and beverage store=>specialty",
		}
	case "19025002":

		return TransactionMap{
			ID:                  "19025002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>food and beverage store=>health food",
		}
	case "19025003":

		return TransactionMap{
			ID:                  "19025003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>food and beverage store=>farmers markets",
		}
	case "19025004":

		return TransactionMap{
			ID:                  "19025004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>food and beverage store=>beer, wine and spirits",
		}
	case "19026000":

		return TransactionMap{
			ID:                  "19026000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>fuel dealer",
		}
	case "19027000":

		return TransactionMap{
			ID:                  "19027000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>furniture and home decor",
		}
	case "19028000":

		return TransactionMap{
			ID:                  "19028000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>gift and novelty",
		}
	case "19029000":

		return TransactionMap{
			ID:                  "19029000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>glasses and optometrist",
		}
	case "19030000":

		return TransactionMap{
			ID:                  "19030000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>hardware store",
		}
	case "19031000":

		return TransactionMap{
			ID:                  "19031000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>hobby and collectibles",
		}
	case "19032000":

		return TransactionMap{
			ID:                  "19032000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>industrial supplies",
		}
	case "19033000":

		return TransactionMap{
			ID:                  "19033000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>jewelry and watches",
		}
	case "19034000":

		return TransactionMap{
			ID:                  "19034000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>luggage",
		}
	case "19035000":

		return TransactionMap{
			ID:                  "19035000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>marine supplies",
		}
	case "19036000":

		return TransactionMap{
			ID:                  "19036000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>music, video and dvd",
		}
	case "19037000":

		return TransactionMap{
			ID:                  "19037000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>musical instruments",
		}
	case "19038000":

		return TransactionMap{
			ID:                  "19038000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>newsstands",
		}
	case "19039000":

		return TransactionMap{
			ID:                  "19039000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>office supplies",
		}
	case "19040000":

		return TransactionMap{
			ID:                  "19040000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>outlet",
		}
	case "19040001":

		return TransactionMap{
			ID:                  "19040001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>outlet=>women's store",
		}
	case "19040002":

		return TransactionMap{
			ID:                  "19040002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>outlet=>swimwear",
		}
	case "19040003":

		return TransactionMap{
			ID:                  "19040003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>outlet=>shoe store",
		}
	case "19040004":

		return TransactionMap{
			ID:                  "19040004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>outlet=>men's store",
		}
	case "19040005":

		return TransactionMap{
			ID:                  "19040005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>outlet=>lingerie store",
		}
	case "19040006":

		return TransactionMap{
			ID:                  "19040006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>outlet=>kids' store",
		}
	case "19040007":

		return TransactionMap{
			ID:                  "19040007",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>outlet=>boutique",
		}
	case "19040008":

		return TransactionMap{
			ID:                  "19040008",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>outlet=>accessories store",
		}
	case "19041000":

		return TransactionMap{
			ID:                  "19041000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>pawn shops",
		}
	case "19042000":

		return TransactionMap{
			ID:                  "19042000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>pets",
		}
	case "19043000":

		return TransactionMap{
			ID:                  "19043000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>pharmacies",
		}
	case "19044000":

		return TransactionMap{
			ID:                  "19044000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>photos and frames",
		}
	case "19045000":

		return TransactionMap{
			ID:                  "19045000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>shopping centers and malls",
		}
	case "19046000":

		return TransactionMap{
			ID:                  "19046000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>sporting goods",
		}
	case "19047000":

		return TransactionMap{
			ID:                  "19047000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>supermarkets and groceries",
		}
	case "19048000":

		return TransactionMap{
			ID:                  "19048000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>tobacco",
		}
	case "19049000":

		return TransactionMap{
			ID:                  "19049000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>toys",
		}
	case "19050000":

		return TransactionMap{
			ID:                  "19050000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>vintage and thrift",
		}
	case "19051000":

		return TransactionMap{
			ID:                  "19051000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>warehouses and wholesale stores",
		}
	case "19052000":

		return TransactionMap{
			ID:                  "19052000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>wedding and bridal",
		}
	case "19053000":

		return TransactionMap{
			ID:                  "19053000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>wholesale",
		}
	case "19054000":

		return TransactionMap{
			ID:                  "19054000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "shops",
			DetailedDescription: "shops=>lawn and garden",
		}
	case "20000000":

		return TransactionMap{
			ID:                  "20000000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "tax",
			DetailedDescription: "tax",
		}
	case "20001000":

		return TransactionMap{
			ID:                  "20001000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "tax",
			DetailedDescription: "tax=>refund",
		}
	case "20002000":

		return TransactionMap{
			ID:                  "20002000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "tax",
			DetailedDescription: "tax=>payment",
		}
	case "21000000":

		return TransactionMap{
			ID:                  "21000000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer",
		}
	case "21001000":

		return TransactionMap{
			ID:                  "21001000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>internal account transfer",
		}
	case "21002000":

		return TransactionMap{
			ID:                  "21002000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>ach",
		}
	case "21003000":

		return TransactionMap{
			ID:                  "21003000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>billpay",
		}
	case "21004000":

		return TransactionMap{
			ID:                  "21004000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>check",
		}
	case "21005000":

		return TransactionMap{
			ID:                  "21005000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>credit",
		}
	case "21006000":

		return TransactionMap{
			ID:                  "21006000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>debit",
		}
	case "21007000":

		return TransactionMap{
			ID:                  "21007000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>deposit",
		}
	case "21007001":

		return TransactionMap{
			ID:                  "21007001",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>deposit=>check",
		}
	case "21007002":

		return TransactionMap{
			ID:                  "21007002",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>deposit=>atm",
		}
	case "21008000":

		return TransactionMap{
			ID:                  "21008000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>keep the change savings program",
		}
	case "21009000":

		return TransactionMap{
			ID:                  "21009000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>payroll",
		}
	case "21009001":

		return TransactionMap{
			ID:                  "21009001",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>payroll=>benefits",
		}
	case "21010000":

		return TransactionMap{
			ID:                  "21010000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party",
		}
	case "21010001":

		return TransactionMap{
			ID:                  "21010001",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>venmo",
		}
	case "21010002":

		return TransactionMap{
			ID:                  "21010002",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>square cash",
		}
	case "21010003":

		return TransactionMap{
			ID:                  "21010003",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>square",
		}
	case "21010004":

		return TransactionMap{
			ID:                  "21010004",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>paypal",
		}
	case "21010005":

		return TransactionMap{
			ID:                  "21010005",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>dwolla",
		}
	case "21010006":

		return TransactionMap{
			ID:                  "21010006",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>coinbase",
		}
	case "21010007":

		return TransactionMap{
			ID:                  "21010007",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>chase quickpay",
		}
	case "21010008":

		return TransactionMap{
			ID:                  "21010008",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>acorns",
		}
	case "21010009":

		return TransactionMap{
			ID:                  "21010009",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>digit",
		}
	case "21010010":

		return TransactionMap{
			ID:                  "21010010",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>betterment",
		}
	case "21010011":

		return TransactionMap{
			ID:                  "21010011",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>third party=>plaid",
		}
	case "21011000":

		return TransactionMap{
			ID:                  "21011000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>wire",
		}
	case "21012000":

		return TransactionMap{
			ID:                  "21012000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>withdrawal",
		}
	case "21012001":

		return TransactionMap{
			ID:                  "21012001",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>withdrawal=>check",
		}
	case "21012002":

		return TransactionMap{
			ID:                  "21012002",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>withdrawal=>atm",
		}
	case "21013000":

		return TransactionMap{
			ID:                  "21013000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "transfer",
			DetailedDescription: "transfer=>save as you go",
		}
	case "22000000":

		return TransactionMap{
			ID:                  "22000000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel",
		}
	case "22001000":

		return TransactionMap{
			ID:                  "22001000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>airlines and aviation services",
		}
	case "22002000":

		return TransactionMap{
			ID:                  "22002000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>airports",
		}
	case "22003000":

		return TransactionMap{
			ID:                  "22003000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>boat",
		}
	case "22004000":

		return TransactionMap{
			ID:                  "22004000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>bus stations",
		}
	case "22005000":

		return TransactionMap{
			ID:                  "22005000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>car and truck rentals",
		}
	case "22006000":

		return TransactionMap{
			ID:                  "22006000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>car service",
		}
	case "22006001":

		return TransactionMap{
			ID:                  "22006001",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>car service=>ride share",
		}
	case "22007000":

		return TransactionMap{
			ID:                  "22007000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>charter buses",
		}
	case "22008000":

		return TransactionMap{
			ID:                  "22008000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>cruises",
		}
	case "22009000":

		return TransactionMap{
			ID:                  "22009000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>gas stations",
		}
	case "22010000":

		return TransactionMap{
			ID:                  "22010000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>heliports",
		}
	case "22011000":

		return TransactionMap{
			ID:                  "22011000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>limos and chauffeurs",
		}
	case "22012000":

		return TransactionMap{
			ID:                  "22012000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>lodging",
		}
	case "22012001":

		return TransactionMap{
			ID:                  "22012001",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>lodging=>resorts",
		}
	case "22012002":

		return TransactionMap{
			ID:                  "22012002",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>lodging=>lodges and vacation rentals",
		}
	case "22012003":

		return TransactionMap{
			ID:                  "22012003",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>lodging=>hotels and motels",
		}
	case "22012004":

		return TransactionMap{
			ID:                  "22012004",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>lodging=>hostels",
		}
	case "22012005":

		return TransactionMap{
			ID:                  "22012005",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>lodging=>cottages and cabins",
		}
	case "22012006":

		return TransactionMap{
			ID:                  "22012006",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>lodging=>bed and breakfasts",
		}
	case "22013000":

		return TransactionMap{
			ID:                  "22013000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>parking",
		}
	case "22014000":

		return TransactionMap{
			ID:                  "22014000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>public transportation services",
		}
	case "22015000":

		return TransactionMap{
			ID:                  "22015000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>rail",
		}
	case "22016000":

		return TransactionMap{
			ID:                  "22016000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>taxi",
		}
	case "22017000":

		return TransactionMap{
			ID:                  "22017000",
			PhysicalLocation:    false,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>tolls and fees",
		}
	case "22018000":

		return TransactionMap{
			ID:                  "22018000",
			PhysicalLocation:    true,
			TransactionType:     XactionCharge,
			Description:         "travel",
			DetailedDescription: "travel=>transportation centers",
		}
	}
	return TransactionMap{
		ID:                  transaction.CategoryID,
		PhysicalLocation:    false,
		TransactionType:     XactionCharge,
		Description:         "unknown",
		DetailedDescription: "unknown",
	}
}
