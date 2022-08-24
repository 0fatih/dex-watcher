package listener

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"dex-watcher/contracts/pair"
	"dex-watcher/db"
	"dex-watcher/global"
	"dex-watcher/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
)

// listenRouter subscribes to `pairs` for new transactions
func listenRouter(pairs []common.Address) error {
	color.Green("[~] Starting to listening...")

	query := ethereum.FilterQuery{
		Addresses: pairs,
	}

	logs := make(chan types.Log)
	sub, err := global.Client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return errors.New("subscription failed")
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(pair.PairABI)))
	if err != nil {
		return errors.New("contract abi construction failed")
	}

	for {
		select {
		case err := <-sub.Err():
			color.Red("[!] Error: listening failed", err)
			return errors.New("listening failed")
		case vLog := <-logs:
			if vLog.Topics[0] == common.HexToHash("0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1") {
				out, err := contractAbi.Unpack("Sync", vLog.Data)
				if err != nil {
					color.Red("[!] Error with parsing transaction")
				}

				reserve0 := out[0].(*big.Int)
				reserve1 := out[1].(*big.Int)

				name0, name1, err := db.GetNames(vLog.Address)
				if err != nil {
					utils.ColoredPrint("[!] error with getting names", "yellow")
					msg := fmt.Sprintf("[+] %v reserve0: %v reserve1: %v", vLog.Address.String(), reserve0.String(), reserve1.String())
					utils.ColoredPrint(msg, "green")
				} else {
					msg := fmt.Sprintf("[+] %v/%v reserve0: %v reserve1: %v", name0, name1, reserve0.String(), reserve1.String())
					utils.ColoredPrint(msg, "green")
				}
			}
		}
	}
}

func StartListening() {
	pairList := db.GetAllPairAddresses()

	listenRouter(pairList)
}
