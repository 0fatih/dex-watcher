package db

import (
	"context"
	"dex-watcher/global"
	dexTypes "dex-watcher/types"
	"dex-watcher/utils"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
)

func GetNames(pairAddress common.Address) (string, string, error) {
	var res dexTypes.PairType
	err := global.PairCollection.FindOne(context.Background(), bson.D{{Key: "address", Value: pairAddress.String()}}).Decode(&res)
	if err != nil {
		return "", "", errors.New("something went wrong while finding pair")
	}

	var token0 dexTypes.TokenType
	err = global.TokenCollection.FindOne(context.Background(), bson.D{{Key: "address", Value: res.Token0Address}}).Decode(&token0)
	if err != nil {
		return "", "", errors.New("something went wrong while finding for token0 name")
	}

	var token1 dexTypes.TokenType
	err = global.TokenCollection.FindOne(context.Background(), bson.D{{Key: "address", Value: res.Token1Address}}).Decode(&token1)
	if err != nil {
		return "", "", errors.New("something went wrong while finding for token1 name")
	}

	return token0.Name, token1.Name, nil
}

func GetAllPairAddresses() []common.Address {
	var results []common.Address

	cur, err := global.PairCollection.Find(context.Background(), bson.D{})
	if err != nil {
		utils.ColoredPrint("[!] Getting all pairs failed! -> "+err.Error(), "red")
		panic("can not get pairs")
	}

	for cur.Next(context.Background()) {
		var res dexTypes.PairType

		err := cur.Decode(&res)
		if err != nil {
			utils.ColoredPrint("[!] Failed to decode pair", "red")
			continue
		}

		results = append(results, common.HexToAddress(res.Address))
	}

	return results
}
