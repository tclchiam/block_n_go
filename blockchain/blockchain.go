package blockchain

import (
	"github.com/tclchiam/block_n_go/blockchain/engine"
	"github.com/tclchiam/block_n_go/blockchain/engine/iter"
	"github.com/tclchiam/block_n_go/blockchain/entity"
	"github.com/tclchiam/block_n_go/blockchain/entity/encoding"
	"github.com/tclchiam/block_n_go/storage"
	"github.com/tclchiam/block_n_go/mining"
	"github.com/tclchiam/block_n_go/wallet"
	"github.com/tclchiam/block_n_go/blockchain/engine/utxo"
)

type Blockchain struct {
	blockRepository storage.BlockRepository
	miner           mining.Miner
	utxoEngine      utxo.Engine
}

func Open(repository storage.BlockRepository, miner mining.Miner, ownerAddress string) (*Blockchain, error) {
	exists, err := genesisBlockExists(repository)
	if err != nil {
		return nil, err
	}

	if !exists {
		transactionEncoder := encoding.TransactionProtoEncoder()

		transactions := entity.Transactions{entity.NewCoinbaseTx(ownerAddress, transactionEncoder)}
		header := entity.NewGenesisBlockHeader(transactions)
		solution := miner.MineBlock(header)

		err := repository.SaveBlock(entity.NewBlock(header, solution, transactions))
		if err != nil {
			return nil, err
		}
	}

	return &Blockchain{
		blockRepository: repository,
		miner:           miner,
		utxoEngine:      utxo.NewCrawlerEngine(repository),
	}, nil
}

func genesisBlockExists(repository storage.BlockRepository) (bool, error) {
	head, err := repository.Head()
	if err != nil {
		return false, err
	}
	if head == nil {
		return false, nil
	}
	return true, nil
}

func (bc *Blockchain) ForEachBlock(consume func(*entity.Block)) error {
	return iter.ForEachBlock(bc.blockRepository, consume)
}

func (bc *Blockchain) ReadBalance(address string) (uint32, error) {
	return engine.ReadBalance(address, bc.utxoEngine)
}

func (bc *Blockchain) Send(sender, receiver, coinbase *wallet.Wallet, expense uint32) error {
	expenseTransaction, err := engine.BuildExpenseTransaction(sender, receiver, expense, bc.utxoEngine)
	if err != nil {
		return err
	}
	rewardTransaction := entity.NewCoinbaseTx(coinbase.GetAddress(), encoding.TransactionProtoEncoder())

	newBlock, err := engine.MineBlock([]*entity.Transaction{expenseTransaction, rewardTransaction}, bc.miner, bc.blockRepository)
	if err != nil {
		return err
	}
	return bc.blockRepository.SaveBlock(newBlock)
}
