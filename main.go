// Package dexlogger is a tool for watching transaction for the dexes who
// implements Uniswap V2
package main

import (
	"context"
	"dex-watcher/db"
	"dex-watcher/globals"
	"dex-watcher/listener"
	"dex-watcher/utils"
	"flag"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var factoryList = []common.Address{}

func init() {
	// Load environment variables from .env
	err := godotenv.Load(".env")
	if err != nil {
		utils.ColoredPrint("[!] Failed to load environment variables from the '.env' file.", utils.PrintColors.RED)
	}

	client, err := ethclient.Dial(os.Getenv("PROVIDER_URI"))
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

	globals.InitGlobalVariables(client, dbConnection)

	for _, factoryAddress := range strings.Split(os.Getenv("FACTORY_ADDRESSES"), ",") {
		factoryList = append(factoryList, common.HexToAddress(factoryAddress))
	}
}

func main() {
	// Parse flags
	shouldListen := flag.Bool("listen", false, "Start listening to the blockchain.")
	initDB := flag.Bool("initialize", false, "Creates documents and writes initial data to database.")
	refreshDB := flag.Bool("refresh", false, "Fetchs latest prices and writes to database for every pair.")
	wantedPairs := flag.Int64("pairs", 5, "Pair amount per factory that you want to subscribe. If you want to subscribe to all pairs in a factory, enter 0.")
	flag.Parse()

	if *initDB {
		db.InitializeDB(factoryList, *wantedPairs)
	}

	if *refreshDB {
		err := db.RefreshDB()
		if err != nil {
			os.Exit(1)
		}
	}

	if *shouldListen {
		err := listener.StartListening()
		if err != nil {
			os.Exit(1)
		}
	}

	utils.ColoredPrint("\n[+] Done!", utils.PrintColors.GREEN)
}
