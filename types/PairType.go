package types

type PairType struct {
	Address       string `json:"address" bson:"address"`
	Token0Address string `json:"token0Address" bson:"token0Address"`
	Token1Address string `json:"token1Address" bson:"token1Address"`
	Reserve0      string `json:"reserve0" bson:"reserve0"`
	Reserve1      string `json:"reserve1" bson:"reserve1"`
}
