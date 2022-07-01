package data

type Account struct {
	AccountId int      `json:"account_id" bson:"account_id"`
	Limit     int      `json:"limit" bson:"limit"`
	Products  []string `json:"products "bson:"products"`
}
