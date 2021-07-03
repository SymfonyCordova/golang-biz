package block

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

//引入区块链
type BlockChain struct {
	//使用数据库代替数组
	db *bolt.DB
	tail []byte //存储最后一个区块的哈希
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

//定义一个区块链
func NewBlockChain(address string)*BlockChain{
	//最后一个区块的哈希, 从数据库中读出来
	var lastHash []byte

	//1 打开数据库
	db,err := bolt.Open(blockChainDb,0660,nil)
	if err != nil {
		log.Panic(fmt.Sprintf("打开数据库失败:%s", err))
	}

	//tx是事务
	err = db.Update(func(tx *bolt.Tx) error {
		//找到抽屉bucket(如果没有就创建)
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil{
			//没有抽屉,我们需要创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(fmt.Sprintf("创建bucket(blockBucket)失败:%s", err))
			}

			//创建一个创世块,并做为第一个区块链添加到区块链中
			genesisBlock := GenesisBlock(address)

			//写数据
			err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(fmt.Printf("创建一个创世块失败:%s", err))
			}

			err = bucket.Put([]byte("lastHashKey"), genesisBlock.Hash)
			if err != nil {
				log.Panic(fmt.Printf("创建一个创世块失败:%s", err))
			}

			lastHash = genesisBlock.Hash
		}else{
			lastHash = bucket.Get([]byte("lastHashKey"))
		}

		return nil
	})

	return &BlockChain{db, lastHash}
}

//创始块
func GenesisBlock(address string)*Block{
	coinbase := NewCoinbaseTx(address, "内蒙古大佬矿工")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//添加区块
func (bc *BlockChain)AddBlock(txs []*Transaction){
	//如何获取前区块的哈希值
	//创建新的区块 添加到区块链db中
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket 不应该为空 请检查！")
		}

		//创建新的区块
		block := NewBlock(txs, bc.tail)

		//添加到区块链db中
		err := bucket.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panic(fmt.Printf("创建一个区块失败:%s", err))
		}

		err = bucket.Put([]byte("lastHashKey"), block.Hash)
		if err != nil {
			log.Panic(fmt.Printf("创建一个区块失败:%s", err))
		}

		//更新一下内存中的区块链,指的是把最后的小尾巴tail更新一下
		bc.tail = block.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

//找到指定地址的所有的utxo 就是未消费的支出
func (bc *BlockChain)FindUTXOs(pubKeyHash []byte)[]TXOutput{
	var utxos []TXOutput

	txs := bc.FindUTXOTransactions(pubKeyHash)
	for _, tx:=range txs{
		for _, output := range tx.TXOutputs{
			if bytes.Equal( pubKeyHash, output.PubKeyHash ){
				utxos = append(utxos, output)
			}
		}
	}

	return utxos
}

func (bc *BlockChain) PrintChain() {
	blockHeight := 0
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		err := b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("lastHashKey")){
				return nil
			}

			block := UnSerialize(v)

			fmt.Printf("=============== 区块高度: %d ====================\n", blockHeight)
			blockHeight++
			fmt.Printf("版本号:%d\n", block.Version)
			fmt.Printf("前区块哈希值: %x\n", block.PrevHash)
			fmt.Printf("梅克儿根: %x\n", block.MerkelRoot)
			fmt.Printf("时间戳: %s\n", time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05"))
			fmt.Printf("难度值: %d\n", block.Difficulty)
			fmt.Printf("随机数: %d\n", block.Nonce)

			fmt.Printf("当前区块哈希值 %x\n", block.Hash)
			fmt.Printf("区块数据 :%s\n", block.Transactions[0].TXInputs[0].PubKey)
			return nil
		})


		if err != nil {
			return err
		}

		return nil
	})
	
	if err != nil {
		return
	}
}

//根据需要找到合理的utxos
func (bc *BlockChain) FindNeedUTXOs(senderPubKeyHash []byte, amount float64) (map[string][]uint64, float64) {
	//找到的合理的utxos集合
	utxos :=  make(map[string][]uint64)
	var calc float64//找到的utxos里面包含钱的总数

	txs := bc.FindUTXOTransactions(senderPubKeyHash)
	for _, tx:=range txs{
		for index, output := range tx.TXOutputs{
			//这个output和我们目标的地址相同,满足条件,加到返回utxo数组中
			if bytes.Equal(senderPubKeyHash, output.PubKeyHash){
				if calc < amount {
					//把utxo加进来
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(index))
					//统计一下当前utxo的总量
					calc += output.TransferAmount
					//比较一下是否满足转帐需求
					//a.满足的话直接返回utxos calc
					//b.不满足继续统计
					if calc >= amount {
						return utxos, calc
					}
				}
			}

		}
	}

	return utxos, calc
}

func (bc *BlockChain)FindUTXOTransactions(senderPubKeyHash []byte) []*Transaction {
	var txs []*Transaction //存储所有包含utxo交易集合
	//map[交易id]
	spentOutputs := make(map[string][]int64)
	//遍历区块
	//遍历交易
	//遍历output 找到和自己相关的utxo(在添加output之前检查一下是否已经消耗过)
	//遍历input 找到自己花费过的utxo的集合(把自己消耗过的标识出来)

	//遍历区块
	it := bc.NewIterator()
	for  {
		block := it.Next()

		//遍历交易
		for _, tx := range block.Transactions {

		OUTPUT:
			//遍历output  找到和自己相关的utxo(在添加output之前检查一下是否已经消耗过)
			for index, output := range tx.TXOutputs {
				//在这里做一个过滤，将所有消耗过的output和当前的所即将添加output对比一下
				//如果相同,则跳过,否则添加
				//map[2222] = []int64{0}
				//map[3333] = []int64{0, 1}
				if spentOutputs[string(tx.TXID)] != nil{
					for _, j := range spentOutputs[string(tx.TXID)]{
						if int64(index) == j{
							//当前准备添加output已经消耗过了,不要再添加了
							continue OUTPUT
						}
					}
				}

				//这个output和我们目标的地址相同,满足条件,加到返回utxo数组中
				//if output.PubKeyHash == address{
				if bytes.Equal(senderPubKeyHash, output.PubKeyHash){
					//utxo = append(utxo, output)
					//查找所有包含我的outx的交易集合
					txs = append(txs, tx)
				}
			}

			//如果当前交易是挖矿交易的话,那么不做遍历,直接跳过
			if !tx.IsCoinbase() {

				//遍历input 找到自己花费过的utxo的集合(把自己消耗过的标识出来)
				for _, input := range tx.TXInputs {
					//判断一下当前这个input和目标李四是否一致 如果相同说明这个是李四消耗过的output 就加进来map
					pubKeyHash := HashPubKey(input.PubKey)
					if bytes.Equal(pubKeyHash, senderPubKeyHash) {
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			}
		}
		if len(block.PrevHash) == 0{
			break
		}
	}

	return txs
}