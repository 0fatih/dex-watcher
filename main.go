// Package dexlogger is a tool for watching transaction for the dexes who
// implements Uniswap V2
package main

import (
	"context"
	"dex-watcher/db"
	"dex-watcher/global"
	"dex-watcher/listener"
	"dex-watcher/utils"
	"flag"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var factoryList = []common.Address{
	common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"), // Uniswap V2
}

func main() {
	isListen := flag.Bool("listen", false, "Start listening to the blockchain.")
	initDB := flag.Bool("initialize", false, "Creates documents and writes initial data to database.")
	refreshDB := flag.Bool("refresh", false, "Fetchs latest prices and writes to database for every pair.")
	wantedPairs := flag.Int64("pairs", 5, "Pair amount per factory that you want to subscribe. If you want to subscribe to all pairs in a factory, enter 0.")
	flag.Parse()

	if *initDB && *refreshDB {
		utils.ColoredPrint("[!] You cant initialize and refresh database at the same time!", "yellow")
		return
	}

	client, err := ethclient.Dial("YOUR_PROVIDER_HERE")
	if err != nil {
		color.Red("[!] Error: client construction failed -> ", err)
		return
	}
	color.Green("[~] Provider connection successfull...")
	defer client.Close()

	dbConnection, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		color.Red("[!] Error: database connection failed -> ", err)
		return
	}
	color.Green("[~] Database connection successfull...")
	defer dbConnection.Disconnect(context.Background())

	global.InitGlobalVariables(client, dbConnection)

	if *initDB {
		wgInitialize := new(sync.WaitGroup)
		wgInitialize.Add(1)
		db.InitializeDB(factoryList, *wantedPairs, wgInitialize)
		wgInitialize.Wait()
	}

	if *refreshDB {
		db.RefreshDB()
	}

	if *isListen {
		listener.StartListening()
	}

	utils.ColoredPrint("\n\n[+] Done!", "green")
}
