// Package global holds some global variables
package global

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client *ethclient.Client
var DBConnection *mongo.Client

var FactoryCollection *mongo.Collection
var PairCollection *mongo.Collection
var TokenCollection *mongo.Collection

func InitGlobalVariables(client *ethclient.Client, dbConnection *mongo.Client) {
	Client = client
	DBConnection = dbConnection

	FactoryCollection = DBConnection.Database("dex-watcher").Collection("factories")
	PairCollection = DBConnection.Database("dex-watcher").Collection("pairs")
	TokenCollection = DBConnection.Database("dex-watcher").Collection("tokens")
}
