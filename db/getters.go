package db

import (
	"context"
	"dex-watcher/globals"
	dexTypes "dex-watcher/types"
	"dex-watcher/utils"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
)

func GetNamesForPairFromDB(pairAddress common.Address) (string, string, error) {
	var pair dexTypes.PairType
	err := globals.PairCollection.FindOne(context.Background(), bson.D{{Key: "address", Value: pairAddress.String()}}).Decode(&pair)
	if err != nil {
		return "", "", errors.New("something went wrong while finding pair")
	}

	token0, err := GetTokenName(common.HexToAddress(pair.Token0Address))
	if err != nil {
		return "", "", err
	}

	token1, err := GetTokenName(common.HexToAddress(pair.Token1Address))
	if err != nil {
		return "", "", err
	}

	return token0.Name, token1.Name, nil
}

func GetTokenName(address common.Address) (dexTypes.TokenType, error) {
	var token dexTypes.TokenType
	err := globals.TokenCollection.FindOne(context.Background(), bson.D{{Key: "address", Value: address}}).Decode(&token)
	if err != nil {
		return dexTypes.TokenType{}, errors.New("failed to finding name for token")
	}

	return token, nil
}

func GetAllPairAddressesFromDB() []common.Address {
	var results []common.Address

	cur, err := globals.PairCollection.Find(context.Background(), bson.D{})
	if err != nil {
		utils.ColoredPrint("[!] Getting all pairs failed! -> "+err.Error(), utils.PrintColors.RED)
		panic("can not get pairs")
	}

	for cur.Next(context.Background()) {
		var res dexTypes.PairType

		err := cur.Decode(&res)
		if err != nil {
			utils.ColoredPrint("[!] Failed to decode pair", utils.PrintColors.RED)
			continue
		}

		results = append(results, common.HexToAddress(res.Address))
	}

	return results
}
