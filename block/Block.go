package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Version uint64 //版本号
	PrevHash []byte //前区块哈希
	//Merkel根 (根梅克尔根,这就是一个哈希值,我们先不管)
	MerkelRoot []byte
	TimeStamp uint64//时间戳
	Difficulty uint64 //难度值
	Nonce uint64//随机数，也就是挖矿要找的数据

	Hash []byte //当前区块哈希
	Transactions []*Transaction  //真实的交易数组
}

//创建区块
func NewBlock(txs []*Transaction, prevBlockHash []byte)*Block{
	block := Block{
		Version: 00,
		PrevHash: prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,

		Nonce: 0,
		Hash: []byte{},
		Transactions: txs,
	}

	block.MerkelRoot = block.MakeMerkelRoot()

	pow := NewProofOfWork(&block)
	//查找随机数,不停的进行哈希运算
	hash, nonce := pow.Run()

	//根据挖矿结果对区块链数据进行更新
	block.Hash = hash
	block.Nonce = nonce

	return &block
}

//实现一个辅助函数,功能是将uint64转成[]byte
func Uint64ToByte(num uint64)[]byte{
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	//使用god进行序列化(编码)得到字节流
	//定义一个编码器
	//使用编码器进行编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic("编码出错")
	}

	return buffer.Bytes()
}

func UnSerialize(data []byte) Block {
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var block Block
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("解码出错")
	}

	return block
}

//模拟梅克尔根 只是对交易的数据做简单的拼接,而不做二叉数处理
func (block *Block)MakeMerkelRoot()[]byte{
	var info []byte
	for _,tx := range block.Transactions{
		//将交易的哈希拼接起来,再整体做哈希处理
		info = append(info, tx.TXID...)
	}

	hash := sha256.Sum256(info)
	return hash[:]
}


/*
//生成哈希
func (block *Block)	SetHash(){
	//拼装数据
	var blockInfo []byte

	//blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
	//blockInfo = append(blockInfo, block.PrevHash...)
	//blockInfo = append(blockInfo, block.MerkelRoot...)
	//blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
	//blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
	//blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
	//blockInfo = append(blockInfo, block.Data...)

	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}

	//将二维切片数组连接起来,生成一个一维切片
	blockInfo := bytes.Join(tmp, []byte{})

	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
*/