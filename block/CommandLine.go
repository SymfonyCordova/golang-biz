package block

import (
	"fmt"
)

//正向打印
func (cli *BlockChainCli) PrintBlockChain() {
	cli.Bc.PrintChain()
}

//反向打印
func (cli *BlockChainCli) PrintBlockChainReverse() {
	it := cli.Bc.NewIterator()
	for {
		block := it.Next()

		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *BlockChainCli) GetBalance(address string) {
	//校验地址 是否有效
	if !IsValidAddress(address) {
		fmt.Printf("地址无效:%s\n", address)
		return
	}
	// 生成公钥哈希
	pubKeyHash := GetPubKeyFromAddress(address)

	utxos := cli.Bc.FindUTXOs(pubKeyHash)

	total := 0.0
	for _, utxo := range utxos {
		total += utxo.TransferAmount
	}

	fmt.Printf("%s 余额为: %f\n", address, total)
}

func (cli *BlockChainCli) Send(from string, to string, amount float64, miner string, data string) {
	//校验地址 是否有效
	if !IsValidAddress(from) {
		fmt.Printf("地址无效from:%s\n", from)
		return
	}

	if !IsValidAddress(to) {
		fmt.Printf("地址无效to:%s\n", to)
		return
	}

	if !IsValidAddress(miner) {
		fmt.Printf("地址无效miner:%s\n", miner)
		return
	}

	//创建挖矿交易
	coinbase := NewCoinbaseTx(miner, data)
	//创建一个普通交易
	tx := NewTransaction(from, to, amount, cli.Bc)
	if tx == nil {
		return
	}
	//添加到区块
	cli.Bc.AddBlock([]*Transaction{coinbase, tx})
	fmt.Printf("转帐成功！")
}

func (cli *BlockChainCli) NewWallet() {
	ws := NewWallets()
	address := ws.CreateWallet()
	fmt.Printf("地址: %s\n", address)
}

func (cli *BlockChainCli) ListAddresses() {
	ws := NewWallets()
	addresses := ws.ListAllAddress()

	for _, address := range addresses {
		fmt.Printf("地址: %s\n", address)
	}
}
