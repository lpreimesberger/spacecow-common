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
)

//XactionType is the type of xaction in the db
type XactionType int

const (
	XactionPayment = iota
	XactionCharge
)

//ShipTypes is the ship level in the game
type ShipTypes int

// classes on richness levels
const (
	ScoutShip    = iota // less than 5k
	AberdeenShip        // 5k-50k
	HerefordShip        // 50-100k
	LonghornShip        // 100k-1m
	WagyuShip           // > 1 million
)

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
