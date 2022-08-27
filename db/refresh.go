package db

import (
	"context"
	PairContract "dex-watcher/contracts/pair"
	"dex-watcher/globals"
	dexTypes "dex-watcher/types"
	"dex-watcher/utils"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
)

// RefreshDB fetchs latest prices and writes to database.
func RefreshDB() error {
	// Find all pairs
	cur, err := globals.PairCollection.Find(context.Background(), bson.D{})
	if err != nil {
		utils.ColoredPrint("[!] Something went wrong wile finding pairs!", utils.PrintColors.RED)
		return errors.New("something went wrong wile finding pairs")
	}
	defer cur.Close(context.Background())

	// Iterate over pairs
	for cur.Next(context.Background()) {
		var pair dexTypes.PairType

		err := cur.Decode(&pair)
		if err != nil {
			utils.ColoredPrint("[!] Failed to decode pair "+pair.Address, utils.PrintColors.RED)
			continue
		}

		pair.Reserve0, pair.Reserve1, err = getReserves(common.HexToAddress(pair.Address))
		if err != nil {
			utils.ColoredPrint("[!] Failed while fetching reserves for: "+pair.Address, utils.PrintColors.RED)
			continue
		}

		_, err = globals.PairCollection.UpdateOne(context.Background(), bson.M{"address": pair.Address}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "reserve0", Value: pair.Reserve0}}},
			{Key: "$set", Value: bson.D{{Key: "reserve1", Value: pair.Reserve1}}},
		})
		if err != nil {
			utils.ColoredPrint("[!] Failed to add pair "+pair.Address+" to database!"+" -> "+err.Error(), utils.PrintColors.RED)
			return errors.New("failed to add pair")
		}

		utils.ColoredPrint("[+] "+pair.Address+" updated.", utils.PrintColors.GREEN)
	}

	return nil
}

func getReserves(address common.Address) (string, string, error) {
	pairInstance, err := PairContract.NewPair(address, globals.Client)
	if err != nil {
		utils.ColoredPrint("[!] Pair instance construction failed: "+address.String()+" -> "+err.Error(), utils.PrintColors.RED)
		return "", "", errors.New("pair instance construction failed")
	}

	reserves, err := pairInstance.GetReserves(nil)
	if err != nil {
		if err != nil {
			utils.ColoredPrint("[!] Failed while getting reserves for: "+address.String()+" -> "+err.Error(), utils.PrintColors.RED)
			return "", "", errors.New("failed while getting reserves")
		}
	}

	return reserves.Reserve0.String(), reserves.Reserve1.String(), nil
}
