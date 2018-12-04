package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tomochain/backend-matching-engine/errors"
	demo "github.com/tomochain/orderbook/common"
)

const (
	// time out 15 seconds
	timeout = 15
)

func (l *Listener) Start(rpcServer string) error {

	demo.LogInfo("EthereumListener starting")

	blockNumber, err := l.Storage.GetEthereumBlockToProcess()
	if err != nil {
		err = errors.Wrap(err, "Error getting ethereum block to process from DB")
		demo.LogError(err.Error())
		return err
	}

	// Check if connected to correct network
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout*time.Second))
	defer cancel()
	id, err := l.Client.NetworkID(ctx)
	if err != nil {
		err = errors.Wrap(err, "Error getting ethereum network ID")
		demo.LogError(err.Error())
		return err
	}

	if id.String() != l.NetworkID {
		return errors.Errorf("Invalid network ID (have=%s, want=%s)", id.String(), l.NetworkID)
	}

	go l.processBlocks(blockNumber)
	return nil
}

func (l *Listener) Stop() error {
	ethClient := l.Client.(*ethclient.Client)
	ethClient.Close()
	l.Enabled = false
	return nil
}

func (l *Listener) processBlocks(blockNumber uint64) {
	if blockNumber == 0 {
		logger.Info("Starting from the latest block")
	} else {
		logger.Infof("Starting from block %d", blockNumber)
	}

	// Time when last new block has been seen
	lastBlockSeen := time.Now()
	noBlockWarningLogged := false

	for {
		if l.Enabled == false {
			// stop listener
			break
		}
		block, err := l.getBlock(blockNumber)
		if err != nil {
			logger.Errorf("Error getting block, blockNumber: %d", blockNumber)
			time.Sleep(1 * time.Second)
			continue
		}

		// Block doesn't exist yet
		if block == nil {
			if time.Since(lastBlockSeen) > 3*time.Minute && !noBlockWarningLogged {
				logger.Warningf("No new block in more than 3 minutes")
				noBlockWarningLogged = true
			}

			time.Sleep(1 * time.Second)
			continue
		}

		// Reset counter when new block appears
		lastBlockSeen = time.Now()
		noBlockWarningLogged = false

		if block.NumberU64() == 0 {
			logger.Error("Ethereum node is not synced yet. Unable to process blocks")
			time.Sleep(30 * time.Second)
			continue
		}

		if l.TransactionHandler == nil {
			// waiting for handler
			time.Sleep(1 * time.Second)
			continue
		}

		err = l.processBlock(block)
		if err != nil {
			logger.Errorf("Error processing block, blockNumber: %d, err: %v", block.NumberU64(), err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Persist block number
		err = l.Storage.SaveLastProcessedEthereumBlock(blockNumber)
		if err != nil {
			logger.Errorf("Error saving last processed block: %s", err)
			time.Sleep(1 * time.Second)
			// We continue to the next block
		}

		blockNumber = block.NumberU64() + 1
	}
}

// getBlock returns (nil, nil) if block has not been found (not exists yet)
func (l *Listener) getBlock(blockNumber uint64) (*types.Block, error) {
	var blockNumberInt *big.Int
	if blockNumber > 0 {
		blockNumberInt = big.NewInt(int64(blockNumber))
	}

	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	block, err := l.Client.BlockByNumber(ctx, blockNumberInt)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		err = errors.Wrap(err, "Error getting block from geth")
		logger.Errorf("Got err: %s, block: %d", err.Error(), blockNumberInt.String())
		return nil, err
	}

	return block, nil
}

func (l *Listener) processBlock(block *types.Block) error {
	transactions := block.Transactions()

	blockTime := time.Unix(block.Time().Int64(), 0)
	logger.Infof("Processing block: blockNumber:%d, blockTime:%v, transactions:%d",
		block.NumberU64(),
		blockTime,
		len(transactions),
	)

	for _, transaction := range transactions {
		to := transaction.To()
		if to == nil {
			// Contract creation
			continue
		}

		// this is the address that we need to check in address association
		// server will store associate like ethereumAddress => userAddress
		// user will store in feed with topic ethereum like {ethereumAddress}
		// if server has problem, client can use this feed to refer to ethereumAddress
		// and check for transaction, then check against current blockchain
		// for server, it is just an indexer to help process faster
		tx := Transaction{
			Hash:     transaction.Hash().Hex(),
			ValueWei: transaction.Value(),
			To:       to.Hex(),
		}
		err := l.TransactionHandler(tx)
		if err != nil {
			logger.Errorf("Error processing transaction: %s", err.Error())
			return errors.Wrap(err, "Error processing transaction")
		}
	}

	// logger.Infof("Processed block")

	return nil
}
