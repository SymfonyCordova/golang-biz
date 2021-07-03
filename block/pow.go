package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
)

func CalcTotalBitCoin(){
	//1.每首21万个块减半
	//2.最初奖励50比特币
	//3.用一个循环来判断
	total := 0.0
	blockInterval := 21.0
	currentRecord := 50.0

	for currentRecord > 0 {
		//每个区间内的总量
		amount1 := blockInterval * currentRecord
		//除效率低,我们使用等价的乘法
		currentRecord *= 0.5
		total += amount1
	}

	fmt.Println("比特币总量：", total, "万")
}

//工作量证明
func Pow(){
	data := "hello world"
	for i:=0; i<1000000; i++ {
		hash := sha256.Sum256([]byte(string(rune(i)) + data))
		fmt.Printf("hash: %x\n", hash[:])
	}
}

func Join(){
	str1 := []string{"hello", "world", "!"}
	res := strings.Join(str1, "")
	fmt.Printf("res : %s\n", res)

	var res1 = bytes.Join([][]byte{[]byte("hello"), []byte("world")}, []byte("+"))
	fmt.Printf("res1 : %s\n", res1)
}

