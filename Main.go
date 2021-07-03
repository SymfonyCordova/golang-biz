package main

import (
	"github.com/SymfonyCordova/golang-biz/block"
)

func main() {
	bc := block.NewBlockChain("1F3fxNMeXDY16sN9k9jdrZtMKdcVEdrMaM")
	cli := block.BlockChainCli{Bc: bc}
	cli.Run()
}
