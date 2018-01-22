package memdb

import (
	"fmt"
	"sync"

	"github.com/tclchiam/block_n_go/storage"
	"github.com/tclchiam/block_n_go/blockchain/entity"
)

type blockMemoryRepository struct {
	head *entity.Block
	db   map[*entity.Hash]*entity.Block
	lock sync.RWMutex
}

func NewBlockRepository() storage.BlockRepository {
	return &blockMemoryRepository{db: make(map[*entity.Hash]*entity.Block)}
}

func (r *blockMemoryRepository) Head() (head *entity.Block, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if r.head != nil {
		return r.head, nil
	}
	return nil, nil
}

func (r *blockMemoryRepository) Block(hash *entity.Hash) (*entity.Block, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if block, ok := r.db[hash]; ok {
		return block, nil
	}
	return nil, fmt.Errorf("no block found with hash: '%s'", hash)
}

func (r *blockMemoryRepository) SaveBlock(block *entity.Block) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.db[block.Hash()] = block
	r.head = block

	return nil
}
