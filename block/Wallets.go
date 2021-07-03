package block

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	bizEncoding "github.com/SymfonyCordova/golang-biz/encoding"
	"io/ioutil"
	"log"
	"os"
)

const WalletFile = "wallet.bat"

//定义一个Wallets结构,它保存所有的wallet以及它的地址
type Wallets struct {
	WalletsMap map[string]*Wallet
}

//创建方法
func NewWallets()*Wallets{
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)
	err := ws.loadFile()
	if err != nil {
		log.Panic(err)
	}
	return &ws
}

//读取文件方法,把所有的wallet读出来
func (ws *Wallets)loadFile() error {
	_, err := os.Stat(WalletFile)
	if os.IsNotExist(err){
		return err
	}

	file, err := ioutil.ReadFile(WalletFile)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(file))

	var wss Wallets
	err = decoder.Decode(&wss)
	if err != nil {
		return err
	}

	//ws = &wss
	ws.WalletsMap = wss.WalletsMap

	return nil
}

func (ws *Wallets)CreateWallet()string{
	wallet := NewWallet()
	address := wallet.NewAddress()

	ws.WalletsMap[address] = wallet

	err := ws.saveToFile(wallet)
	if err != nil {
		log.Panic(err)
	}
	return address
}

//保存方法,把新建的wallet添加进去
func (ws *Wallets)saveToFile(wallet *Wallet) error {
	var buffer bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(ws)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(WalletFile, buffer.Bytes(), 0600)
	if err != nil {
		return err
	}

	return nil
}

//获取所有钱包
func (ws *Wallets) ListAllAddress()[]string{
	var addresses []string

	//遍历钱包,将所有的key取出来返回
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}

	return addresses
}


//通过地址返回公钥哈希
func GetPubKeyFromAddress(address string)[]byte{
	//1.解码
	addressByte := bizEncoding.Base58Decode([]byte(address))
	//2.截取出公钥哈希: 去除校验码(4字节)
	length := len(addressByte)
	pubKeyHash := addressByte[1:length-4]
	return pubKeyHash
}





