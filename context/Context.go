package context

//Context是协程安全的 代码中可以将单个Context传递给任意数量的goroutine,并在取消该Context时可以将信号传递给所有的goroutine
//Context可以派生 组成一个树 父的context退出,子也会退出

/*
type Context interface {
    Deadline() (deadline time.Time, ok bool)
		//是获取设置的截止时间的意思，第一个返回式是截止时间，到了这个时间点，Context会自动发起取消请求；
		//第二个返回值ok==false时表示没有设置截止时间，如果需要取消的话，需要调用取消函数进行取消
    Done() <-chan struct{}
		//返回一个只读的chan，类型为struct{}，
		//我们在goroutine中，如果该方法返回的chan可以读取，则意味着parent context已经发起了取消请求，
		//我们通过Done方法收到这个信号后，就应该做清理操作，然后退出goroutine，释放资源
    Err() error
		//返回取消的错误原因，因为什么Context被取消
    Value(key interface{}) interface{}
		//获取该Context上绑定的值，是一个键值对，所以要通过一个Key才可以获取对应的值，这个值一般是线程安全的
		//用于获取特定于当前任务树的额外信息
}
*/

/*
可以看到Done方法返回的channel正是用来传递结束信号以抢占并中断当前任务；
Deadline方法指示一段时间后当前goroutine是否会被取消；
以及一个Err方法，来解释goroutine被取消的原因；
而Value则用于获取特定于当前任务树的额外信息。
而context所包含的额外信息键值对是如何存储的呢？
	其实可以想象一颗树，树的每个节点可能携带一组键值对,如果当前节点上无法找到key所对应的值，就会向上去父节点里找，直到根节点

context.Background() 树的根节点(父节点) 属于main函数的
context.TODO()
//Context使用原则
//不要把Context放在结构体中，要以参数的方式进行传递
//以Context作为参数的函数方法，应该把Context作为第一个参数，放在第一位
//给一个函数方法传递Context的时候，不要传递nil，如果不知道传递什么，就使用context.TODO()
//Context的Value相关方法应该传递必须的数据，不要什么数据都使用这个传递

在没用呢? 父 cancel()  子 <-ctx.done

*/

//context包为我们提供的With系列的函数了 						根据父节点context生成子节点context
//func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
//传递一个父Context作为参数，返回子Context，以及一个取消函数用来取消Context
//func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
//和WithCancel差不多，它会多传递一个截止时间参数，意味着到了这个时间点，会自动取消Context，当然我们也可以不等到这个时候，可以提前通过取消函数进行取消
//func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
//和WithDeadline基本上一样，这个表示是超时自动取消，是多少时间后自动取消Context的意思
//func WithValue(parent Context, key, val interface{}) Context
//函数和取消Context无关，它是为了生成一个绑定了一个键值对数据的Context，这个绑定的数据可以通过Context.Value方法访问到
//
import (
	"context"
	"fmt"
	"time"
)

func a(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	i := 0
	ctx, cancle := context.WithCancel(ctx)
	for _ = range ticker.C {
		select {
		case <-ctx.Done():
			fmt.Println("a break")
			return
		default:
			i++
			if i == 5 {
				fmt.Println(i, "start cancle")
				cancle()
			}
			fmt.Println("a send num = ", 10)
			go b(ctx, 10)
		}
	}
}

func b(ctx context.Context, num int) {
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		select {
		case <-ctx.Done():
			fmt.Println("b break")
			return
		default:
			fmt.Println("b sen num = ", num)
			go c(ctx, num)
		}
	}
}

func c(ctx context.Context, num int) {
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		select {
		case <-ctx.Done():
			fmt.Println("c break")
			return
		default:
			fmt.Println(num)
		}
	}
}

func DemoRelease() {
	//ctx, _ := context.WithCancel(context.Background())

	go a(context.TODO())

	//fmt.Println(".")
	//time.Sleep(10*time.Second)
	//cancel()

	for true {
		fmt.Println("..")
		time.Sleep(1 * time.Second)
	}
}

var key1 string = "key1"
var key2 string = "key2"

func DemoContextKey() {
	ctx, cancel := context.WithCancel(context.Background())

	//value
	valueCtx := context.WithValue(ctx, key1, "key1 value")
	valueCtx2 := context.WithValue(valueCtx, key2, "key2 value")

	go watch(valueCtx2)

	cancel()
}

func watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Value(key1))
			fmt.Println(ctx.Value(key2))
			break
		default:
			fmt.Println(ctx.Value(key1), "...")
		}
	}
}
