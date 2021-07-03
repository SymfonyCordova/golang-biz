package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	bizEncoding "github.com/SymfonyCordova/golang-biz/encoding"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//这里的钱包是一个结构,每一个钱包保存了公钥,私钥对
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	//约定,这里的PubKey不存储原始的公钥,而是存储X和Y的拼接的字符串,在校验端重新拆分(参考r, s传递)
	PubKey []byte
}

//创建钱包
func NewWallet()*Wallet{
	//创建曲线
	cure := elliptic.P256()

	//生成私钥
	privateKey, err := ecdsa.GenerateKey(cure, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	//生成公钥
	pubKey := privateKey.PublicKey

	//拼接 X，Y
	wrapperPubKey := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)

	return &Wallet{privateKey, wrapperPubKey}
}

//生成地址
func (wallet *Wallet)NewAddress()string{
	pubKey := wallet.PubKey

	rip160Hash := HashPubKey(pubKey)

	//拼接version
	version := byte(00)
	payload := append([]byte{version}, rip160Hash...)

	//checksum
	checkCode := CheckSum(payload)

	//25字节数据
	payload = append(payload, checkCode...)

	//go语言有一个库,叫做btcd,这个是go语言实现的比特币节点源码
	address := bizEncoding.Base58Encode(payload)

	return string(address)
}

func HashPubKey(data []byte)[]byte{
	hash := sha256.Sum256(data)

	//理解成编码器
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}

	//返回rip160哈希结果
	rip160Hash :=rip160hasher.Sum(nil)
	return rip160Hash
}

func CheckSum(data []byte)[]byte{
	//两次sha256
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])

	//前4字节校验码
	checkCode := hash2[:4]
	return checkCode
}

func IsValidAddress(address string) bool{
	//1.解码
	addressByte := bizEncoding.Base58Decode([]byte(address))

	if len(addressByte) < 4{
		return false
	}

	//2.取数据
	payload := addressByte[:len(addressByte) - 4]
	checksum1 := addressByte[len(addressByte)-4:]

	//3.做checkSum函数
	checksum2 := CheckSum(payload)
	fmt.Printf("checksum1 : %x\n", checksum1)
	fmt.Printf("checksum2 : %x\n", checksum2)

	//4.比较
	return bytes.Equal(checksum1, checksum2)
}

