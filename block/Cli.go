package block

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

//这是一个用来接收命令并且控制区块链操作的文件
type BlockChainCli struct {
	Bc *BlockChain
}

const Usage = `golang-biz version v1.0.0 

command:
	printChain                      "正向打印区块链"
	printChainReverse               "反向打印区块链"
	getBalance --address ADDRESS    "获取指定地址的余额"
	send FROM TO AMOUNT MINER DATA  "有FROM转AMOUNT给TO,由MINTER挖矿,同时写入DATA"
	newWallet			"创建一个新的钱包(私钥公钥对)"
	listAddresses	                 "列举所有钱包地址"
`

//接收参数的动作,我们放到一个函数中
// 分析命令 执行相应的动作
func (cli *BlockChainCli) Run(){
	args := os.Args

	//得到所有的命令
	if len(args) < 2{
		fmt.Printf(Usage)
		return
	}

	cmd := args[1]
	switch cmd {
	case "printChain":
		//打印区块
		cli.PrintBlockChain()
		break
	case "printChainReverse":
		cli.PrintBlockChainReverse()
		break
	case "getBalance":
		if len(args) == 4 && args[2] == "--address" {
			address := args[3]
			cli.GetBalance(address)
		}else{
			fmt.Println("invalid params")
			fmt.Printf(Usage)
		}

		break
	case "send":
		if len(args) != 7{
			fmt.Println("invalid params")
			fmt.Printf(Usage)
		}
		from := args[2]
		to := args[3]
		amount, err := strconv.ParseFloat(args[4], 64)
		if err != nil {
			log.Panic(err)
		}
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	case "newWallet":
		fmt.Printf("创建一个新的钱包(私钥公钥对)......\n")
		cli.NewWallet()
	case "listAddresses":
		fmt.Printf("列举所有钱包地址......\n")
		cli.ListAddresses()
	default:
		//无效的命令请检查
		fmt.Printf(Usage)
		break
	}
}








