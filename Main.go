package main

import (
	"github.com/SymfonyCordova/golang-biz/block"
)

func main(){
	bc := block.NewBlockChain("1AzfPzSiv2EW6ACm6XDc7Qhg7gtMSBRFPB")
	cli := block.BlockChainCli{ Bc: bc }
	cli.Run()
}
