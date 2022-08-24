// Package db is for writing/reading data from the MondoDB
package db

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	FactoryContract "dex-watcher/contracts/factory"
	PairContract "dex-watcher/contracts/pair"
	TokenContract "dex-watcher/contracts/token"
	"dex-watcher/global"
	dexTypes "dex-watcher/types"
	"dex-watcher/utils"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var wantedPairs int64

// InitializeDB creates collections for factories, pairs and tokens.
func InitializeDB(factoryAddresses []common.Address, _wantedPairs int64, wgInitialize *sync.WaitGroup) {
	wantedPairs = _wantedPairs

	fmt.Println("factories:", factoryAddresses)
	wgFactory := new(sync.WaitGroup)
	for _, address := range factoryAddresses {
		wgFactory.Add(1)
		go factoryHandler(address, wgFactory)
	}
	wgFactory.Wait()
	utils.ColoredPrint("[+] Initializing done!", "green")
	wgInitialize.Done()
}

// factoryHandler adds pairs to database in a factory.
func factoryHandler(address common.Address, wgFactory *sync.WaitGroup) {
	utils.ColoredPrint("[~] Handling factory "+address.String(), "yellow")

	instance, err := FactoryContract.NewFactory(address, global.Client)
	if err != nil {
		utils.ColoredPrint("[!] Factory instance construction failed: "+address.String()+" -> "+err.Error(), "red")
		wgFactory.Done()
		return
	}

	// add factory to database
	_, err = global.FactoryCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "address", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		utils.ColoredPrint("[!] Failed to creating index -> "+err.Error(), "red")
		wgFactory.Done()
		return
	}

	f := dexTypes.FactoryType{Address: address.String()}
	_, err = global.FactoryCollection.InsertOne(context.Background(), f)
	if err != nil {
		if !isDup(err) {
			utils.ColoredPrint("[!] Failed to add factory "+address.String()+" to database! -> "+err.Error(), "red")
			wgFactory.Done()
			return
		}
	}

	pairsLength := big.NewInt(wantedPairs)
	if wantedPairs == 0 {
		pairsLength, err = instance.AllPairsLength(nil)
		if err != nil {
			utils.ColoredPrint("[!] Failed to fetch pairs length: "+address.String()+" -> "+err.Error(), "red")
			wgFactory.Done()
			return
		}
	}

	wgPair := new(sync.WaitGroup)
	var i int64
	for ; i < pairsLength.Int64(); i++ {
		wgPair.Add(1)
		go pairHandler(instance, i, wgPair)
	}
	wgPair.Wait()
	wgFactory.Done()
}

func pairHandler(factoryInstance *FactoryContract.Factory, id int64, wgPair *sync.WaitGroup) {
	pairAddress, err := factoryInstance.AllPairs(nil, big.NewInt(id))
	if err != nil {
		utils.ColoredPrint("[!] Failed to fetch pair address: "+strconv.Itoa(int(id))+" -> "+err.Error(), "red")
		wgPair.Done()
		return
	}

	utils.ColoredPrint("[~] Handling pair "+pairAddress.String(), "yellow")

	pairInstance, err := PairContract.NewPair(pairAddress, global.Client)
	if err != nil {
		utils.ColoredPrint("[!] Pair instance construction failed: "+pairAddress.String()+" -> "+err.Error(), "red")
		wgPair.Done()
		return
	}

	token0Address, err := pairInstance.Token0(nil)
	if err != nil {
		utils.ColoredPrint("[!] Failed to fetch token0 for:"+pairAddress.String()+" -> "+err.Error(), "red")
		wgPair.Done()
		return
	}

	token1Address, err := pairInstance.Token1(nil)
	if err != nil {
		utils.ColoredPrint("[!] Failed to fetch token1 for:"+pairAddress.String()+" -> "+err.Error(), "red")
		wgPair.Done()
		return
	}

	wgToken := new(sync.WaitGroup)
	wgToken.Add(1)
	go handleToken(token0Address, wgToken)
	wgToken.Add(1)
	go handleToken(token1Address, wgToken)
	wgToken.Wait()

	_, err = global.PairCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "address", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		utils.ColoredPrint("[!] Failed to creating index -> "+err.Error(), "red")
		wgPair.Done()
		return
	}

	p := dexTypes.PairType{Address: pairAddress.String(), Token0Address: token0Address.String(), Token1Address: token1Address.String()}
	_, err = global.PairCollection.InsertOne(context.Background(), p)
	if err != nil {
		if !isDup(err) {
			utils.ColoredPrint("[!] Failed to add pair"+pairAddress.String()+"to database!"+" -> "+err.Error(), "red")
			wgPair.Done()
			return
		}
	}
	wgPair.Done()
}

func handleToken(address common.Address, wgToken *sync.WaitGroup) {
	utils.ColoredPrint("[~] Handling token "+address.String(), "yellow")

	tokenInstance, err := TokenContract.NewToken(address, global.Client)
	if err != nil {
		utils.ColoredPrint("[!] Token0 instance construction failed: "+address.String()+" -> "+err.Error(), "red")
		wgToken.Done()
		return
	}

	name, err := tokenInstance.Name(nil)
	if err != nil {
		utils.ColoredPrint("[!] Token0 name fetching failed: "+address.String()+" -> "+err.Error(), "red")
		wgToken.Done()
		return
	}

	decimals, err := tokenInstance.Decimals(nil)
	if err != nil {
		utils.ColoredPrint("[!] Token0 decimals fetching failed: "+address.String()+" -> "+err.Error(), "red")
		wgToken.Done()
		return
	}

	_, err = global.TokenCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "address", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		utils.ColoredPrint("[!] Failed to creating index -> "+err.Error(), "red")
		wgToken.Done()
		return
	}

	t := dexTypes.TokenType{Address: address.String(), Name: name, Decimals: decimals}
	_, err = global.TokenCollection.InsertOne(context.Background(), t)
	if err != nil {
		if !isDup(err) {
			utils.ColoredPrint("[!] Failed to add pair"+address.String()+"to database!"+" -> "+err.Error(), "red")
			wgToken.Done()
			return
		}
	}

	wgToken.Done()
}

func isDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}
