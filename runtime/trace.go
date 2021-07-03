package runtime

import (
	"os"
	"runtime/trace"
	"time"
)

/*
执行追踪器,跟踪器捕获各种各样的时间,
	如go协程的创建、阻塞、解锁，syscall 进入、退出、阻塞、GC相关时间,堆大小变化，处理器启动、停止等，将这些事件写入到io.writor中，大多数时间都会捕获到精确的纳秒精度时间戳

	func Start(w io.Writer) error 开始go执行追踪器 未当前的程序启用追踪器，追踪的数据将会写入w 中，不能重复创建
	func Stop() 停止当前的追踪器，当所有追踪完全写入w后，才返回

	go tool trace xxxx
*/


func CreateTraceFile(){
	file, err := os.Create("./trace_file")
	if err != nil {
		panic(err)
	}

	trace.Start(file)
	defer trace.Stop()

	data := make(chan int)
	go test(data)
	<-data
}


func test(s chan int){
	time.Sleep(time.Second)
	go test2(s)
}
func test2(s chan int){
	time.Sleep(time.Second)
	s <- 3
}

//go tool trace trace_file
	//系统会自动启动浏览器 Opening browser. Trace viewer is listening on http://127.0.0.1:39877