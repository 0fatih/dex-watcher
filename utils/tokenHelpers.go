package utils

import (
	"context"
	"dex-watcher/global"
	dexTypes "dex-watcher/types"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
)

// GetPrices returns price0 and price1 for token0 and token1
func GetPrices(r0, r1 *big.Int, pairAddress common.Address) (float64, float64, error) {
	var pair dexTypes.PairType

	err := global.PairCollection.FindOne(context.Background(), bson.D{{Key: "address", Value: pairAddress.String()}}).Decode(&pair)
	if err != nil {
		return 0, 0, errors.New("something went wrong while finding pairs")
	}

	decimal0, err := GetDecimals(pair.Token0Address)
	if err != nil {
		return 0, 0, err
	}

	decimal1, err := GetDecimals(pair.Token1Address)
	if err != nil {
		return 0, 0, err
	}

	amount0In := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(decimal0), nil)) // 1 token
	amount1In := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(decimal1), nil)) // 1 token

	fmt.Println(amount0In.String(), amount1In.String())

	amount0Out := getAmountOut(amount1In, r1, r0)
	amount1Out := getAmountOut(amount0In, r0, r1)

	price0 := amount0Out / float64(amount0In.Uint64())
	price1 := amount1Out / float64(amount1In.Uint64())

	price2 := amount0Out / float64(amount0In.Uint64())
	price3 := amount1Out / float64(amount1In.Uint64())

	fmt.Println("alternative: price0: ", price2, "price1: ", price3)

	return price0, price1, nil
}

func GetDecimals(address string) (int64, error) {
	var token dexTypes.TokenType

	err := global.TokenCollection.FindOne(context.Background(), bson.D{{Key: "address", Value: address}}).Decode(&token)
	if err != nil {
		ColoredPrint("[!] Finding decimals for:"+address+" failed!", "red")
		return 0, errors.New("finding decimals failed")
	}

	return int64(token.Decimals), nil
}

func getAmountOut(amountIn, reserveIn, reserveOut *big.Int) float64 {
	// TODO: Add fee
	amountInWithFee := new(big.Int).Mul(amountIn, big.NewInt(9975))
	numerator := new(big.Int).Mul(amountInWithFee, reserveOut)
	denominator := new(big.Int).Add(new(big.Int).Mul(reserveIn, big.NewInt(10000)), amountInWithFee)
	amountOut := new(big.Int).Quo(numerator, denominator)

	return float64(amountOut.Uint64())
}
