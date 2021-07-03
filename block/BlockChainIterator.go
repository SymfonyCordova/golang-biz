package block

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	db *bolt.DB
	//游标，用于不断索引
	currentHashPointer []byte
}

func (bc *BlockChain)NewIterator() *BlockChainIterator{
	return &BlockChainIterator {
		bc.db,
		//最初指向区块链的最后一个区块,随着Next的调用,不断变化
		bc.tail,
	}
}

//迭代器属于区块链的
//Next属于迭代器的
//1.返回当前的区块
//2.指针前移
func (it *BlockChainIterator)Next() *Block {
	var block Block
	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("迭代器遍历bucket不应该为空,请检查！")
		}

		blockTmp := bucket.Get(it.currentHashPointer)
		//解码
		block = UnSerialize(blockTmp)
		//游标哈希向前移动
		it.currentHashPointer = block.PrevHash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &block
}