package proofofwork

import (
	"math"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/tclchiam/oxidize-go/blockchain/engine/mining"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/identity"
)

const (
	maxNonce = math.MaxUint64
)

var (
	defaultWorkerCount = uint(runtime.NumCPU())
)

type miner struct {
	workerCount uint
	beneficiary *identity.Address
}

func NewMiner(workerCount uint, beneficiary *identity.Address) mining.Miner {
	return &miner{workerCount: workerCount, beneficiary: beneficiary}
}

func NewDefaultMiner(beneficiary *identity.Address) mining.Miner {
	return NewMiner(defaultWorkerCount, beneficiary)
}

func (miner *miner) MineBlock(parent *entity.BlockHeader, transactions entity.Transactions) *entity.Block {
	reward := entity.NewRewardTx(miner.beneficiary)

	return miner.mineBlock(parent, transactions.Add(reward), uint64(time.Now().Unix()))
}

func (miner *miner) mineBlock(parent *entity.BlockHeader, transactions entity.Transactions, now uint64) *entity.Block {
	transactionsHash := mining.CalculateTransactionsHash(transactions)

	input := &mining.BlockHashingInput{
		Index:            parent.Index + 1,
		PreviousHash:     parent.Hash,
		Timestamp:        now,
		TransactionsHash: transactionsHash,
		Difficulty:       parent.Difficulty,
	}

	solutions := make(chan *BlockSolution)
	nonces := make(chan uint64, miner.workerCount)
	defer close(nonces)

	for workerNum := uint(0); workerNum < miner.workerCount; workerNum++ {
		go worker(input, nonces, solutions)
	}

	for nonce := uint64(0); nonce < maxNonce; nonce++ {
		select {
		case solution := <-solutions:
			header := entity.NewBlockHeader(
				input.Index,
				input.PreviousHash,
				input.TransactionsHash,
				input.Timestamp,
				solution.Nonce,
				solution.Hash,
				input.Difficulty,
			)
			return entity.NewBlock(header, transactions)
		default:
			nonces <- nonce
		}
	}

	log.Panic(MaxNonceOverflowError)
	return nil
}

func worker(work *mining.BlockHashingInput, nonces <-chan uint64, solutions chan<- *BlockSolution) {
	for nonce := range nonces {
		hash := mining.CalculateBlockHash(work, nonce)

		if mining.HasDifficulty(hash, work.Difficulty) {
			solutions <- &BlockSolution{Nonce: nonce, Hash: hash}
		}
	}
}
