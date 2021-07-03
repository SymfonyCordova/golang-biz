package flag

import (
	"flag"
	"fmt"
	"log"
)

/*
flag包实现了命令行参数的解析。
*/

//使用flag.String(), Bool(), Int()等函数注册flag，下例声明了一个整数flag，解析结果保存在*int指针ip里：
func TestFlag()  {
	ip := flag.Int("flagname", 1234, "help message for flagname")
	log.Println(ip)
}

//mysql client cmd
func MySqlClientCmd(){
	//第一个参数是命令行key，第二个参数是默认值，第三个参数是 mysql -h 提示
	//var user = flag.String("user", "root", "用户名")
	//var port = flag.Int("port", 3306, "端口")
	//var ip = flag.String("host", "localhost", "主机地址")

	//上述你可能也发现了问题，需要用*变量才能取到值，是不是感觉不太方便，那就来看看flag.TypeVar()。
	//声明变量用于接收命令行参数
	var user string
	var port int
	var ip string

	//从命令行扫描参数赋值到变量
	flag.StringVar(&user, "user", "root", "用户名")
	flag.IntVar(&port, "port", 3306, "端口")
	flag.StringVar(&ip, "l", "localhost", "mysql ip")


	//必须使用flag.Parse()解析一下命令行参数
	flag.Parse()

	fmt.Println(user, port, ip)
	// test command ./test -l 127.0.0.1 -port 3307 -user root


	///////////////// 其他方法
	//返回命令行参数后的其他参数
	fmt.Println(flag.Args())
	//返回命令行参数后的其他参数个数
	fmt.Println(flag.NArg())  //3
	//返回使用的命令行参数个数
	fmt.Println(flag.NFlag()) //3

	// test command ./test -l 127.0.0.1 -port 3307 -user root 11 22 33
}