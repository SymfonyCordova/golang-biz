package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//bufio 包实现了有缓冲的I/O。它包装一个io.Reader或io.Writer接口对象，创建另一个也实现了该接口，且同时还提供了缓冲和一些文本I/O的帮助函数的对象。
func main(){
	//声明并初始化带缓冲的读取器
	//准备从标准输入读取内容
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("please input your name:")
	//以\n为分割符读取一段内容
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("found an error: %s\n", err)
	}else{
		//对input 进行切片操作,去掉内容中最后一个字节\n
		input = input[:len(input)-1]
		fmt.Printf("hello, %s\n", input)
	}

	for {
		input, err = inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("an error occurred: %s\n", err)
			continue
		}
		input = input[:len(input)-1]
		// 全部转换小写
		input = strings.ToLower(input)
		switch input {
		case "":
			continue
		case "noting", "bye":
			fmt.Println("Bye!")
			//正常退出
			os.Exit(0)
		default:
			fmt.Println("Sorry, I did't catch you.")
		}
	}
}