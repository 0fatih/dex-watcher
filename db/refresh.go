package db

import (
	"context"
	PairContract "dex-watcher/contracts/pair"
	"dex-watcher/global"
	dexTypes "dex-watcher/types"
	"dex-watcher/utils"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
)

// RefreshDB fetchs latest prices and writes to database.
func RefreshDB() {
	cur, err := global.PairCollection.Find(context.Background(), bson.D{})
	if err != nil {
		utils.ColoredPrint("[!] Something went wrong wile finding pairs!", "red")
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var result dexTypes.PairType

		err := cur.Decode(&result)
		if err != nil {
			utils.ColoredPrint("[!] Failed to decode pair", "red")
			continue
		}

		pairInstance, err := PairContract.NewPair(common.HexToAddress(result.Address), global.Client)
		if err != nil {
			utils.ColoredPrint("[!] Pair instance construction failed: "+result.Address+" -> "+err.Error(), "red")
			return
		}

		reserves, err := pairInstance.GetReserves(nil)
		if err != nil {
			if err != nil {
				utils.ColoredPrint("[!] Failed while getting reserves for: "+result.Address+" -> "+err.Error(), "red")
				return
			}
		}

		result.Reserve0 = reserves.Reserve0.String()
		result.Reserve1 = reserves.Reserve1.String()

		_, err = global.PairCollection.UpdateOne(context.Background(), bson.M{"address": result.Address}, bson.D{
			{"$set", bson.D{{"reserve0", reserves.Reserve0.String()}}},
			{"$set", bson.D{{"reserve1", reserves.Reserve1.String()}}},
		})
		if err != nil {
			utils.ColoredPrint("[!] Failed to add pair"+result.Address+"to database!"+" -> "+err.Error(), "red")
			return
		}

		utils.ColoredPrint("[+] "+result.Address+" updated.", "green")
	}
}
