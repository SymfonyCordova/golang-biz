package runtime

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

/*
go的pprof工具可以用来监测进程的运行数据，用于监控程序的性能，对内存使用和CPU使用的情况统信息进行分析。
官方提供了两个包：runtime/pprof和net/http/pprof，前者用于普通代码的性能分析，后者用于web服务器的性能分析。

runtime/pprof的使用
	该包提供了一系列用于调试信息的方法，可以很方便的对堆栈进行调试
	通常用得多得是以下几个：
		StartCPUProfile：开始监控cpu。
		StopCPUProfile：停止监控cpu，使用StartCPUProfile后一定要调用该函数停止监控。
		WriteHeapProfile：把堆中的内存分配信息写入分析文件中。
*/

const (
	Cols int = 1000
	Rows int  = 1000
)

func CreateCpuAndMemProFile(){
	var cpuProfile string
	var memProfile string

	//从命令行扫描参数赋值到变量
	flag.StringVar(&cpuProfile, "cpuprofile", "", "write cpu profile `file`")
	flag.StringVar(&memProfile, "memprofile", "", "write memory profile `file`")

	flag.Parse()

	if cpuProfile == ""{
		return
	}

	if memProfile == ""{
		return
	}


	cpuFile, err := os.Create(cpuProfile)
	if err != nil {
		log.Fatal("could not create CPU profile ")
		return
	}
	defer cpuFile.Close()

	err = pprof.StartCPUProfile(cpuFile) //开始监控cpu
	if err != nil{
		log.Fatal("could not start CPU profile ")
		return
	}
	defer pprof.StopCPUProfile() //停止监控cpu

	// 主逻辑区，进行一些简单的代码运算
	x := [Rows][Cols]int{}
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < Rows; i++{
		for j := 0; j < Cols; j++ {
			x[i][j] = s.Intn(100000)
		}
	}


	for i := 0; i < Rows; i++{
		tmp := 0
		for j := 0; j < Cols; j++ {
			tmp += x[i][j]
		}
	}


	memfile, err:= os.Create(memProfile)
	if err != nil {
		log.Fatal("could not create Memory profile ")
	}
	defer memfile.Close()

	runtime.GC() //会让运行时系统进行一次强制性的垃圾收集 获取最新的数据信息

	err = pprof.WriteHeapProfile(memfile)
	if err != nil{
		log.Fatal("could not write memory profile: ", err)
		return
	}

	//go build
	// ./pprof -cpuprofile cpu.prof -memprofile mem.prof
	// ./test -cpuprofile cpu.prof -memprofile mem.prof
	//生成数据文件后使用go tool pprof file进入交互式界面进行数据分析，输入help可以查看命令。 go tool pprof cpu.prof
	//
	//top 命令格式：top [n]，查看排名前n个数据，默认为10。
	//tree [n]，以树状图形式显示，默认显示10个。
	//以web形式查看，在web服务的时候经常被用到，需要安装gv工具，官方网页：http://www.graphviz.org/。

	//在web服务器中监测只需要在import部分加上监测包即可：
	//import(
	//	_ "net/http/pprof"
	//)
	//当服务开启后，在当前服务环境的http://ip:port/debug/pprof页面可以看到当前的系统信息：
	//通常可以对服务器在一段时间内进行数据采样，然后分析服务器的耗时和性能: go tool pprof http://*:*/debug/pprof/profile
		//go tool pprof http://127.0.0.1:8080/debug/pprof/profile
}