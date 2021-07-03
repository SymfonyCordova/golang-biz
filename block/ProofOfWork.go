package block

import (
	"bytes"
	"crypto/sha256"
	"log"
	"math/big"
)

//定义一个工作量证明的结构ProofOfWork
type ProofOfWork struct {
	block *Block //block
	target *big.Int//目标值 big.Int 一个非常大的数,它有很多丰富的方法:比较,赋值方法
}

//提供创建POW的函数
func NewProofOfWork(block *Block)*ProofOfWork{
	pow := ProofOfWork{
		block: block,
	}

	//我们指定的难度值
	targetStr := "0000f00000000000000000000000000000000000000000000000000000000000"
	//引入的辅助变量,目的是将上面的难度值转换成big.int
	tmpInt := big.Int{}
	//将难度值赋值给big.int,指定16进制的格式
	tmpInt.SetString(targetStr, 16)

	pow.target = &tmpInt
	return &pow
}

//提供不断计算的hash函数
func (pow *ProofOfWork)Run()([]byte, uint64){
	var nonce uint64
	var hash [32]byte
	block := pow.block

	log.Println("开始挖矿")
	for  {
		//拼接数据(区块的数据,还有不断变化的随机数)
		//挖矿的本质是对区块头进行哈希
			//而根梅克尔根是从区块体中里面一个个交易二叉树进行哈希得到的
			//这样确保区块体的数据能够影响到区块头的哈希
		tmp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nonce),
			//只对区块头做哈希,不对区块体哈希 区块体通过MerkelRoot对产生影响
			//block.Data,
		}

		//将二维切片数组连接起来,生成一个一维切片
		blockInfo := bytes.Join(tmp, []byte{})

		//做哈希运算
		hash = sha256.Sum256(blockInfo)
		//与pow中target进行比较
		tmpInt := big.Int{}
		//将我们得到hash	数组转成一个big.int
		tmpInt.SetBytes(hash[:])

		//比较当前的哈希值与目标的哈希值,如果当前的哈希值小于目标的哈希值,就说明找到了,否则继续找
		// -1 x < y
		//  0 x == y
		// +1 x > y
		ret := tmpInt.Cmp(pow.target)
		if ret == -1{
			//找到了,退出返回
			log.Printf("挖矿成功 hash: %x, nonce:%d\n", hash, nonce)
			break
		}else{
			//没找到,继续找,随机数加1
			nonce++
		}
	}

	return hash[:], nonce
}

//提供一个校验函数
func IsValid(){

}

