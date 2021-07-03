https://blog.csdn.net/wxd1234567890/article/details/116308372
https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-defer/

1.值类型
bool
    int(32 or 64), int8, int16, int32, int64
    uint(32 or 64), uint8(byte), uint16, uint32, uint64
    float32, float64
    string
    complex64, complex128
    array    -- 固定长度的数组


2.字符
    Go 语言的字符有以下两种
        uint8类型，或者叫byte 型，代表了ASCII码的一个字符。
        rune类型，代表一个 UTF-8字符。
            当需要处理中文、日文或者其他复合字符时，则需要用到rune类型。
            rune类型实际是一个int32。
            Go 使用了特殊的 rune 类型来处理 Unicode，让基于 Unicode的文本处理更为方便，也可以使用 byte 型进行默认字符串处理，性能和扩展性都有照顾
            因为UTF8编码下一个中文汉字由3~4个字节组成，所以我们不能简单的按照字节去遍历一个包含中文的字符串,否则就会出现wenti
            一个rune字符由一个或多个byte组成
            rune类型用来表示utf8字符，一个rune字符由一个或多个byte组成

3.字符串
    字符串底层是一个byte数组，所以可以和[]byte类型相互转换。
    字符串是只读不能修改的，字符串是由byte字节组成，所以字符串的长度是byte字节的长度
    要修改字符串，需要先将其转换成[]rune或[]byte，完成后再转换为string。无论哪种转换，都会重新分配内存，并复制字节数组
    常用操作
        len(str) 求长度
        +或fmt.Sprintf 拼接字符串
        strings.Split 分割
        strings.Contains 判断是否包含
        strings.HasPrefix,strings.HasSuffix  前缀/后缀判断
        strings.Index(),strings.LastIndex() 子串出现的位置
        strings.Join(a[]string, sep string) join操作

    遍历字符串
          func traversalString() {
                s := "pprof.cn博客"
                for i := 0; i < len(s); i++ { //byte
                    fmt.Printf("%v(%c) ", s[i], s[i])
                }
                fmt.Println()
                for _, r := range s { //rune
                    fmt.Printf("%v(%c) ", r, r)
                }
                fmt.Println()
            }
            //err
            112(p) 112(p) 114(r) 111(o) 102(f) 46(.) 99(c) 110(n) 229(å) 141() 154() 229(å) 174(®) 162(¢)
            112(p) 112(p) 114(r) 111(o) 102(f) 46(.) 99(c) 110(n) 21338(博) 23458(客)

    修改字符串
        s1 := "123"
        byteS1 := []byte(s1)
        byteS1[0] = 'a'
        fmt.Print(string(byteS1))

        s2 := "中国123"
        byteS2 := []rune(s2)
        byteS2[0] = '我'
        fmt.Print(string(byteS2))

数组
    func arrayTest() {
    	a := [2]int{1, 2}
    	b := [...]int{1, 2, 3}
    	c := [4]int{0: 1, 3: 10}
    	d := [...]struct {
    		name string
    		age  uint8
    	}{
    		{"a1", 10},
    		{"a2", 10},
    	}
    	fmt.Println(a, b, c, d)
    	//
    	aM := [2][3]int{{1, 2, 3}, {4, 5, 6}}
    	bM := [...][2]int{{1, 2}} //纬度不能用 "..."。
    	fmt.Println(aM, bM)

    }
    // 值拷贝和传指针
    func test(x [2]int) {
    	fmt.Printf("x: %p\n", &x)
    	x[1] = 1000
    }
    func test2(arr *[5]int) {
    	arr[0] = 10
    }

    func testCopy() {
    	a := [2]int{}
    	fmt.Printf("a: %p\n", &a)
    	//值拷贝
    	test(a)
    	fmt.Println(a)

    	var arr1 [5]int
    	//数组指针
    	test2(&arr1)
    	fmt.Println(arr1)
    	arr2 := [...]int{2, 4, 6, 8, 10}
    	test2(&arr2)
    	fmt.Println(arr2)
    }
    //a: 0xc000064210
    //x: 0xc000064220
    //[0 0]
    //[10 0 0 0 0]
    //[10 4 6 8 10]
    //数组指针和指针数组  ....

引用类型：(指针类型)
     slice   -- 动态长度的数组
     map     -- 映射
     chan    -- 管道

