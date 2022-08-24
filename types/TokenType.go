package types

type TokenType struct {
	Address  string `json:"address" bson:"address"`
	Name     string `json:"name" bson:"name"`
	Decimals uint8  `json:"decimals" bson:"decimals"`
}
