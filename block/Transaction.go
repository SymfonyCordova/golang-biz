package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

const reward = 50.10

//交易结构
type Transaction struct {
	TXID []byte //交易ID
	TXInputs []TXInput //交易输入数组
	TXOutputs []TXOutput //交易输出数组
}

//交易输入
type TXInput struct {
	//引用的交易ID
	TXid []byte
	//引用的output的索引值
	Index int64
	//解锁脚本
	//数字签名 由 r, s拼成的[]byte
	Signature []byte
	//约定,这里的PubKey不存储原始的公钥,而是存储X和Y的拼接的字符串,在校验端重新拆分(参考r, s传递)
	//注意是公钥 不是哈希,也不是地址
	PubKey []byte
}

//交易输出
type TXOutput struct {
	//转帐金额
	TransferAmount float64
	//锁定脚本
	//收款方的公钥哈希 注意 是哈希而不是公钥 也不是地址
	PubKeyHash []byte
}

//由于现在存储的字段是地址的公钥哈希,所以无法直接创建TXOutput,
//为了能够得到公钥哈希,我们需要处理一下,写一个Lock函数
func (output *TXOutput)Lock(address string){
	//3.真正的锁定动作
	output.PubKeyHash = GetPubKeyFromAddress(address)
}

//给TXOutput提供一个创建的方法,否则无法调用Lock
func NewTxOutput(value float64, address string) *TXOutput{
	output := TXOutput{
		TransferAmount: value,
	}

	output.Lock(address)
	return &output
}

//设置交易ID
func (tx *Transaction)SetHash(){
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

//创建交易
//创建挖矿交易
func NewCoinbaseTx(address string, data string) *Transaction {
	//挖矿交易特点
	//只有一个input
	//无需引用交易id
	//矿工由于挖矿时无需指定签名,所以这个sig字段可以由矿工自由填写,一般是填写矿池的名字
	input := TXInput{
		TXid: []byte{},
		Index:-1,
		Signature: nil,
		PubKey: []byte(data),
	}

	output := NewTxOutput(reward, address)

	//对于挖矿交易来说,只有一个input和一个output
	tx := Transaction{
		TXID:[]byte{},
		TXInputs: []TXInput{input},
		TXOutputs: []TXOutput{*output},
	}
	tx.SetHash()

	return &tx
}

//判断当前交易是否为挖矿交易
func (tx *Transaction)IsCoinbase()bool{
	//对于挖矿交易来说,只有一个input和一个output 交易index为-1
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index == -1 {
		return true
	}

	return false
}

//创建普通的转账交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction{
	//1.创建交易之后要进行数字签名->所以需要私钥->打开钱包
	ws := NewWallets()
	//2.找到自己的钱包 根据地址返回自己的wallet
	wallet := ws.WalletsMap[from]
	if wallet == nil{
		log.Print("没有找到该地址的钱包,交易失败")
		return nil
	}
	//3.得到对应的公钥,私钥
	pubKey := wallet.PubKey
	//TODO privateKey := wallet.PrivateKey
	pubKeyHash := HashPubKey(pubKey)

	//找到最合理UTXO集合 map[string]uint64
	utxos, resValue := bc.FindNeedUTXOs(pubKeyHash, amount)
	if resValue < amount {
		log.Println("余额不足")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput

	//将这些utxo逐一转成inputs
	for id, indexArray := range utxos {
		for _, i := range indexArray{
			input := TXInput{
				TXid:[]byte(id),
				Index: int64(i),
				Signature: nil,
				PubKey: pubKey,
			}
			inputs = append(inputs, input)
		}
	}

	//创建交易输出outputs
	output := NewTxOutput(amount, to)
	outputs = append(outputs, *output)

	if resValue > amount {
		//如果有零钱,要找零 创建找零outputs
		output = NewTxOutput(resValue-amount, from)
		outputs = append(outputs, *output)
	}

	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()

	return &tx
}