内置函数
       append          -- 用来追加元素到数组、slice中,返回修改后的数组、slice
        close           -- 主要用来关闭channel
        delete            -- 从map中删除key对应的value
        panic            -- 停止常规的goroutine  （panic和recover：用来做错误处理）
        recover         -- 允许程序定义goroutine的panic动作
        real            -- 返回complex的实部   （complex、real imag：用于创建和操作复数）
        imag            -- 返回complex的虚部
        make            -- 用来分配内存，返回Type本身(只能应用于slice, map, channel)
        new                -- 用来分配内存，主要用来分配值类型，比如int、struct。返回指向Type的指针
        cap                -- capacity是容量的意思，用于返回某个类型的最大容量（只能用于切片和 map）
        copy            -- 用于复制和连接slice，返回复制的数目
        len                -- 来求长度，比如string、array、slice、map、channel ，返回长度

对于引用类型（slice、map、chan）的变量，我们在使用的时候不仅要声明它，还要为它分配内存空间，否则我们的值就没办法存储。
而对于值类型的声明不需要分配内存空间，是因为它们在声明的时候已经默认分配好了内存空间。要分配内存，new和make
   make只用于slice、map以及channel的初始化，返回的还是这三个引用类型本身
    make函数是无可替代的，我们在使用slice、map以及channel的时候，都需要使用make进行初始化，然后才可以对它们进行操作
   new用于值类型和struct的内存分配，并且内存对应的值为类型零值，返回的是指向类型的指针。根据传入的类型分配一片内存空间并返回指向这片内存空间的指针。使用new函数得到的是一个类型的指针


defer
    会在当前函数返回前执行传入的函数，它会经常被用于关闭文件描述符、关闭数据库连接、关闭http的响应body，回滚数据库的事务以及解锁资源
    func main() {
    	startedAt := time.Now()
    	defer fmt.Println(time.Since(startedAt))
    	time.Sleep(time.Second)
    }
    $ go run main.go
    0s
    调用 defer 关键字会立刻拷贝函数中引用的外部参数，所以 time.Since(startedAt) 的结果不是在 main 函数退出之前计算的，而是在 defer 关键字调用时计算的，最终导致上述代码输出 0s
    func main() {
    	startedAt := time.Now()
    	defer func() { fmt.Println(time.Since(startedAt)) }()
    	time.Sleep(time.Second)
    }
    $ go run main.go
    1s
    type _defer struct {
    	siz       int32//参数和结果的内存大小
    	started   bool
    	openDefer bool
    	sp        uintptr//栈指针
    	pc        uintptr//调用方的程序计数器
    	fn        *funcval//defer关键字中传入的函数
    	_panic    *_panic
    	link      *_defer // 延迟调用链表
    }
    编译期将 defer 关键字被转换 runtime.deferproc 并在调用 defer 关键字的函数返回之前插入 runtime.deferreturn；
    运行时调用 runtime.deferproc 会将一个新的 runtime._defer 结构体追加到当前 Goroutine 的链表头；
    运行时调用 runtime.deferreturn 会从 Goroutine 的链表中取出 runtime._defer 结构并依次执行；

    栈上分配
        当defer关键字在函数体中最多执行一次时，编译期间的 cmd/compile/internal/gc.state.call 会将结构体分配到栈上并调用
    开放编码 · 1.14
        通过开放编码（Open Coded）实现 defer 关键字，该设计使用代码内联优化 defer 关键的额外开销并引入函数数据 funcdata 管理 panic 的调用。
        开放编码只会在满足以下的条件时启用：
            函数的 defer 数量少于或者等于 8 个；
            函数的 defer 关键字不能在循环中执行；
            函数的 return 语句与 defer 语句的乘积小于或者等于 15 个；
            编译期间判断 defer 关键字、return 语句的个数确定是否开启开放编码优化
            通过 deferBits 和 cmd/compile/internal/gc.openDeferInfo 存储 defer 关键字的相关信息
            如果 defer 关键字的执行可以在编译期间确定，会在函数返回前直接插入相应的代码，否则会由运行时的 runtime.deferreturn 处理

结构体
    Tag
    Tag是Struct的一部分，只有在反射场景中才有用，而反射包中提供了操作Tag的方法。常见的tag用法，主要是JSON数据解析、ORM映射等。struct结合反射获取tag中的键值
    type Server struct {
        ServerName string `key1:"value1" key11:"value11"`
        ServerIP   string `key2:"value2"`
    }

    func main() {
        s := Server{}
        st := reflect.TypeOf(s)

        field1 := st.Field(0)
        fmt.Printf("key1:%v\n", field1.Tag.Get("key1"))
        fmt.Printf("key11:%v\n", field1.Tag.Get("key11"))

        filed2 := st.Field(1)
        fmt.Printf("key2:%v\n", filed2.Tag.Get("key2"))
    }

interface{}
    任意的结构体都能转换为空接口类型
    任何类型都可以被Any类型引用，Any类型就是空接口，即interface{}
    c.(类型) 类型转换 类型断言

反射
    很多框架（gorm、json）都依赖 Go 语言的反射机制简化代码。
    因为 Go 语言的语法元素很少、设计简单，所以它没有特别强的表达能力，
    但是 Go 语言的 reflect 包能够弥补它在语法上reflect.Type的一些劣势。
    reflect 实现了运行时的反射能力，能够让程序操作不同类型的对象。
    反射包中有两对非常重要的函数和类型，两个函数分别是：
        reflect.TypeOf 能获取类型信息；
        reflect.ValueOf 能获取数据的运行时表示；

gdb 调试工具
    go build -gcflags "-N -l" gdbfile.go
    gdb gdbfile

单元测试
    testing内置库，go test
性能压力测试
    压力测试用来检测函数(方法）的性能。go test不会默认执行压力测试的函数，
    如果要执行压力测试需要带上参数-test.bench，
    语法:test.bench=“test_name_regex”,例如go test -test.bench=".*"表示测试全部的压力测试函数

json
    写字母开头的字段成员是无法被外部直接访问的，所以 struct 在进行 json、xml、gob 等格式的 encode 操作时，这些私有字段会被忽略。

http-client
    使用 HTTP 标准库发起请求、获取响应时，即使你不从响应中读取任何数据或响应为空，都需要手动关闭响应体。应该先检查 HTTP 响应错误为 nil，再调用 resp.Body.Close() 来关闭响应体：
    func main() {
        resp, err := http.Get("http://www.baidu.com")
        // 关闭 resp.Body 的正确姿势
        if resp != nil {
            defer resp.Body.Close()
        }
        checkError(err)
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        checkError(err)

        fmt.Println(string(body))
    }

切片扩容的策略：
    如果切片的容量小于 1024 个元素，于是扩容的时候就翻倍增加容量。
    一旦元素个数超过 1024 个元素，那么增长因子就变成 1.25 ，即每次增加原来容量的四分之一。
    注意：扩容扩大的容量都是针对原来的容量而言的，而不是针对原来数组的长度而言的

关闭channel
    关闭channel时会把recvq中的G全部唤醒，本该写入G的数据位置为nil。把sendq中的G全部唤醒，但这些G会panic
    除此之外，panic出现的常见场景还有
        关闭值为nil的channel
        关闭已经被关闭的channel
        向已经关闭的channel写数据
        ha := a
        hv1, hb := <-ha
        for ; hb != false; hv1, hb = <-ha { //如果不存在当前值，意味着当前的管道已经被关闭 //如果存在当前值，会为 v1 赋值并清除 hv1 变量中的数据，然后重新陷入阻塞等待新数据
            v1 := hv1
            hv1 = nil
            ...
        }
channel
    道像一个传送带或者队列 先入先出（First In First Out）的规则，保证收发数据的顺序

select
    select 也能够让 Goroutine 同时等待多个 Channel 可读或者可写，在多个文件或者 Channel状态改变之前，select 会一直阻塞当前线程或者 Goroutine
      select {
        case <-chan1:
           // 如果chan1成功读到数据，则进行该case处理语句
        case chan2 <- 1:
           // 如果成功向chan2写入数据，则进行该case处理语句
        default:
           // 如果上面都没有成功，则进入default处理流程
        }
        1、select的非阻塞：通过default保证。
        2、随机性：select的两个case如果都是同时满足执行条件的，如果我们按照顺序依次判断，那么后面的条件永远都会得不到执行，而随机的引入就是为了避免饥饿问题的发生

Mutex
    当一个变量被上了互斥锁后，其他访问该变量的线程会被堵塞，不可对该变量进行读写操作，直到锁被释放。互斥锁是一种常用的控制共享资源访问的方法，它能够保证同时只有一个goroutine可以访问共享资源

RWMutex
    RWMutex是基于互斥锁Mutex实现的，包含了读锁Rlock()和写锁Lock()，上读锁时，数据可以被多个goroutine并发访问但不可写，而上写锁时，数据不可被其他goroutine读或写。读写锁非常适合读多写少的场景。

sync.WaitGroup
    waitgroup在golang中的实现都依赖于 原子操作 & 信号量，go中的信号量是在runtime包中实现的。Golang中的信号量，提供了goroutine的阻塞和唤醒机制
        var x int64
        var wg sync.WaitGroup

        func add() {
            for i := 0; i < 5000; i++ {
                x = x + 1
            }
            wg.Done() //-1
        }
        func main() {
            wg.Add(2)
            go add()
            go add()
            wg.Wait() //同步等待完成
            fmt.Println(x)
        }

sync.Once
    需要确保某些操作在高并发的场景下只执行一次，例如只加载一次配置文件、只关闭一次通道等。
    sync.Once其实内部包含一个互斥锁和一个布尔值，互斥锁保证布尔值和数据的安全，而布尔值用来记录初始化是否完成。
    这样设计就能保证初始化操作的时候是并发安全的并且初始化操作也不会被执行多次。
    加载载配置文件。延迟一个开销很大的初始化操作到真正用到它的时候再执行是一个很好的实践。
    因为预先初始化一个变量（比如在init函数中完成初始化）会增加程序的启动耗时，而且有可能实际执行过程中这个变量没有用上，那么这个初始化操作就不是必须要做的。

sync.NewCond

sync.Pool
    大量重复地创建许多对象，造成 GC 的工作量巨大，CPU 频繁掉底。可以使用 sync.Pool 来缓存对象，减轻 GC 的消耗。可以作为保存临时取还对象的一个“池子”。
    sync.Pool 是协程安全的，这对于使用者来说是极其方便的。使用前，设置好对象的 New 函数，用于在 Pool 里没有缓存的对象时，创建一个。
    之后，在程序的任何地方、任何时候仅通过 Get()、Put() 方法就可以取、还对象了。pool作为临时对象池子

sync.Map
    Go语言中内置的map不是并发安全的。Go语言的sync包中提供了一个开箱即用的并发安全版map–sync.Map

atomic
    代码中的加锁操作因为涉及内核态的上下文切换会比较耗时、代价比较高。
    针对基本数据类型我们还可以使用原子操作来保证并发安全，因为原子操作是Go语言提供的方法它在用户态就可以完成，
    因此性能比加锁操作更好。Go语言中原子操作原子类由内置的标准库sync/atomic提供。
    这里原子操作，是保证多个cpu（协程）对同一块内存区域的操作是原子的。CAS操作，为了保证原子性，golang是通过汇编指令来实现的

M P G
    M:内核线程(thread线程)
    P: Processor调度器
    G: (goroutine协程)
        M 与 P 是 1：1 的关系
        P 与 G 是 1：n 的关系
全局队列（Global Queue）：存放等待运行的 G。
P 的本地队列：同全局队列类似，存放的也是等待运行的 G，存的数量有限，不超过 256 个。新建 G’时，G’优先加入到 P 的本地队列，如果队列满了，则会把本地队列中一半的 G 移动到全局队列。
P 列表：所有的 P 都在程序启动时创建，并保存在数组中，最多有 GOMAXPROCS(可配置) 个。
M：线程想运行任务就得获取 P，从 P 的本地队列获取 G，P 队列为空时，M 也会尝试从全局队列拿一批 G 放到 P 的本地队列，或从其他 P 的本地队列偷一半放到自己 P 的本地队列。M 运行 G，G 执行之后，M 会从 P 获取下一个 G，不断重复下去

调度器P的策略:
    复用线程：避免频繁的创建、销毁线程，而是对线程的复用。
    1）work stealing 机制
        当本线程无可运行的 G 时，尝试从其他线程绑定的 P 偷取 G，而不是销毁线程。
    2）hand off 机制
        当本线程因为 G 进行系统调用阻塞时，线程释放绑定的 P，把 P 转移给其他空闲的线程执行。
        利用并行：GOMAXPROCS 设置 P 的数量，最多有 GOMAXPROCS 个线程分布在多个 CPU 上同时运行。GOMAXPROCS 也限制了并发的程度，比如 GOMAXPROCS = 核数/2，则最多利用了一半的 CPU 核进行并行。
        抢占：在 coroutine 中要等待一个协程主动让出 CPU 才执行下一个协程，在 Go 中，一个 goroutine 最多占用 CPU 10ms，防止其他 goroutine 被饿死，这就是 goroutine 不同于 coroutine 的一个地方。
    全局 G 队列：在新的调度器中依然有全局 G 队列，但功能已经被弱化了，当 M 执行 work stealing 从其他 P 偷不到 G 时，它可以从全局 G 队列获取 G

debug.SetMaxThreads(threads int)//设置最多可以创建多少操作系统线程M
runtime.GOMAXPROCS(1)//设置逻辑处理器个数,也就是设置GMP中的P的个数 P
runtime.Gosched()//类似Java中线程的yeild方法，当一个goroutine执行该方法时候意味着当前goroutine放弃当前cpu的使用权，然后运行时会调度系统会调度其他goroutine占用cpu进行运行，放弃CPU使用权的goroutine并没有被阻塞，而是处于就绪状态，可以在随时获取到cpu情况下继续运行



